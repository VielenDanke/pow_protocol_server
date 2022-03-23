package server

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/vielendanke/pow_protocol_server/mocks"
	"os"
	"testing"
)

func TestDefaultServer_Start(t *testing.T) {
	// given
	os.Setenv("USER_FILE_PATH", "../users.txt")
	os.Setenv("WISDOM_WORDS_FILE_PATH", "../wisdom_words.txt")
	server, _ := NewDefaultServer("", "")
	ctx, cancel := context.WithCancel(context.Background())

	defer os.Clearenv()

	cancel()

	// when
	err := server.Start(ctx)

	// then
	assert.Nil(t, err)
	assert.Equal(t, defaultNetworkType, server.GetNetworkType())
	assert.Equal(t, defaultAddress, server.GetAddress())
}

func TestDefaultServer_Start_UserFilePathIsEmpty(t *testing.T) {
	// given
	os.Setenv("WISDOM_WORDS_FILE_PATH", "../wisdom_words.txt")

	defer os.Clearenv()

	// when
	_, err := NewDefaultServer("", "")

	// then
	assert.NotNil(t, err)
}

func TestDefaultServer_Start_WisdomWordsFilePathIsEmpty(t *testing.T) {
	// given
	os.Setenv("USER_FILE_PATH", "../users.txt")

	defer os.Clearenv()

	// when
	_, err := NewDefaultServer("", "")

	// then
	assert.NotNil(t, err)
}

func TestDefaultServer_Start_WithCustomAddressAndNetwork(t *testing.T) {
	// given
	networkType := "tcp4"
	address := ":8090"
	os.Setenv("WISDOM_WORDS_FILE_PATH", "../wisdom_words.txt")
	os.Setenv("USER_FILE_PATH", "../users.txt")
	server, _ := NewDefaultServer(networkType, address)
	ctx, cancel := context.WithCancel(context.Background())

	defer os.Clearenv()
	cancel()

	// when
	err := server.Start(ctx)

	// then
	assert.Nil(t, err)
	assert.Equal(t, networkType, server.GetNetworkType())
	assert.Equal(t, address, server.GetAddress())
}

func TestDefaultServer_HandleConn(t *testing.T) {
	// given
	os.Setenv("USER_FILE_PATH", "../users.txt")
	os.Setenv("WISDOM_WORDS_FILE_PATH", "../wisdom_words.txt")
	server, err := NewDefaultServer("", "")
	ch := make(chan error, 1)
	conn := mocks.NewCustomConn()

	defer os.Clearenv()

	// when
	server.HandleConn(conn, ch)

	// then
	assert.Nil(t, err)
	assert.Empty(t, ch)
}
