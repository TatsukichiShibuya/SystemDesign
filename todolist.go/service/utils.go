package service

import (
	"fmt"
	"crypto/sha256"
)

func hash(passward string) string {
	h := sha256.Sum256([]byte(passward))
	return fmt.Sprintf("%x", h)
}
