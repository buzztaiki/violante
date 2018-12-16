package main

import (
	"flag"
	"log"
	"os"

	"github.com/williballenthin/govt"
)

func main() {
	addr := flag.String("addr", ":51289", "listen address")
	serverMode := flag.Bool("server", false, "run a server instead of client")
	flag.Parse()

	if *serverMode {
		vtClient, err := govt.New(
			govt.SetApikey(os.Getenv("VT_API_KEY")),
		)
		if err != nil {
			log.Fatal(err)
		}
		notifier := &SlackNotifier{os.Getenv("SLACK_WEBHOOK_URL"), os.Getenv("SLACK_CHANNEL")}

		det := NewDetector(vtClient, notifier)
		go det.Start()
		defer det.Shutdown()

		server := NewServer(*addr, det)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := NewClient(*addr).Add(flag.Args()); err != nil {
			log.Fatal(err)
		}
	}
}
