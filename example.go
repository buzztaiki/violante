package main

import (
	"flag"
	"log"
	"os"

	"fmt"

	"github.com/k0kubun/pp"
	"github.com/williballenthin/govt"
)

func example() error {
	mode := flag.String("mode", "report", "{report,scan,reporthash}")
	flag.Parse()

	client, err := govt.New(
		govt.SetApikey(os.Getenv("VT_API_KEY")),
	)
	if err != nil {
		return err
	}
	snf := slackNotifier{os.Getenv("SLACK_WEBHOOK_URL"), os.Getenv("SLACK_CHANNEL")}

	switch *mode {
	case "scan":
		sfr, err := client.ScanFile(flag.Arg(0))
		if err != nil {
			return err
		}
		log.Print(pp.Sprintf("ScanFileResult: %+v", sfr))

		fr, err := client.GetFileReport(sfr.Sha256)
		if err != nil {
			return err
		}
		log.Print(pp.Sprintf("FileReport: %+v", fr))
	case "report":
		h, err := sha256Sum(flag.Arg(0))
		if err != nil {
			return err
		}
		log.Printf("hash %s", h)

		fr, err := client.GetFileReport(h)

		if err != nil {
			return err
		}
		return snf.SendReport(flag.Arg(0), fr)
	case "reporthash":
		fr, err := client.GetFileReport(flag.Arg(0))

		if err != nil {
			return err
		}
		snf.SendReport(flag.Arg(0), fr)
	default:
		return fmt.Errorf("unknown mode %s", *mode)
	}

	return nil
}
