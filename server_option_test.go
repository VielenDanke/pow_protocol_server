package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithServerMaxRepeatNumber(t *testing.T) {
	// given
	maxRepeatNumber := 4096

	// when
	server, err := NewDefaultServer("", "", WithServerMaxRepeatNumber(maxRepeatNumber))

	// then
	assert.Equal(t, maxRepeatNumber, server.GetMaxRepeatNumber())
	assert.Nil(t, err)
}

func TestWithServerMaxRepeatNumber_IncorrectMaxNumber(t *testing.T) {
	// given
	maxRepeatNumber := 5555

	// when
	server, err := NewDefaultServer("", "", WithServerMaxRepeatNumber(maxRepeatNumber))

	// then
	assert.NotEqual(t, maxRepeatNumber, server.GetMaxRepeatNumber())
	assert.Equal(t, defaultServerMaxRepeatNumber, server.GetMaxRepeatNumber())
	assert.Nil(t, err)
}

func TestWithServerMinRepeatNumber(t *testing.T) {
	// given
	minRepeatNumber := 2144

	// when
	server, err := NewDefaultServer("", "", WithServerMinRepeatNumber(minRepeatNumber))

	// then
	assert.Equal(t, minRepeatNumber, server.GetMinRepeatNumber())
	assert.Nil(t, err)
}

func TestWithServerMinRepeatNumber_IncorrectMinNumber(t *testing.T) {
	// given
	minRepeatNumber := 1000

	// when
	server, err := NewDefaultServer("", "", WithServerMinRepeatNumber(minRepeatNumber))

	// then
	assert.NotEqual(t, minRepeatNumber, server.GetMinRepeatNumber())
	assert.Equal(t, defaultServerMinRepeatNumber, server.GetMinRepeatNumber())
	assert.Nil(t, err)
}

func TestWithServerNonceNumber(t *testing.T) {
	// given
	nonceNumber := 24

	// when
	server, err := NewDefaultServer("", "", WithServerNonceNumber(nonceNumber))

	// then
	assert.Equal(t, nonceNumber, server.GetNonceNumber())
	assert.Nil(t, err)
}

func TestWithServerNonceNumber_IncorrectNonceNumber(t *testing.T) {
	// given
	nonceNumber := 1

	// when
	server, err := NewDefaultServer("", "", WithServerNonceNumber(nonceNumber))

	// then
	assert.NotEqual(t, nonceNumber, server.GetNonceNumber())
	assert.Equal(t, defaultNonceNumber, server.GetNonceNumber())
	assert.Nil(t, err)
}

func TestWithServerSaltNumber(t *testing.T) {
	// given
	saltNumber := 24

	// when
	server, err := NewDefaultServer("", "", WithServerSaltNumber(saltNumber))

	// then
	assert.Equal(t, saltNumber, server.GetSaltNumber())
	assert.Nil(t, err)
}

func TestWithServerSaltNumber_IncorrectSaltNumber(t *testing.T) {
	// given
	saltNumber := 1

	// when
	server, err := NewDefaultServer("", "", WithServerSaltNumber(saltNumber))

	// then
	assert.NotEqual(t, saltNumber, server.GetSaltNumber())
	assert.Equal(t, defaultSaltNumber, server.GetSaltNumber())
	assert.Nil(t, err)
}
