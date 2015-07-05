package runtime

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/jbgo/mission_control/docker_client"
)

func Bootstrap() error {
	fmt.Println("[mc_admin] checking for existence of required services")

	client, err := docker_client.New()
	if err != nil {
		return err
	}

	storageConf := docker_client.ContainerConfig{
		Cmd:    []string{"sleep", "90"},
		Image:  "debian:wheezy",
		Labels: map[string]string{"service": "mc_storage", "service_type": "persistence"},
	}
	err = startService(&client, &storageConf)
	if err != nil {
		return err
	}

	agentConf := docker_client.ContainerConfig{
		Cmd:    []string{"sleep", "90"},
		Image:  "debian:wheezy",
		Labels: map[string]string{"service": "mc_agent", "service_type": "daemon"},
	}
	err = startService(&client, &agentConf)
	if err != nil {
		return err
	}

	// TODO empty but working config file
	// TODO need to add options: --net=host -p 80:80
	proxyConf := docker_client.ContainerConfig{
		Image:        "haproxy:1.5",
		Labels:       map[string]string{"service": "mc_proxy", "service_type": "proxy"},
		Volumes:      map[string]struct{}{"/usr/local/etc/haproxy": struct{}{}},
		NetworkMode:  "host",
		PortSpecs:    []string{":80:80"},
		Name:         "mc_haproxy",
		PortBindings: map[docker.Port][]docker.PortBinding{"80/tcp": []docker.PortBinding{docker.PortBinding{HostPort: "80/tcp"}}},
	}
	err = startService(&client, &proxyConf)
	if err != nil {
		return err
	}

	return nil
}

func startService(client *docker_client.DockerClient, conf *docker_client.ContainerConfig) error {
	container, err := client.FindContainerWithLabel("service=" + conf.Labels["service"])
	if err != nil {
		return err
	}

	if container == nil {
		fmt.Printf("[mc_admin] %s service not found\n", conf.Labels["service"])
		fmt.Printf("[mc_admin] creating %s service\n", conf.Labels["service"])
		container, err = client.CreateContainer(conf)
		if err != nil {
			return err
		} else {
			fmt.Printf("[mc_admin] created %s service\n", conf.Labels["service"])
		}
	} else {
		fmt.Printf("[mc_admin] %s service already created\n", conf.Labels["service"])
	}

	if container.State != "running" {
		fmt.Printf("[mc_admin] %s service not started\n", conf.Labels["service"])
		err = client.StartContainer(container, conf)
		if err != nil {
			return err
		} else {
			fmt.Printf("[mc_admin] started %s service\n", conf.Labels["service"])
		}
	} else {
		fmt.Printf("[mc_admin] %s service already started\n", conf.Labels["service"])
	}

	return nil
}
