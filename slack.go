package violante

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/williballenthin/govt"
)

type slackAttachment struct {
	Title     string `json:"title,omitempty"`
	TitleLink string `json:"title_link,omitempty"`
	Text      string `json:"text,omitempty"`
	Color     string `json:"color,omitempty"`
}

type slackMessage struct {
	Channel     string            `json:"channel"`
	Attachments []slackAttachment `json:"attachments"`
}

// SlackNotifier ...
type SlackNotifier struct {
	webhookURL string
	channel    string
}

// NewSlackNotifier ...
func NewSlackNotifier(webhookURL, channel string) *SlackNotifier {
	return &SlackNotifier{webhookURL, channel}
}

func (n *SlackNotifier) send(a slackAttachment) error {
	payload, err := json.Marshal(slackMessage{
		Channel:     n.channel,
		Attachments: []slackAttachment{a},
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(n.webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to sent to slack %s (%s)", resp.Status, body)
	}

	return nil
}

// SendReport ...
func (n *SlackNotifier) SendReport(file string, fr *govt.FileReport) error {
	if fr.Positives == 0 {
		return nil
	}

	return n.send(slackAttachment{
		Title:     fmt.Sprintf("Virus detected in %s", file),
		TitleLink: fr.Permalink,
		Text:      fmt.Sprintf("Detected: %d/%d\nURL: %s", fr.Positives, fr.Total, fr.Permalink),
	})
}
