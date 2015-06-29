package runtime

import (
	"fmt"
	"github.com/jbgo/mission_control/docker_client"
)

type ServiceConf struct {
	Cmd    []string
	Image  string
	Labels map[string]string
}

func Bootstrap() error {
	fmt.Println("[mc_admin] checking for existence of required services")

	client, err := docker_client.New()
	if err != nil {
		return err
	}

	storageConf := ServiceConf{
		Cmd:    []string{"sleep", "90"},
		Image:  "debian:wheezy",
		Labels: map[string]string{"service": "mc_storage"},
	}
	err = startService(&client, &storageConf)
	if err != nil {
		return err
	}

	agentConf := ServiceConf{
		Cmd:    []string{"sleep", "90"},
		Image:  "debian:wheezy",
		Labels: map[string]string{"service": "mc_agent"},
	}
	err = startService(&client, &agentConf)
	if err != nil {
		return err
	}

	proxyConf := ServiceConf{
		Cmd:    []string{"sleep", "90"},
		Image:  "debian:wheezy",
		Labels: map[string]string{"service": "mc_proxy"},
	}
	err = startService(&client, &proxyConf)
	if err != nil {
		return err
	}

	return nil
}

func startService(client *docker_client.DockerClient, conf *ServiceConf) error {
	container, err := client.FindContainerWithLabel("service=" + conf.Labels["service"])
	if err != nil {
		return err
	}

	if container == nil {
		fmt.Printf("[mc_admin] %s service not found\n", conf.Labels["service"])
		fmt.Printf("[mc_admin] creating %s service\n", conf.Labels["service"])
		container, err = client.CreateContainer(conf.Image, conf.Labels, conf.Cmd)
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
		err = client.StartContainer(container)
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
