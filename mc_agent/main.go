package main

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/jbgo/mission_control/dockerclient"
	"log"
)

func handleEvent(msg *docker.APIEvents) {
	client, err := dockerclient.New()
	if err != nil {
		log.Fatal(err)
		return
	}

	container, err := client.InspectContainer(msg.ID)
	if err != nil {
		fmt.Errorf("[mc_agent] failed to inspect container %s. %v\n", msg.ID, err)
		return
	}

	fmt.Printf("[mc_agent][%s] %s %s\n", msg.Status, container.Name, msg.ID[0:12])

	reconfigureEvents := map[string]bool{
		"start":   true,
		"stop":    true,
		"kill":    true,
		"pause":   true,
		"unpause": true,
	}

	if container.Config.Labels["service_type"] == "web" && reconfigureEvents[msg.Status] {
		fmt.Printf("[mc_agent] matched event status: %s, label: service_type=web\n", msg.Status)

		dataContainerId, err := findOrCreateDataContainer(client)
		if err != nil {
			fmt.Errorf("[mc_agent] could not find data container %v\n", err)
			return
		}
		// TODO find or build reconfigure image
		fmt.Printf("           dataContainerId: %s\n", dataContainerId)
		err = reconfigure(client, dataContainerId)
		if err != nil {
			fmt.Errorf("[mc_agent] could not reconfigure data container %v\n", err)
			return
		}
		// TODO find and start/restart mc_proxy
	}
}

func main() {
	fmt.Println("[mc_agent] starting...")

	client, err := dockerclient.New()
	if err != nil {
		log.Fatal(err)
	}

	err = bootstrap(client)
	if err != nil {
		log.Fatal(err)
	}

	listener := make(chan *docker.APIEvents, 10)
	err = client.AddEventListener(listener)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case msg := <-listener:
			go handleEvent(msg)
		}
	}
}
