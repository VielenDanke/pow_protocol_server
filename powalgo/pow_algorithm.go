package powalgo

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

func HMACGenerator(password, generatedSalt, updatedNonce string, repeatedNumber int) string {
	saltNonce := generatedSalt + updatedNonce
	h := hmac.New(sha1.New, []byte(password))
	for i := 0; i < repeatedNumber; i++ {
		h.Write([]byte(saltNonce))
	}
	return hex.EncodeToString(h.Sum(nil))
}
