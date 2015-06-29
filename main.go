package main

import (
	"github.com/jbgo/mission_control/runtime"
	"github.com/jbgo/mission_control/web"
	"log"
)

func main() {
	err := runtime.Bootstrap()
	if err != nil {
		log.Fatal(err)
	}

	web.Serve(":8001")
}
