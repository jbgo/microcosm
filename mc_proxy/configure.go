package main

import (
	"bytes"
	"fmt"
	"github.com/jbgo/mission_control/docker_client"
	"io/ioutil"
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

func generateNginxConfig(serviceGroups map[string][]*docker_client.Container) bytes.Buffer {
	var buffer bytes.Buffer
	templatePath := "/go/src/app/nginx.conf.template"
	cfg := template.Must(template.ParseFiles(templatePath))
	cfg.ExecuteTemplate(&buffer, "nginx.conf.template", serviceGroups)
	return buffer
}

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
	nginxConfig := generateNginxConfig(serviceGroups)
	fmt.Println(nginxConfig.String())
	err = ioutil.WriteFile("/etc/nginx/conf.d/microcosm.conf", nginxConfig.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
