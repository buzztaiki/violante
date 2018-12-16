package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client ...
type Client struct {
	addr string
}

// NewClient ...
func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

// Add ..
func (c *Client) Add(files []string) error {
	r := AddRequest{Files: files}
	body, err := json.Marshal(&r)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://"+c.addr+"/", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return fmt.Errorf("server error %s", resp.Status)
	}

	return nil
}
