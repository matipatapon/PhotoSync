package helper

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
)

type IHasher interface {
	Hash(file []byte) (string, error)
}

type Hasher struct {
	logger *log.Logger
}

func NewHasher() Hasher {
	return Hasher{logger: log.New(os.Stdout, "[Hasher]: ", log.LstdFlags)}
}

func (h *Hasher) Hash(file []byte) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, bytes.NewReader(file)); err != nil {
		h.logger.Printf("Failed to hash a file: '%s'", err.Error())
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
