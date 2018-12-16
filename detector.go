package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/williballenthin/govt"
	"strings"
)

const (
	// The requested resource is not among the finished, queued or pending scans
	responseCodeNotFound = 0
	// Your resource is queued for analysis
	responseCodeQueued = -2
	// Scan finished, information embedded
	responseCodeSuccess = 1
	// Max number of files to getFileReports
	// see https://www.virustotal.com/ja/documentation/public-api/#getting-file-scans
	maxGetFileReports = 4
)

// Detector ...
type Detector struct {
	client     *govt.Client
	notifier   Notifier
	remains    []string
	remainsMux sync.Mutex
	ticker     *time.Ticker
}

type report struct {
	file string
	r    *govt.FileReport
}

// NewDetector ...
func NewDetector(client *govt.Client, notifier Notifier) *Detector {
	return &Detector{client: client, notifier: notifier, ticker: time.NewTicker(time.Minute / 4)}
}

// Start ...
func (d *Detector) Start() {
	for range d.ticker.C {
		d.detect()
	}
}

// Shutdown ...
func (d *Detector) Shutdown() {
	d.ticker.Stop()
}

// Add ...
func (d *Detector) Add(file string) {
	d.remainsMux.Lock()
	defer d.remainsMux.Unlock()
	d.remains = append(d.remains, file)
}

func (d *Detector) drain(n int) []string {
	d.remainsMux.Lock()
	defer d.remainsMux.Unlock()

	l := len(d.remains)
	if l == 0 {
		return nil
	}

	if l < n {
		n = l
	}

	rs := d.remains[0:n]
	rest := d.remains[n:]
	d.remains = make([]string, len(rest))
	copy(d.remains, rest)

	return rs
}

func (d *Detector) collectReports() ([]report, error) {
	files := d.drain(maxGetFileReports)
	hmap, hashes := d.collectHashes(files)
	reports, err := d.getFileReports(hashes)
	if err != nil {
		if strings.Contains(err.Error(), "No Content") {
			log.Printf("rate limit exceeded, requeue all")
			for _, f := range files {
				d.Add(f)
			}
			return nil, nil
		}
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
		r, err := d.client.GetFileReport(hashes[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get single report %s", err)
		}
		return []govt.FileReport{*r}, nil
	default:
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
		rs = append(rs, report{file: f, r: &r})
	}

	return rs
}

func (d *Detector) scanAndPut(file string) error {
	r, err := d.client.ScanFile(file)
	if err != nil {
		return err
	}
	if r.ResponseCode != responseCodeSuccess {
		return fmt.Errorf("%s: %d %s", file, r.ResponseCode, r.VerboseMsg)
	}
	d.Add(file)
	return nil
}

func (d *Detector) detect() {
	rs, err := d.collectReports()
	if err != nil {
		log.Print(err)
	}

	succeeded := make([]report, 0, len(rs))
	for _, r := range rs {
		switch r.r.ResponseCode {
		case responseCodeSuccess:
			log.Printf("%s is succeeded to report", r.file)
			succeeded = append(succeeded, r)
		case responseCodeQueued:
			log.Printf("%s is not yet processed, requeue", r.file)
			d.Add(r.file)
		case responseCodeNotFound:
			log.Printf("%s is scheduled to scan", r.file)
			r := r // create new r
			go func() {
				if err := d.scanAndPut(r.file); err != nil {
					log.Printf("scan failed %s", err)
				}
			}()
		default:
			log.Printf("%s: %d %s", r.file, r.r.Status.ResponseCode, r.r.Status.VerboseMsg)
		}
	}

	for _, r := range succeeded {
		d.notifier.SendReport(r.file, r.r)
	}
}

func (d *Detector) loop(ticker *time.Ticker) {
	for range ticker.C {
		log.Print("start detection...")
		d.detect()
	}
}
