package main

import (
	"flag"
	"log"
	"os"

	"github.com/williballenthin/govt"
)

func example() error {
	mode := flag.String("mode", "report", "{report,scan}")
	flag.Parse()
	file := flag.Arg(0)
	log.Println(file)

	client, err := govt.New(
		govt.SetApikey(os.Getenv("VT_API_KEY")),
	)
	if err != nil {
		return err
	}

	switch *mode {
	case "scan":
		sfr, err := client.ScanFile(file)
		if err != nil {
			return err
		}
		log.Printf("ScanFileResult: %+v", sfr)

		fr, err := client.GetFileReport(sfr.Sha256)
		if err != nil {
			return err
		}
		log.Printf("FileReport: %+v", fr)
	case "report":
		h, err := sha256Sum(file)
		if err != nil {
			return err
		}
		log.Printf("hash %s", h)

		fr, err := client.GetFileReport(h)

		if err != nil {
			return err
		}
		log.Printf("FileReport: %+v", fr)
	}

	return nil
}
