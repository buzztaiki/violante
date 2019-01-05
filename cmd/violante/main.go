package main

import (
	"flag"
	"log"
	"os"

	"github.com/buzztaiki/violante"
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
		notifier := violante.NewSlackNotifier(os.Getenv("SLACK_WEBHOOK_URL"), os.Getenv("SLACK_CHANNEL"))

		det := violante.NewDetector(vtClient, notifier)
		go det.Start()
		defer det.Shutdown()

		server := violante.NewServer(*addr, det)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := violante.NewClient(*addr).Add(flag.Args()); err != nil {
			log.Fatal(err)
		}
	}
}
