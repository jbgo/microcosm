package docker_client

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"os"
	"strings"
)

type DockerClient struct {
	Client *docker.Client
}

func New() (DockerClient, error) {
	var client *docker.Client
	var err error

	endpoint := os.Getenv("DOCKER_HOST")

	if strings.HasPrefix(endpoint, "unix") {
		client, err = docker.NewClient(endpoint)
	} else {
		path := os.Getenv("DOCKER_CERT_PATH")
		ca := fmt.Sprintf("%s/ca.pem", path)
		cert := fmt.Sprintf("%s/cert.pem", path)
		key := fmt.Sprintf("%s/key.pem", path)

		client, err = docker.NewTLSClient(endpoint, cert, key, ca)
	}

	return DockerClient{Client: client}, err
}

func Version(d DockerClient) (map[string]string, error) {
	env, err := d.Client.Version()
	if err != nil {
		return nil, err
	} else {
		return env.Map(), nil
	}
}
