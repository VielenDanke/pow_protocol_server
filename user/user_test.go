package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserCreate(t *testing.T) {
	// given
	password := "bla"

	// when
	u := NewUser(password)

	// then
	assert.Equal(t, password, u.GetPassword())
}
