package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("[mc_agent] starting...")

	client, err := newDockerClient()
	if err != nil {
		log.Fatal(err)
	}

	err = bootstrap(client)
	if err != nil {
		log.Fatal(err)
	}
}
