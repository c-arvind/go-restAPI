package main

import (
	"crypto/sha1"
	"encoding/hex"
)

func hash(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash
}
