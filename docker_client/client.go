package docker_client

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"os"
)

type DockerClient struct {
	client *docker.Client
}

func New() (DockerClient, error) {
	endpoint := os.Getenv("DOCKER_HOST")
	path := os.Getenv("DOCKER_CERT_PATH")
	ca := fmt.Sprintf("%s/ca.pem", path)
	cert := fmt.Sprintf("%s/cert.pem", path)
	key := fmt.Sprintf("%s/key.pem", path)

	client, err := docker.NewTLSClient(endpoint, cert, key, ca)
	return DockerClient{client: client}, err
}

func Version(d DockerClient) (map[string]string, error) {
	env, err := d.client.Version()
	if err != nil {
		return nil, err
	} else {
		return env.Map(), nil
	}
}
