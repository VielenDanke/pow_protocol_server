package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserCreate(t *testing.T) {
	// given
	password := "bla"

	// when
	u := &user{password: password}

	// then
	assert.Equal(t, password, u.password)
}
