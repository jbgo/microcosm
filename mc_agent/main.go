package main

import (
	"fmt"
	"github.com/jbgo/mission_control/dockerclient"
	"log"
)

func main() {
	fmt.Println("[mc_agent] starting...")

	client, err := dockerclient.New()
	if err != nil {
		log.Fatal(err)
	}

	err = bootstrap(client)
	if err != nil {
		log.Fatal(err)
	}
}
