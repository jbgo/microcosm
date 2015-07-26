/*
Package dockerclient provides a high-level interface for docker.

Package developers should be familiar with the docker remote API.
https://docs.docker.com/reference/api/docker_remote_api_v1.19/
*/
package dockerclient

//go:generate mockgen -destination mock_dockerclient.go github.com/jbgo/microcosm/dockerclient DockerClient

import (
	docker "github.com/fsouza/go-dockerclient"
	"os"
	"path"
)

type goDockerClient struct {
	Client *docker.Client
}

func New() (DockerClient, error) {
	endpoint := os.Getenv("DOCKER_HOST")
	tlsVerify := os.Getenv("DOCKER_TLS_VERIFY")

	var client *docker.Client
	var err error

	if tlsVerify == "1" {
		client, err = newTLSDockerClient(endpoint)
	} else {
		client, err = docker.NewClient(endpoint)
	}

	return DockerClient(&goDockerClient{Client: client}), err
}

func newTLSDockerClient(endpoint string) (*docker.Client, error) {
	certPath := os.Getenv("DOCKER_CERT_PATH")
	certFile := path.Join(certPath, "cert.pem")
	keyFile := path.Join(certPath, "key.pem")
	caFile := path.Join(certPath, "ca.pem")

	return docker.NewTLSClient(endpoint, certFile, keyFile, caFile)
}

func (c *Container) initContainer(d *docker.Container) {
	c.ID = d.ID
	c.Labels = d.Config.Labels
}

func (d *goDockerClient) InspectContainer(id string) (*Container, error) {
	c := &Container{}
	inspect, err := d.Client.InspectContainer(id)
	if err != nil {
		return nil, err
	}
	c.initContainer(inspect)
	return c, err
}

func initEvent(dockerEvent *docker.APIEvents) *Event {
	return &Event{
		ContainerID: dockerEvent.ID,
		Status:      dockerEvent.Status,
		ImageID:     dockerEvent.From,
		Timestamp:   dockerEvent.Time,
	}
}

func proxyEvents(listener chan<- *Event, privateListener <-chan *docker.APIEvents) {
	for {
		select {
		case dockerEvent := <-privateListener:
			listener <- initEvent(dockerEvent)
		}
	}
}

func (d *goDockerClient) AddEventListener(listener chan<- *Event) error {
	privateListener := make(chan *docker.APIEvents, 10)
	go proxyEvents(listener, privateListener)
	return d.Client.AddEventListener(privateListener)
}
