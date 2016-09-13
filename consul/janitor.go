package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
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

		config := consulapi.DefaultConfig()
		config.Address = consulUrl
		consul, err := consulapi.NewClient(config)

		if err != nil {
			fmt.Printf("Error while trying to clean up container registry: %s\n", err)
			continue
		}

		kv := consul.KV()

		allContainerKeys, _, err := kv.Keys("containers", "", nil)
		if err != nil {
			fmt.Printf("Error while trying to clean up container registry: %s\n", err)
			continue
		}
		instanceIds := DecodeInstanceIds(allContainerKeys)
		fmt.Printf("%d containers listed in registry.\n", len(instanceIds))

		for _, instanceId := range instanceIds {
			kvp, _, err := kv.Get("containers/"+instanceId+"/unixEpochTimestamp", nil)
			if err != nil {
				fmt.Printf("Error while trying to clean up container registry: %s\n", err)
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
