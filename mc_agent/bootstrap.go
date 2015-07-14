package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"time"
)

func bootstrap(client *docker.Client) error {
	dataContainerId, err := findOrCreateDataContainer(client)
	if err != nil {
		return err
	}

	opts := docker.ListContainersOptions{
		All:     true,
		Filters: map[string][]string{"label": []string{"service=mc_proxy", "service_type=reverse-proxy"}},
	}
	containers, err := client.ListContainers(opts)
	if err != nil {
		return err
	}

	var proxyContainerId string
	if len(containers) > 0 {
		proxyContainerId = containers[0].ID
		fmt.Printf("[mc_agent] found mc_proxy container: %s\n", proxyContainerId)
	} else {
		proxyContainer, err := createProxyContainer(client, dataContainerId)
		if err != nil {
			return err
		}
		proxyContainerId = proxyContainer.ID
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
		fmt.Println("[mc_agent] restarting mc_proxy")
		return client.RestartContainer(proxyContainerId, 60)
	} else {
		fmt.Println("[mc_agent] starting mc_proxy")
		return client.StartContainer(proxyContainerId, proxyInfo.HostConfig)
	}
}

func findOrCreateDataContainer(client *docker.Client) (string, error) {
	opts := docker.ListContainersOptions{
		All:     true,
		Filters: map[string][]string{"label": []string{"service=mc_proxy", "service_type=data-container"}},
	}
	containers, err := client.ListContainers(opts)
	if err != nil {
		return "", err
	}

	var dataContainerId string
	if len(containers) > 0 {
		dataContainerId = containers[0].ID
		fmt.Printf("[mc_agent] found mc_proxy_data container: %s\n", dataContainerId)
	} else {
		dataContainer, err := createDataContainer(client)
		if err != nil {
			return "", err
		}
		dataContainerId = dataContainer.ID
		fmt.Printf("[mc_agent] created mc_proxy_data container: %s\n", dataContainerId)
	}

	return dataContainerId, nil
}

func createDataContainer(client *docker.Client) (*docker.Container, error) {
	opts := docker.CreateContainerOptions{
		Name: "mc_proxy_data",
		Config: &docker.Config{
			Image:   "nginx",
			Labels:  map[string]string{"service": "mc_proxy", "service_type": "data-container"},
			Volumes: map[string]struct{}{"/etc/nginx": struct{}{}},
		},
	}
	return client.CreateContainer(opts)
}

func createProxyContainer(client *docker.Client, dataContainerId string) (*docker.Container, error) {
	opts := docker.CreateContainerOptions{
		Name: "mc_proxy",
		Config: &docker.Config{
			Image:  "nginx",
			Labels: map[string]string{"service": "mc_proxy", "service_type": "reverse-proxy"},
		},
		HostConfig: &docker.HostConfig{
			NetworkMode: "host",
			VolumesFrom: []string{dataContainerId},
		},
	}
	return client.CreateContainer(opts)
}

// docker run --rm -it -e DOCKER_HOST=unix:///var/run/docker.sock -v /var/run/docker.sock:/var/run/docker.sock --volumes-from=mc_proxy_data mc_proxy
func reconfigure(client *docker.Client, dataContainerId string) error {
	opts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: "mc_reconfigure_proxy",
			Env:   []string{"DOCKER_HOST=unix:///var/run/docker.sock"},
		},
		HostConfig: &docker.HostConfig{
			Binds:       []string{"/var/run/docker.sock:/var/run/docker.sock"},
			VolumesFrom: []string{"mc_proxy_data"},
		},
	}

	fmt.Println("[mc_agent] creating container to reconfigure nginx")
	container, err := client.CreateContainer(opts)
	if err != nil {
		return err
	}

	fmt.Println("[mc_agent] reconfiguring nginx")
	err = client.StartContainer(container.ID, container.HostConfig)

	// TODO we should actually poll for status or use some other method to
	// actually verify the task completes before removing the container.
	time.Sleep(1 * time.Second)
	client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})

	return err
}
