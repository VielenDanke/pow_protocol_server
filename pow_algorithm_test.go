package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPowAlgorithm(t *testing.T) {
	// given
	preCalculatedProofCh := make(chan string, 1)
	password := "abc"
	salt := "salt"
	nonce := "nonce"
	repeatedNumber := 5

	// when
	hmacGenerator(password, salt, nonce, repeatedNumber, preCalculatedProofCh)

	// then
	assert.NotEmpty(t, <-preCalculatedProofCh)
}

func TestPowAlgorithm_PasswordIsEmpty(t *testing.T) {
	// given
	preCalculatedProofCh := make(chan string, 1)
	salt := "salt"
	nonce := "nonce"
	repeatedNumber := 5

	// when
	hmacGenerator("", salt, nonce, repeatedNumber, preCalculatedProofCh)

	// then
	assert.NotEmpty(t, <-preCalculatedProofCh)
}

func TestPowAlgorithm_SaltIsEmpty(t *testing.T) {
	// given
	preCalculatedProofCh := make(chan string, 1)
	password := "abc"
	nonce := "nonce"
	repeatedNumber := 5

	// when
	hmacGenerator(password, "", nonce, repeatedNumber, preCalculatedProofCh)

	// then
	assert.NotEmpty(t, <-preCalculatedProofCh)
}

func TestPowAlgorithm_NonceIsEmpty(t *testing.T) {
	// given
	preCalculatedProofCh := make(chan string, 1)
	salt := "salt"
	password := "abc"
	repeatedNumber := 5

	// when
	hmacGenerator(password, salt, "", repeatedNumber, preCalculatedProofCh)

	// then
	assert.NotEmpty(t, <-preCalculatedProofCh)
}

func TestPowAlgorithm_RepeatedNumberIsLessThanZero(t *testing.T) {
	// given
	preCalculatedProofCh := make(chan string, 1)
	salt := "salt"
	password := "abc"
	repeatedNumber := -1
	nonce := "nonce"

	// when
	hmacGenerator(password, salt, nonce, repeatedNumber, preCalculatedProofCh)

	// then
	assert.NotEmpty(t, <-preCalculatedProofCh)
}
