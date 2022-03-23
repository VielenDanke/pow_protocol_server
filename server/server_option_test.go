package server

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestWithServerMaxRepeatNumber(t *testing.T) {
	// given
	os.Setenv("USER_FILE_PATH", "../users.txt")
	os.Setenv("WISDOM_WORDS_FILE_PATH", "../wisdom_words.txt")
	maxRepeatNumber := 4096

	defer os.Clearenv()

	// when
	server, err := NewDefaultServer("", "", WithServerMaxRepeatNumber(maxRepeatNumber))

	// then
	assert.Equal(t, maxRepeatNumber, server.GetMaxRepeatNumber())
	assert.Nil(t, err)
}

func TestWithServerMaxRepeatNumber_IncorrectMaxNumber(t *testing.T) {
	// given
	os.Setenv("USER_FILE_PATH", "../users.txt")
	os.Setenv("WISDOM_WORDS_FILE_PATH", "../wisdom_words.txt")
	maxRepeatNumber := 5555

	defer os.Clearenv()

	// when
	server, err := NewDefaultServer("", "", WithServerMaxRepeatNumber(maxRepeatNumber))

	// then
	assert.NotEqual(t, maxRepeatNumber, server.GetMaxRepeatNumber())
	assert.Equal(t, defaultServerMaxRepeatNumber, server.GetMaxRepeatNumber())
	assert.Nil(t, err)
}

func TestWithServerMinRepeatNumber(t *testing.T) {
	// given
	os.Setenv("USER_FILE_PATH", "../users.txt")
	os.Setenv("WISDOM_WORDS_FILE_PATH", "../wisdom_words.txt")
	minRepeatNumber := 2144

	defer os.Clearenv()

	// when
	server, err := NewDefaultServer("", "", WithServerMinRepeatNumber(minRepeatNumber))

	// then
	assert.Equal(t, minRepeatNumber, server.GetMinRepeatNumber())
	assert.Nil(t, err)
}

func TestWithServerMinRepeatNumber_IncorrectMinNumber(t *testing.T) {
	// given
	os.Setenv("USER_FILE_PATH", "../users.txt")
	os.Setenv("WISDOM_WORDS_FILE_PATH", "../wisdom_words.txt")
	minRepeatNumber := 1000

	defer os.Clearenv()

	// when
	server, err := NewDefaultServer("", "", WithServerMinRepeatNumber(minRepeatNumber))

	// then
	assert.NotEqual(t, minRepeatNumber, server.GetMinRepeatNumber())
	assert.Equal(t, defaultServerMinRepeatNumber, server.GetMinRepeatNumber())
	assert.Nil(t, err)
}

func TestWithServerNonceNumber(t *testing.T) {
	// given
	os.Setenv("USER_FILE_PATH", "../users.txt")
	os.Setenv("WISDOM_WORDS_FILE_PATH", "../wisdom_words.txt")
	nonceNumber := 24

	defer os.Clearenv()

	// when
	server, err := NewDefaultServer("", "", WithServerNonceNumber(nonceNumber))

	// then
	assert.Equal(t, nonceNumber, server.GetNonceNumber())
	assert.Nil(t, err)
}

func TestWithServerNonceNumber_IncorrectNonceNumber(t *testing.T) {
	// given
	os.Setenv("USER_FILE_PATH", "../users.txt")
	os.Setenv("WISDOM_WORDS_FILE_PATH", "../wisdom_words.txt")
	nonceNumber := 1

	defer os.Clearenv()

	// when
	server, err := NewDefaultServer("", "", WithServerNonceNumber(nonceNumber))

	// then
	assert.NotEqual(t, nonceNumber, server.GetNonceNumber())
	assert.Equal(t, defaultNonceNumber, server.GetNonceNumber())
	assert.Nil(t, err)
}

func TestWithServerSaltNumber(t *testing.T) {
	// given
	os.Setenv("USER_FILE_PATH", "../users.txt")
	os.Setenv("WISDOM_WORDS_FILE_PATH", "../wisdom_words.txt")
	saltNumber := 24

	defer os.Clearenv()

	// when
	server, err := NewDefaultServer("", "", WithServerSaltNumber(saltNumber))

	// then
	assert.Equal(t, saltNumber, server.GetSaltNumber())
	assert.Nil(t, err)
}

func TestWithServerSaltNumber_IncorrectSaltNumber(t *testing.T) {
	// given
	os.Setenv("USER_FILE_PATH", "../users.txt")
	os.Setenv("WISDOM_WORDS_FILE_PATH", "../wisdom_words.txt")
	saltNumber := 1

	defer os.Clearenv()

	// when
	server, err := NewDefaultServer("", "", WithServerSaltNumber(saltNumber))

	// then
	assert.NotEqual(t, saltNumber, server.GetSaltNumber())
	assert.Equal(t, defaultSaltNumber, server.GetSaltNumber())
	assert.Nil(t, err)
}
