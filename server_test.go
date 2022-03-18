package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/vielendanke/pow_protocol_server/mocks"
	"os"
	"testing"
)

func TestDefaultServer_Start(t *testing.T) {
	// given
	server, _ := NewDefaultServer("", "")
	os.Setenv("USER_FILE_PATH", "users.txt")
	os.Setenv("WISDOM_WORDS_FILE_PATH", "wisdom_words.txt")
	ctx, cancel := context.WithCancel(context.Background())

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
	server, _ := NewDefaultServer("", "")
	os.Setenv("WISDOM_WORDS_FILE_PATH", "wisdom_words.txt")

	// when
	err := server.Start(context.Background())

	// then
	assert.NotNil(t, err)
}

func TestDefaultServer_Start_WisdomWordsFilePathIsEmpty(t *testing.T) {
	// given
	server, _ := NewDefaultServer("", "")
	os.Setenv("USER_FILE_PATH", "users.txt")
	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	// when
	err := server.Start(ctx)

	// then
	assert.NotNil(t, err)
}

func TestDefaultServer_Start_WithCustomAddressAndNetwork(t *testing.T) {
	// given
	networkType := "tcp4"
	address := ":8090"
	server, _ := NewDefaultServer(networkType, address)
	os.Setenv("USER_FILE_PATH", "users.txt")
	os.Setenv("WISDOM_WORDS_FILE_PATH", "wisdom_words.txt")
	ctx, cancel := context.WithCancel(context.Background())

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
	os.Setenv("USER_FILE_PATH", "users.txt")
	os.Setenv("WISDOM_WORDS_FILE_PATH", "wisdom_words.txt")
	server, err := NewDefaultServer("", "")
	ch := make(chan error, 1)
	conn := mocks.NewCustomConn()

	// when
	server.HandleConn(conn, ch)

	// then
	assert.Nil(t, err)
	assert.Empty(t, ch)
}