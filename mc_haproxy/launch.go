package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	docker "github.com/samalba/dockerclient"
	"io/ioutil"
	"log"
	"path"
)

func main() {
	// establish TLS config
	certPath := "/Users/jordan/.docker/machine/machines/dev"

	cert, err := tls.LoadX509KeyPair(path.Join(certPath, "cert.pem"), path.Join(certPath, "key.pem"))
	if err != nil {
		log.Fatal(err)
	}

	caCert, err := ioutil.ReadFile(path.Join(certPath, "ca.pem"))
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlc := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		RootCAs:            caCertPool,
	}

	// create docker client
	client, err := docker.NewDockerClient("tcp://192.168.99.100:2376", tlc)
	if err != nil {
		log.Fatal(err)
	}

	containerConfig := &docker.ContainerConfig{
		Image:   "haproxy:1.5",
		Labels:  map[string]string{"service": "mc_proxy", "service_type": "proxy"},
		Volumes: map[string]struct{}{"/usr/local/etc/haproxy": struct{}{}},
		HostConfig: docker.HostConfig{
			NetworkMode: "host",
		},
	}

	// create container
	containerId, err := client.CreateContainer(containerConfig, "mc_haproxy_test")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("created container: %s", containerId)

	// start container
}
