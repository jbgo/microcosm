package main

import (
	"fmt"
	docker "github.com/jbgo/microcosm/dockerclient"
	"log"
)

func main() {
	fmt.Println("[agent] starting...")

	client, err := docker.New()
	if err != nil {
		log.Fatal(err)
	}

	agent := New(client)
	agent.Listen()
}
