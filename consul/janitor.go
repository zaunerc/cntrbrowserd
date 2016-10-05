package consul

import (
	"fmt"
	consulapi "github.com/zaunerc/consul/api"
	cc "github.com/zaunerc/go_consul_commons"
	"strconv"
	"strings"
	"time"
)

const MaxInstanceAgeInSeconds = 15
const CleanUpTaskIntervalInSeconds = 5

/**
 * ScheduleRegistration return immediately after the
 * container registration job is scheduled.
 */
func ScheduleCleanUpTask(consulUrl string) {
	fmt.Printf("Scheduling registry cleanup task using consul URL >%s<.\n", consulUrl)
	go runCleanUpTask(consulUrl, CleanUpTaskIntervalInSeconds)
}

func runCleanUpTask(consulUrl string, sleepSeconds int) {
	for {
		fmt.Printf("Container registry will be cleaned up in %d seconds...\n", sleepSeconds)
		time.Sleep(time.Duration(sleepSeconds) * time.Second)

		consul, err := cc.GetConsulClientForUrl(consulUrl)
		if err != nil {
			fmt.Printf("Skipping container registry cleanup run. Error while getting Consul HTTP client: %s\n", err)
			continue
		}

		kv := consul.KV()

		allContainerKeys, _, err := kv.Keys("containers", "", nil)
		if err != nil {
			fmt.Printf("Error while trying to clean up container registry: %s\n", err)
			continue
		}
		if allContainerKeys == nil {
			fmt.Printf("No containers to clean up: Top-level key \"containers/\" does not exist.")
			continue
		}

		instanceIds := DecodeInstanceIds(allContainerKeys)
		fmt.Printf("%d containers listed in registry.\n", len(instanceIds))

		for _, instanceId := range instanceIds {
			key := "containers/" + instanceId + "/unixEpochTimestamp"
			kvp, _, err := kv.Get(key, nil)
			if err != nil {
				fmt.Printf("Error while trying to clean up container registry: %s\n", err)
				continue
			}
			if kvp == nil {
				fmt.Printf("Key >%s< does not exist in registry. Therefore skipping further processing steps specific to this key.", key)
				continue
			}
			ageInSeconds := DetermineAgeInSeconds(kvp.Value)
			fmt.Printf("Instance ID %s registered %d seconds ago.\n", instanceId, ageInSeconds)
			if ageInSeconds > MaxInstanceAgeInSeconds {
				deleteInstanceFromRegistry(instanceId, kv)
			}
		}

		fmt.Printf("Successfully cleaned up container registry.\n")
	}
}

func deleteInstanceFromRegistry(instanceId string, kv *consulapi.KV) {
	_, err := kv.DeleteTree("containers/"+instanceId, nil)
	if err != nil {
		fmt.Printf("Error while trying to remove instance with ID %s from container registry: %s\n", err)
	}
	fmt.Printf("Successfully removed instance with ID %s from container registry.\n", instanceId)
}

func DecodeInstanceIds(containerKeys []string) []string {
	instanceIdMap := make(map[string]string)
	for _, containerKey := range containerKeys {
		keyParts := strings.Split(containerKey, "/")
		if len(keyParts) > 2 {
			instanceIdMap[keyParts[1]] = keyParts[1]
		}
	}

	instanceIds := make([]string, len(instanceIdMap))
	i := 0
	for k := range instanceIdMap {
		instanceIds[i] = k
		i++
	}
	return instanceIds
}

func DetermineAgeInSeconds(unixEpochTimestamp []byte) int64 {
	now := time.Now().Unix()
	unixEpochTimestampAsInt64, _ := strconv.ParseInt(string(unixEpochTimestamp), 10, 64)

	return now - unixEpochTimestampAsInt64
}
