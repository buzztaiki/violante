package main

import (
	"bufio"
	"context"
	"log"
	"os"

	"github.com/williballenthin/govt"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := govt.New(
		govt.SetApikey(os.Getenv("VT_API_KEY")),
	)
	if err != nil {
		log.Fatal(err)
	}
	notifier := &slackNotifier{os.Getenv("SLACK_WEBHOOK_URL"), os.Getenv("SLACK_CHANNEL")}

	det := newDetector(client, notifier)
	det.start(ctx)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		if s != "" {
			det.put(s)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
