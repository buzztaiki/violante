package violante

import "github.com/williballenthin/govt"

// Notifier ...
type Notifier interface {
	SendReport(file string, fr *govt.FileReport) error
}
