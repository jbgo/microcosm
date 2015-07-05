package main

// package proxy

import (
	"bytes"
	"fmt"
	"github.com/jbgo/mission_control/docker_client"
	// "io/ioutil"
	"log"
	"text/template"
)

func groupContainersByService(containers []*docker_client.Container) map[string][]*docker_client.Container {
	groups := make(map[string][]*docker_client.Container)

	for _, c := range containers {
		c.HostIP = c.Original.NetworkSettings.Ports["4567/tcp"][0].HostIP
		c.HostPort = c.Original.NetworkSettings.Ports["4567/tcp"][0].HostPort
		list, _ := groups[c.Labels["service"]]
		groups[c.Labels["service"]] = append(list, c)
	}

	return groups
}

func generateHaproxyConfig(serviceGroups map[string][]*docker_client.Container) bytes.Buffer {
	var buffer bytes.Buffer
	templatePath := "/go/src/app/haproxy.cfg.template"
	// templatePath := "/Users/jordan/projects/gocode/src/github.com/jbgo/mission_control/mc_proxy/haproxy.cfg.template"
	cfg := template.Must(template.ParseFiles(templatePath))
	cfg.ExecuteTemplate(&buffer, "haproxy.cfg.template", serviceGroups)
	return buffer
}

// func generateProxyConfigs(client *DockerClient, proxyContainer, webContainers []*Container) {
func main() {
	client, err := docker_client.New()
	if err != nil {
		log.Fatal(err)
	}

	webContainers, err := client.FindContainersWithLabel("service_type=web")
	if err != nil {
		log.Fatal(err)
	}

	serviceGroups := groupContainersByService(webContainers)
	haproxyConfig := generateHaproxyConfig(serviceGroups)
	fmt.Println(haproxyConfig.String())
	err = ioutil.WriteFile("/usr/local/etc/haproxy/haproxy.cfg", haproxyConfig.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}

	proxyContainer, err := client.FindContainerWithLabel("service=mc_proxy")
	if err != nil {
		log.Fatal(err)
	} else if proxyContainer == nil {
		log.Fatalf("[mc_proxy] could not find container labeled service=mc_proxy")
	}

	log.Println("[mc_proxy] restarting mc_proxy")
	err = client.RestartContainer(proxyContainer)
	if err != nil {
		log.Fatal(err)
	}
}
