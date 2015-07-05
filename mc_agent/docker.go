package main

import (
	"crypto/tls"
	"crypto/x509"
	docker "github.com/samalba/dockerclient"
	"io/ioutil"
	"os"
	"path"
)

func newDockerClient() (*docker.DockerClient, error) {
	endpoint := os.Getenv("DOCKER_HOST")
	tlsVerify := os.Getenv("DOCKER_TLS_VERIFY")

	if tlsVerify == "1" {
		return newTLSDockerClient(endpoint)
	} else {
		return docker.NewDockerClient(endpoint, nil)
	}
}

func newTLSDockerClient(endpoint string) (*docker.DockerClient, error) {
	certPath := os.Getenv("DOCKER_CERT_PATH")

	certFile := path.Join(certPath, "cert.pem")
	keyFile := path.Join(certPath, "key.pem")
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(path.Join(certPath, "ca.pem"))
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlc := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		RootCAs:            caCertPool,
	}

	return docker.NewDockerClient(endpoint, tlc)
}
