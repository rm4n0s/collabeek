package common

import (
	"crypto/sha256"
	"fmt"
)

func HashPassword(pass string) string {
	s := sha256.New()
	s.Write([]byte(pass))
	hash := s.Sum(nil)
	return fmt.Sprintf("%x", hash)
}
