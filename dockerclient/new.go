package dockerclient

import (
	docker "github.com/fsouza/go-dockerclient"
	"os"
	"path"
)

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
