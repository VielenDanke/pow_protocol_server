package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomStringGenerator(t *testing.T) {
	// given
	type test struct {
		lenOfStr int
		expected int
		success  bool
	}
	tests := []test{
		{lenOfStr: 10, expected: 10, success: true},
		{lenOfStr: 15, expected: 15, success: true},
		{lenOfStr: 1, expected: 1, success: true},
		{lenOfStr: 25, expected: 15, success: false},
	}

	for _, v := range tests {
		// when
		str := randomStringGenerator(v.lenOfStr)

		// then
		if v.success {
			assert.Equal(t, v.expected, len(str))
		} else {
			assert.NotEqual(t, v.expected, len(str))
		}
	}
}
