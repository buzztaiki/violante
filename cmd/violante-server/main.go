package main

import (
	"flag"
	"log"
	"os"

	"github.com/buzztaiki/violante"
)

func main() {
	addr := flag.String("addr", ":51289", "listen address")
	flag.Parse()

	notifier := violante.NewSlackNotifier(os.Getenv("SLACK_WEBHOOK_URL"), os.Getenv("SLACK_CHANNEL"))

	det, err := violante.NewDetector(os.Getenv("VT_API_KEY"), notifier)
	if err != nil {
		log.Fatal(err)
	}

	go det.Start()
	defer det.Shutdown()

	server := violante.NewServer(*addr, det)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
