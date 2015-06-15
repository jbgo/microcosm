package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
)

var client docker.Client

type indexData struct {
	Containers []docker.APIContainers
}

func getClient() *docker.Client {
	endpoint := os.Getenv("DOCKER_HOST")
	path := os.Getenv("DOCKER_CERT_PATH")
	ca := fmt.Sprintf("%s/ca.pem", path)
	cert := fmt.Sprintf("%s/cert.pem", path)
	key := fmt.Sprintf("%s/key.pem", path)

	client, err := docker.NewTLSClient(endpoint, cert, key, ca)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

type ContainerData struct {
	State     string
	ID        string
	Name      string
	Image     string
	Command   string
	Container *docker.Container
}

type ContainerOverview struct {
	Running    []ContainerData
	NotRunning []ContainerData
}

type ByAge []ContainerData

func (s ByAge) Len() int {
	return len(s)
}

func (s ByAge) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByAge) Less(i, j int) bool {
	a, b := s[i].Container.State, s[j].Container.State
	if a.Running && b.Running {
		return a.StartedAt.Unix() < b.StartedAt.Unix()
	} else if a.Running && !b.Running {
		return true
	} else if !a.Running && b.Running {
		return false
	} else {
		return a.FinishedAt.Unix() > b.FinishedAt.Unix()
	}
}

type ByName []ContainerData

func (s ByName) Len() int {
	return len(s)
}

func (s ByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func (d *ContainerData) Fill(c *docker.Container) {
	d.Container = c
	d.Name = c.Name
	d.ID = c.ID

	if c.State.Running {
		d.State = "running"
	} else if c.State.Paused {
		d.State = "paused"
	} else if c.State.Restarting {
		d.State = "restarting"
	} else if c.State.OOMKilled {
		d.State = "oom_killed"
	} else {
		d.State = "stopped"
	}
}

var templateFuncs = template.FuncMap{
	"join":     strings.Join,
	"noslash":  func(s string) string { return s[1:] },
	"truncate": func(s string, length int) string { return s[0:length] },
}

func fillContainerData(client *docker.Client, containers []docker.APIContainers) []ContainerData {
	var data []ContainerData
	for _, c := range containers {
		inspect, err := client.InspectContainer(c.ID)
		if err != nil {
			// TODO could we show an error state here instead?
			log.Fatal(err)
		}
		d := &ContainerData{}
		d.Fill(inspect)
		data = append(data, *d)
	}
	return data
}

func runningContainers(cd []ContainerData) []ContainerData {
	var running []ContainerData
	for _, d := range cd {
		if d.State == "running" {
			running = append(running, d)
		}
	}
	return running
}

func notRunningContainers(cd []ContainerData) []ContainerData {
	var notRunning []ContainerData
	for _, d := range cd {
		if d.State != "running" {
			notRunning = append(notRunning, d)
		}
	}
	return notRunning
}

func showContainersHandler(w http.ResponseWriter, r *http.Request) {
	client := getClient()

	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		fmt.Fprintf(w, "TODO: unkown error")
	} else {
		cd := fillContainerData(client, containers)
		data := ContainerOverview{
			Running:    runningContainers(cd),
			NotRunning: notRunningContainers(cd),
		}
		sort.Sort(ByName(data.Running))
		sort.Sort(ByName(data.NotRunning))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t := template.Must(template.New("").Funcs(templateFuncs).ParseFiles("views/index.html"))
		t.ExecuteTemplate(w, "index.html", data)
	}
}

func startWebserver() {
	http.HandleFunc("/", showContainersHandler)
	_, filename, _, _ := runtime.Caller(0)

	assetsPath := path.Join(path.Dir(filename), "assets")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsPath))))

	fmt.Println("Listening on port 8001")
	http.ListenAndServe(":8081", nil)
}

func main() {
	go startWebserver()

	client := getClient()

	env, err := client.Version()
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range env.Map() {
		fmt.Printf("%s => %v\n", k, v)
	}

	listener := make(chan *docker.APIEvents, 10)
	err = client.AddEventListener(listener)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case msg := <-listener:
			fmt.Printf("[event:%d] status: %s, image: %s, cid: %s\n", msg.Time, msg.Status, msg.From, msg.ID)
		}
	}
}
