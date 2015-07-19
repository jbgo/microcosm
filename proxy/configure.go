package main

import (
	"bytes"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/jbgo/microcosm/dockerclient"
	"io/ioutil"
	"log"
	"path"
	"runtime"
	"text/template"
)

type ContainerInfo struct {
	HostIP   string
	HostPort int64
}

func groupContainersByService(client *docker.Client, containers []docker.APIContainers) map[string][]*ContainerInfo {
	groups := make(map[string][]*ContainerInfo)

	for _, c := range containers {
		port := c.Ports[0]
		info := &ContainerInfo{
			HostIP:   port.IP,
			HostPort: port.PublicPort,
		}

		inspect, _ := client.InspectContainer(c.ID)
		service := inspect.Config.Labels["microcosm.service"]

		list, _ := groups[service]
		groups[service] = append(list, info)
	}

	return groups
}

func findTemplatePath(template string) string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filename), template)
}

func generateNginxConfig(serviceGroups map[string][]*ContainerInfo) bytes.Buffer {
	var buffer bytes.Buffer
	templatePath := findTemplatePath("nginx.conf.template")
	cfg := template.Must(template.ParseFiles(templatePath))
	cfg.ExecuteTemplate(&buffer, "nginx.conf.template", serviceGroups)
	return buffer
}

func main() {
	client, err := dockerclient.New()
	if err != nil {
		log.Fatal(err)
	}

	webContainers, err := client.ListContainers(docker.ListContainersOptions{
		All:     false,
		Filters: map[string][]string{"label": []string{"microcosm.type=web"}},
	})
	if err != nil {
		log.Fatal(err)
	}

	serviceGroups := groupContainersByService(client, webContainers)
	nginxConfig := generateNginxConfig(serviceGroups)
	fmt.Println(nginxConfig.String())
	err = ioutil.WriteFile("/etc/nginx/conf.d/microcosm.conf", nginxConfig.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
