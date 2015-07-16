package main

import (
	"github.com/jbgo/microcosm/admin"
	"log"
)

func main() {
	admin.Serve(":8001")
}
