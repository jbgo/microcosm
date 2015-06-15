package main

import (
	"fmt"
)

var a = "foo"

func main() {
	var b = "foo"
	if a == b {
		fmt.Println("EQUAL")
	} else {
		fmt.Println("NOT EQUAL")
	}
}
