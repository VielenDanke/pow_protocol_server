package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomStringGenerator(t *testing.T) {
	// when
	lenOfRandomStr := 10
	str := randomStringGenerator(lenOfRandomStr)

	// then
	assert.Equal(t, lenOfRandomStr, len(str))
}
