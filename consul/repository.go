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
		kvp, _, err := kv.Get("containers/"+instanceId+"/cntrInfodHttpUrl", nil)
		if err != nil {
			fmt.Printf("Error while trying to read container registry: %s\n", err)
			return containers
		}
		cntrInfodHttpUrl := string(kvp.Value)

		// MAC
		kvp, _, err = kv.Get("containers/"+instanceId+"/macAdress", nil)
		if err != nil {
			fmt.Printf("Error while trying to read container registry: %s\n", err)
			return containers
		}
		macAddress := string(kvp.Value)

		// IP
		kvp, _, err = kv.Get("containers/"+instanceId+"/ipAdress", nil)
		if err != nil {
			fmt.Printf("Error while trying to read container registry: %s\n", err)
			return containers
		}
		ipAddress := string(kvp.Value)

		// Unix Epoch Timestamp
		kvp, _, err = kv.Get("containers/"+instanceId+"/unixEpochTimestamp", nil)
		if err != nil {
			fmt.Printf("Error while trying to read container registry: %s\n", err)
			return containers
		}
		ageInSeconds := DetermineAgeInSeconds(kvp.Value)

		// Hostname
		kvp, _, err = kv.Get("containers/"+instanceId+"/hostname", nil)
		if err != nil {
			fmt.Printf("Error while trying to read container registry: %s\n", err)
			return containers
		}
		hostname := string(kvp.Value)

		// HostHostname
		kvp, _, err = kv.Get("containers/"+instanceId+"/hostinfo/hostname", nil)
		if err != nil {
			fmt.Printf("Error while trying to read container registry: %s\n", err)
			return containers
		}
		hostHostname := string(kvp.Value)

		container := Container{Hostname: hostname, CntrInfodHttpUrl: cntrInfodHttpUrl, MacAddress: macAddress,
			IpAddress: ipAddress, HostHostname: hostHostname, AgeInSeconds: ageInSeconds}

		containers = append(containers, container)
		fmt.Printf("Successfully read info for instance ID %s from container registry.\n", instanceId)
	}

	return containers
}
