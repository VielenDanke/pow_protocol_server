package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

func hmacGenerator(password, generatedSalt, updatedNonce string, repeatedNumber int, preCalculatedProofCh chan<- string) {
	saltNonce := generatedSalt + updatedNonce
	h := hmac.New(sha1.New, []byte(password))
	for i := 0; i < repeatedNumber; i++ {
		h.Write([]byte(saltNonce))
	}
	preCalculatedProofCh <- hex.EncodeToString(h.Sum(nil))
}
