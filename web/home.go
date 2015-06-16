package web

import (
	"github.com/jbgo/mission_control/docker_client"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type ContainerOverview struct {
	Running    docker_client.Containers
	NotRunning docker_client.Containers
}

var templateFuncs = template.FuncMap{
	"join":     strings.Join,
	"noslash":  func(s string) string { return s[1:] },
	"truncate": func(s string, length int) string { return s[0:length] },
}

func (app WebApp) Home(w http.ResponseWriter, r *http.Request) {
	client, err := docker_client.New()
	if err != nil {
		log.Fatal(err)
	}

	containers, err := client.GetContainers()
	if err != nil {
		log.Fatal(err)
	}

	data := ContainerOverview{
		Running:    containers.Running(),
		NotRunning: containers.NotRunning(),
	}

	viewFiles := path.Join(app.Root, "views/index.html")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.Must(template.New("").Funcs(templateFuncs).ParseFiles(viewFiles))
	t.ExecuteTemplate(w, "index.html", data)
}
