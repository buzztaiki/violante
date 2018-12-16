package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/williballenthin/govt"
)

func main() {
	addr := flag.String("addr", ":51289", "listen address")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := govt.New(
		govt.SetApikey(os.Getenv("VT_API_KEY")),
	)
	if err != nil {
		log.Fatal(err)
	}
	notifier := &SlackNotifier{os.Getenv("SLACK_WEBHOOK_URL"), os.Getenv("SLACK_CHANNEL")}

	det := NewDetector(client, notifier)
	det.Start(ctx)

	server := NewServer(*addr, det)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
