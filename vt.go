package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/williballenthin/govt"
)

// Detector ...
type Detector struct {
	client     *govt.Client
	notifier   Notifier
	remains    []string
	remainsMux sync.Mutex
}

type report struct {
	file string
	r    *govt.FileReport
}

// NewDetector ...
func NewDetector(client *govt.Client, notifier Notifier) *Detector {
	return &Detector{client: client, notifier: notifier}
}

// Start ...
func (d *Detector) Start(ctx context.Context) {
	go d.loop(ctx, time.NewTicker(time.Minute/4))
}

// Put ...
func (d *Detector) Put(file string) {
	d.remainsMux.Lock()
	defer d.remainsMux.Unlock()
	d.remains = append(d.remains, file)
}

func (d *Detector) drainAll() []string {
	d.remainsMux.Lock()
	defer d.remainsMux.Unlock()
	rs := d.remains
	d.remains = nil
	return rs
}

func (d *Detector) collectReports() ([]report, error) {
	files := d.drainAll()
	hmap, hashes := d.collectHashes(files)
	reports, err := d.getFileReports(hashes)
	if err != nil {
		return nil, err
	}
	return d.convertReports(hmap, reports), nil
}

func (d *Detector) collectHashes(files []string) (hmap map[string]string, hashes []string) {
	hmap = make(map[string]string, len(files))
	hashes = make([]string, 0, len(files))

	for _, f := range files {
		h, err := sha256Sum(f)
		if err != nil {
			log.Print(err)
			continue
		}
		hmap[h] = f
		hashes = append(hashes, h)
	}
	return hmap, hashes
}

func (d *Detector) getFileReports(hashes []string) ([]govt.FileReport, error) {
	switch len(hashes) {
	case 0:
		return nil, nil
	case 1:
		log.Printf("1")
		r, err := d.client.GetFileReport(hashes[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get single report %s", err)
		}
		return []govt.FileReport{*r}, nil
	default:
		log.Printf("2")
		r, err := d.client.GetFileReports(hashes)
		if err != nil {
			return nil, fmt.Errorf("failed to get reports %s", err)
		}
		return *r, nil
	}
}
func (d *Detector) convertReports(hmap map[string]string, reports []govt.FileReport) []report {
	rs := make([]report, 0, len(hmap))
	for _, r := range reports {
		f, ok := hmap[r.Resource]
		if !ok {
			log.Printf("%s: not found", r.Resource)
			continue
		}

		if r.ScanId == "" {
			log.Printf("%s: %d %s", f, r.Status.ResponseCode, r.Status.VerboseMsg)
			continue
		}

		rs = append(rs, report{file: f, r: &r})
	}

	return rs
}

func (d *Detector) detect() {
	rs, err := d.collectReports()
	if err != nil {
		log.Print(err)
	}

	for _, r := range rs {
		d.notifier.SendReport(r.file, r.r)
	}
}

func (d *Detector) loop(ctx context.Context, ticker *time.Ticker) {
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			log.Print("start detection...")
			d.detect()
		}
	}
}
