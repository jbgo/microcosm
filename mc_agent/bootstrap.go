package main

import (
	"fmt"
	docker "github.com/samalba/dockerclient"
	"net/url"
	"time"
)

func bootstrap(client *docker.DockerClient) error {
	filters := url.QueryEscape(`{"label":["service=mc_proxy", "service_type=data-container"]}`)
	containers, err := client.ListContainers(true, false, filters)
	if err != nil {
		return err
	}

	var dataContainerId string
	if len(containers) > 0 {
		dataContainerId = containers[0].Id
		fmt.Printf("[mc_agent] found mc_proxy_data data container: %s\n", dataContainerId)
	} else {
		dataContainerId, err = createDataContainer(client)
		if err != nil {
			return err
		}
		fmt.Printf("[mc_agent] created mc_proxy_data data container: %s\n", dataContainerId)
	}

	filters = url.QueryEscape(`{"label":["service=mc_proxy", "service_type=reverse-proxy"]}`)
	containers, err = client.ListContainers(true, false, filters)
	if err != nil {
		return err
	}

	proxyHostConfig := docker.HostConfig{
		NetworkMode: "host",
		VolumesFrom: []string{dataContainerId},
	}

	var proxyContainerId string
	if len(containers) > 0 {
		proxyContainerId = containers[0].Id
		fmt.Printf("[mc_agent] found mc_proxy container: %s\n", proxyContainerId)
	} else {
		proxyContainerId, err = createProxyContainer(client, dataContainerId, proxyHostConfig)
		if err != nil {
			return err
		}
		fmt.Printf("[mc_agent] created mc_proxy container: %s\n", proxyContainerId)
	}

	// TODO find configure image or build if not found

	// reconfigure nginx
	err = reconfigure(client, dataContainerId)
	if err != nil {
		return err
	}

	proxyInfo, err := client.InspectContainer(proxyContainerId)
	if err != nil {
		return err
	}

	if proxyInfo.State.Running {
		return client.RestartContainer(proxyContainerId, 60)
	} else {
		return client.StartContainer(proxyContainerId, &proxyHostConfig)
	}
}

func createDataContainer(client *docker.DockerClient) (string, error) {
	containerConfig := &docker.ContainerConfig{
		Image:   "nginx",
		Labels:  map[string]string{"service": "mc_proxy", "service_type": "data-container"},
		Volumes: map[string]struct{}{"/etc/nginx": struct{}{}},
	}
	return client.CreateContainer(containerConfig, "mc_proxy_data")
}

func createProxyContainer(client *docker.DockerClient, dataContainerId string, hostConfig docker.HostConfig) (string, error) {
	containerConfig := &docker.ContainerConfig{
		Image:      "nginx",
		Labels:     map[string]string{"service": "mc_proxy", "service_type": "reverse-proxy"},
		HostConfig: hostConfig,
	}
	return client.CreateContainer(containerConfig, "mc_proxy")
}

// docker run --rm -it -e DOCKER_HOST=unix:///var/run/docker.sock -v /var/run/docker.sock:/var/run/docker.sock --volumes-from=mc_proxy_data mc_proxy
func reconfigure(client *docker.DockerClient, dataContainerId string) error {
	hostConfig := docker.HostConfig{
		Binds:       []string{"/var/run/docker.sock:/var/run/docker.sock"},
		VolumesFrom: []string{"mc_proxy_data"},
	}

	containerConfig := &docker.ContainerConfig{
		Image:      "mc_reconfigure_proxy",
		Env:        []string{"DOCKER_HOST=unix:///var/run/docker.sock"},
		HostConfig: hostConfig,
	}

	fmt.Println("[mc_agent] creating container to reconfigure nginx")
	containerId, err := client.CreateContainer(containerConfig, "tmp_mc_proxy")
	if err != nil {
		return err
	}

	fmt.Println("[mc_agent] reconfiguring nginx")
	err = client.StartContainer(containerId, &hostConfig)

	// wait for reconfigure to complete before removing container
	time.Sleep(1 * time.Second)
	client.RemoveContainer(containerId, false, false)

	return err
}
