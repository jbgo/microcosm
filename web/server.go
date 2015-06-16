package web

import (
	"fmt"
	"log"
	"net/http"
)

func Serve(address string) {
	app := BuildApp()

	http.HandleFunc("/", app.Home)
	http.Handle("/assets/", app.AssetsHandler())

	fmt.Println("Listening on ", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
