package violante

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func sha256Sum(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
