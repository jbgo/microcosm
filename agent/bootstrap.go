package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"time"
)

func logError(msg string, err error) error {
	fmt.Printf("[agent][ERROR] %s. Reason: %v\n", msg, err)
	return err
}

func bootstrap(client *docker.Client) error {
	dataContainerId, err := findOrCreateDataContainer(client)
	if err != nil {
		return logError("failed to find data container", err)
	}
	fmt.Println("AAA")

	opts := docker.ListContainersOptions{
		All:     true,
		Filters: filtersForService("proxy"),
	}
	containers, err := client.ListContainers(opts)
	if err != nil {
		return logError("failed to list proxy containers", err)
	}

	fmt.Println("BBB")

	// TODO pull nginx image before we can create this container

	var proxyContainerId string
	if len(containers) > 0 {
		proxyContainerId = containers[0].ID
		fmt.Printf("[agent] found microcosm-proxy container: %s\n", proxyContainerId)
	} else {
		proxyContainer, err := createProxyContainer(client, dataContainerId)
		if err != nil {
			return logError("failed to create proxy container", err)
		}
		proxyContainerId = proxyContainer.ID
		fmt.Printf("[agent] created microcosm-proxy container: %s\n", proxyContainerId)
	}

	fmt.Println("CCC")

	// TODO find configure image or build if not found

	// reconfigure nginx
	err = reconfigure(client, dataContainerId)
	if err != nil {
		return logError("failed to reconfigure nginx", err)
	}

	fmt.Println("DDD")

	proxyInfo, err := client.InspectContainer(proxyContainerId)
	if err != nil {
		return logError("failed to inspect proxy container", err)
	}

	if proxyInfo.State.Running {
		fmt.Println("[agent] restarting microcosm-proxy")
		return client.RestartContainer(proxyContainerId, 60)
	} else {
		fmt.Println("[agent] starting microcosm-proxy")
		return client.StartContainer(proxyContainerId, proxyInfo.HostConfig)
	}
}

func findOrCreateDataContainer(client *docker.Client) (string, error) {
	opts := docker.ListContainersOptions{
		All:     true,
		Filters: filtersForService("data-container"),
	}
	containers, err := client.ListContainers(opts)
	if err != nil {
		logError("[agent] failed to find data container", err)
		return "", err
	}

	var dataContainerId string
	if len(containers) > 0 {
		dataContainerId = containers[0].ID
		fmt.Printf("[agent] found microcosm-proxy-data container: %s\n", dataContainerId)
	} else {
		dataContainer, err := createDataContainer(client)
		if err != nil {
			logError("[agent] failed to create data container", err)
			return "", err
		}
		dataContainerId = dataContainer.ID
		fmt.Printf("[agent] created microcosm-proxy-data container: %s\n", dataContainerId)
	}

	return dataContainerId, nil
}

func createDataContainer(client *docker.Client) (*docker.Container, error) {
	opts := docker.CreateContainerOptions{
		Name: "microcosm-proxy-data",
		Config: &docker.Config{
			Image:   "nginx",
			Labels:  labelsForService("data-container"),
			Volumes: map[string]struct{}{"/etc/nginx": struct{}{}},
		},
	}
	return client.CreateContainer(opts)
}

func createProxyContainer(client *docker.Client, dataContainerId string) (*docker.Container, error) {
	opts := docker.CreateContainerOptions{
		Name: "microcosm-proxy",
		Config: &docker.Config{
			Image:  "nginx",
			Labels: labelsForService("proxy"),
		},
		HostConfig: &docker.HostConfig{
			NetworkMode: "host",
			VolumesFrom: []string{"microcosm-code", "microcosm-proxy-data"},
		},
	}
	return client.CreateContainer(opts)
}

func reconfigure(client *docker.Client, dataContainerId string) error {
	fmt.Println("begin reconfigure")
	opts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: "microcosm/proxy",
			Env:   []string{"DOCKER_HOST=unix:///var/run/docker.sock"},
		},
		HostConfig: &docker.HostConfig{
			Binds:       []string{"/var/run/docker.sock:/var/run/docker.sock"},
			VolumesFrom: []string{"microcosm-code", "microcosm-proxy-data"},
		},
	}

	fmt.Println("[agent] creating container to reconfigure nginx")
	container, err := client.CreateContainer(opts)
	if err != nil {
		return logError("failed to create reconfigure container", err)
	}

	fmt.Println("[agent] reconfiguring nginx")
	err = client.StartContainer(container.ID, container.HostConfig)
	if err != nil {
		// don't return for this one because we want to remove it
		logError("failed to start reconfigure container", err)
	}

	// TODO we should actually poll for status or use some other method to
	// actually verify the task completes before removing the container.
	time.Sleep(1 * time.Second)
	err2 := client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})
	if err != nil {
		logError("failed to remove reconfigure container", err2)
	}

	return err
}

func labelsForService(serviceType string) map[string]string {
	return map[string]string{
		"microcosm.service": "microcosm",
		"microcosm.type":    serviceType,
	}
}

func filtersForService(serviceType string) map[string][]string {
	return map[string][]string{
		"label": []string{"microcosm.service=microcosm", "microcosm.type=" + serviceType},
	}
}
