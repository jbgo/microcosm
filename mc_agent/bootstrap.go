package main

import (
	"fmt"
	docker "github.com/samalba/dockerclient"
)

func bootstrap(client *docker.DockerClient) error {
	dataContainerId, err := createDataContainer(client, "mc_proxy_data", "/usr/local/etc/haproxy")
	if err != nil {
		return err
	}
	fmt.Printf("[mc_agent] created mc_proxy_data data container: %s\n", dataContainerId)

	// mcProxy, err := createMcProxyContainer(mcProxyData)
	// if err != nil { return err }

	// configureProxy(mcProxy)
	// if err != nil { return err }
	return nil
}

func createDataContainer(client *docker.DockerClient, name, volume string) (string, error) {
	containerConfig := &docker.ContainerConfig{
		Image:   "haproxy:1.5",
		Labels:  map[string]string{"service": "mc_proxy", "service_type": "persistence"},
		Volumes: map[string]struct{}{volume: struct{}{}},
	}
	return client.CreateContainer(containerConfig, name)
}
