package main

import (
	"bufio"
	"context"
	"log"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	det := newDetector()
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
