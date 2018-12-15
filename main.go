package main

import (
	"log"
)

func main() {
	err := example()
	if err != nil {
		log.Fatalf("%#v", err)
	}
}
