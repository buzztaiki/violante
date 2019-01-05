package main

import (
	"flag"
	"log"

	"github.com/buzztaiki/violante"
)

func main() {
	addr := flag.String("addr", ":51289", "server listen address")
	flag.Parse()

	client := violante.NewClient(*addr)
	if err := client.Add(flag.Args()); err != nil {
		log.Fatal(err)
	}
}
