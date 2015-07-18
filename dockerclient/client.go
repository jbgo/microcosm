/*
Package dockerclient provides high-level functions for working with
docker containers and images.

It provides convenient helper functions for performing common series
of docker remote API calls via the go-dockerclient package. It is not
intended to abstract away the API. Client code can and should use the
go-dockerclient types and functions directly.

Docker Remote API Documentation
https://docs.docker.com/reference/api/docker_remote_api_v1.19/

Go-dockerclient Documentation
https://godoc.org/github.com/fsouza/go-dockerclient
*/
package dockerclient

import (
	docker "github.com/fsouza/go-dockerclient"
	"os"
	"path"
)

type DockerClient struct {
	Client *docker.Client
}

func New() (*docker.Client, error) {
	endpoint := os.Getenv("DOCKER_HOST")
	tlsVerify := os.Getenv("DOCKER_TLS_VERIFY")

	if tlsVerify == "1" {
		return newTLSDockerClient(endpoint)
	} else {
		return docker.NewClient(endpoint)
	}
}

func newTLSDockerClient(endpoint string) (*docker.Client, error) {
	certPath := os.Getenv("DOCKER_CERT_PATH")
	certFile := path.Join(certPath, "cert.pem")
	keyFile := path.Join(certPath, "key.pem")
	caFile := path.Join(certPath, "ca.pem")

	return docker.NewTLSClient(endpoint, certFile, keyFile, caFile)
}

func NewSimpleClient() (DockerClient, error) {
	client, err := New()
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
