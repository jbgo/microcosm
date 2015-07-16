package admin

import (
	"github.com/gorilla/mux"
	"github.com/jbgo/microcosm/dockerclient"
	"log"
	"net/http"
)

type ContainerOverview struct {
	Running    dockerclient.Containers
	NotRunning dockerclient.Containers
}

func (app WebApp) ListContainers(w http.ResponseWriter, r *http.Request) {
	client, err := dockerclient.NewSimpleClient()
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
	client, err := dockerclient.NewSimpleClient()
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
