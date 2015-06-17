package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Serve(address string) {
	app := BuildApp()

	router := mux.NewRouter().StrictSlash(false)

	router.HandleFunc("/", app.Home)

	containers := router.Path("/containers").Subrouter()
	containers.Methods("GET").HandlerFunc(app.ListContainers)

	container := router.Path("/containers/{id}").Subrouter()
	container.Methods("GET").HandlerFunc(app.ShowContainer)

	repos := router.Path("/repos").Subrouter()
	repos.Methods("GET").HandlerFunc(app.ListRepos)

	services := router.Path("/services").Subrouter()
	services.Methods("GET").HandlerFunc(app.ListServices)

	router.PathPrefix("/assets/").Handler(app.AssetsHandler())

	fmt.Println("Listening on ", address)
	app.Log.Fatal(http.ListenAndServe(address, router))
}
