package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/jbgo/mission_control/docker_client"
	"log"
	"strings"
)

type Action interface {
	Execute(container *docker_client.Container) error
}

type DockerRunAction struct {
	Image string
	Cmd   []string
}

func (action DockerRunAction) Execute(container *docker_client.Container) error {
	return nil
}

func findMatchingActions(eventType string, container *docker_client.Container) []Action {
	action := DockerRunAction{Image: "mc_proxy", Cmd: []string{"configure"}}
	return []Action{action}
}

func handleEvent(
	client *docker_client.DockerClient,
	event *docker.APIEvents) {
	container, err := client.FindContainer(event.ID)
	if err != nil {
		log.Printf("[mc_agent] [error] %v\n", err)
	}

	actions := findMatchingActions(event.Status, container)
	for _, action := range actions {
		switch a := action.(type) {
		case DockerRunAction:
			log.Printf("[mc_agent] [action] image: %s, cmd: %s", a.Image, strings.Join(a.Cmd, " "))
		}
		err = action.Execute(container)
		if err != nil {
			log.Printf("[mc_agent] [error] %v\n", err)
		}
	}
}

func main() {
	fmt.Println("[mc_agent] starting...")

	client, err := docker_client.New()
	if err != nil {
		log.Fatal(err)
	}

	listener := make(chan *docker.APIEvents, 10)
	err = client.Client.AddEventListener(listener)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[mc_agent] listening for events")

	for {
		select {
		case event := <-listener:
			fmt.Printf("[mc_agent] [%d:%7s] %s %s\n", event.Time, event.Status, event.ID[0:12], event.From)
			go handleEvent(&client, event)
		}
	}
}
