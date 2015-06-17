package docker_client

import (
	"github.com/fsouza/go-dockerclient"
	"log"
	"sort"
)

type Container struct {
	State    string
	ID       string
	Name     string
	Image    string
	Command  string
	Original *docker.Container
}

type Containers []*Container

func (c *Container) init(d *docker.Container) {
	c.Original = d
	c.Name = d.Name
	c.ID = d.ID

	if d.State.Running {
		c.State = "running"
	} else if d.State.Paused {
		c.State = "paused"
	} else if d.State.Restarting {
		c.State = "restarting"
	} else if d.State.OOMKilled {
		c.State = "oom_killed"
	} else {
		c.State = "stopped"
	}
}

func (d DockerClient) FindContainer(containerID string) (*Container, error) {
	c := &Container{}
	inspect, err := d.client.InspectContainer(containerID)
	if err != nil {
		return c, err
	}
	c.init(inspect)
	return c, err
}

func (d DockerClient) GetContainers() (Containers, error) {
	var containers Containers
	results, err := d.client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		log.Fatal(err)
		return containers, err
	}

	for _, r := range results {
		c, err := d.FindContainer(r.ID)
		if err != nil {
			return containers, err
		}
		containers = append(containers, c)
	}

	return containers, err
}

func (containers Containers) Running() Containers {
	var running Containers

	for _, c := range containers {
		if c.State == "running" {
			running = append(running, c)
		}
	}

	sort.Sort(ByName(running))
	return running
}

func (containers Containers) NotRunning() Containers {
	var notRunning Containers

	for _, c := range containers {
		if c.State != "running" {
			notRunning = append(notRunning, c)
		}
	}

	sort.Sort(ByAge(notRunning))
	return notRunning
}
