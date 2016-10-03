package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
)

type Container struct {
	Hostname         string
	CntrInfodHttpUrl string
	MacAddress       string
	IpAddress        string
	HostHostname     string
	AgeInSeconds     int64
}

func FetchContainerData(consulUrl string) []Container {

	var containers []Container

	config := consulapi.DefaultConfig()
	config.Address = consulUrl
	consul, err := consulapi.NewClient(config)

	if err != nil {
		fmt.Printf("Error while trying to read container registry: %s\n", err)
		return containers
	}

	kv := consul.KV()

	allContainerKeys, _, err := kv.Keys("containers", "", nil)
	if err != nil {
		fmt.Printf("Error while trying to read container registry: %s\n", err)
		return containers
	}

	instanceIds := DecodeInstanceIds(allContainerKeys)

	for _, instanceId := range instanceIds {

		// cntrInfodUrl
		valueAsBytes := internalGet(kv, "containers/"+instanceId+"/cntrInfodHttpUrl")
		if valueAsBytes == nil {
			return containers
		}
		cntrInfodHttpUrl := string(valueAsBytes)

		// MAC
		valueAsBytes = internalGet(kv, "containers/"+instanceId+"/macAdress")
		if valueAsBytes == nil {
			return containers
		}
		macAddress := string(valueAsBytes)

		// IP
		valueAsBytes = internalGet(kv, "containers/"+instanceId+"/ipAdress")
		if valueAsBytes == nil {
			return containers
		}
		ipAddress := string(valueAsBytes)

		// Unix Epoch Timestamp
		valueAsBytes = internalGet(kv, "containers/"+instanceId+"/unixEpochTimestamp")
		if valueAsBytes == nil {
			return containers
		}
		ageInSeconds := DetermineAgeInSeconds(valueAsBytes)

		// Hostname
		valueAsBytes = internalGet(kv, "containers/"+instanceId+"/hostname")
		if valueAsBytes == nil {
			return containers
		}
		hostname := string(valueAsBytes)

		// HostHostname
		valueAsBytes = internalGet(kv, "containers/"+instanceId+"/hostinfo/hostname")
		if valueAsBytes == nil {
			return containers
		}
		hostHostname := string(valueAsBytes)

		container := Container{Hostname: hostname, CntrInfodHttpUrl: cntrInfodHttpUrl, MacAddress: macAddress,
			IpAddress: ipAddress, HostHostname: hostHostname, AgeInSeconds: ageInSeconds}

		containers = append(containers, container)
		fmt.Printf("Successfully read info for instance ID %s from container registry.\n", instanceId)
	}

	return containers
}
