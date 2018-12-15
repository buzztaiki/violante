package main

import (
	"sync"

	"context"
	"github.com/k0kubun/pp"
	"time"
)

type detector struct {
	remains    []string
	remainsMux sync.Mutex
}

func newDetector() *detector {
	return &detector{remains: nil}
}

func (d *detector) start(ctx context.Context) {
	go d.loop(ctx, time.NewTicker(time.Minute/4))
}

func (d *detector) put(file string) {
	d.remainsMux.Lock()
	defer d.remainsMux.Unlock()
	d.remains = append(d.remains, file)
}

func (d *detector) drainAll() []string {
	d.remainsMux.Lock()
	defer d.remainsMux.Unlock()
	rs := d.remains
	d.remains = nil
	return rs
}

func (d *detector) getReports() {
	rs := d.drainAll()
	pp.Printf("remains %v", rs)
}

func (d *detector) loop(ctx context.Context, ticker *time.Ticker) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.getReports()
		}
	}
}
