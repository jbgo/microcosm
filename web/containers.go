package web

import (
	"github.com/gorilla/mux"
	"github.com/jbgo/mission_control/docker_client"
	"log"
	"net/http"
)

type ContainerOverview struct {
	Running    docker_client.Containers
	NotRunning docker_client.Containers
}

func (app WebApp) ListContainers(w http.ResponseWriter, r *http.Request) {
	client, err := docker_client.New()
	if err != nil {
		log.Fatal(err)
	}

	containers, err := client.GetContainers()
	if err != nil {
		log.Fatal(err)
	}

	overview := ContainerOverview{
		Running:    containers.Running(),
		NotRunning: containers.NotRunning(),
	}

	app.RenderHTML(w, "main", "containers/index", overview)
}

func (app WebApp) ShowContainer(w http.ResponseWriter, r *http.Request) {
	client, err := docker_client.New()
	if err != nil {
		log.Fatal(err)
	}

	containerID := mux.Vars(r)["id"]
	container, err := client.FindContainer(containerID)
	if err != nil {
		log.Fatal(err)
	}

	app.RenderHTML(w, "main", "containers/show", container)
}
