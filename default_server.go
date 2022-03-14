package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var users map[string]*user
var wisdomWords = make([]string, 0)
var defaultNonceNumber = 18
var defaultSaltNumber = 16
var defaultServerMaxRepeatNumber = 4096
var defaultServerMinRepeatNumber = 2048

type DefaultServer struct {
	serverMinRepeatNumber int
	serverMaxRepeatNumber int
	nonceNumber           int
	saltNumber            int
	networkType           string
	address               string
}

func NewDefaultServer(networkType, address string) (Server, error) {
	ds := &DefaultServer{
		networkType:           networkType,
		address:               address,
		serverMinRepeatNumber: defaultServerMinRepeatNumber,
		serverMaxRepeatNumber: defaultServerMaxRepeatNumber,
		nonceNumber:           defaultNonceNumber,
		saltNumber:            defaultSaltNumber,
	}
	if nonceNumber := os.Getenv("SERVER_NONCE_NUMBER"); len(nonceNumber) != 0 {
		nonceNumberFromEnv, nonceNumberParseErr := strconv.Atoi(nonceNumber)
		if nonceNumberParseErr != nil {
			return nil, nonceNumberParseErr
		}
		if nonceNumberFromEnv > defaultNonceNumber {
			ds.nonceNumber = nonceNumberFromEnv
		}
	}
	if saltNumber := os.Getenv("SERVER_SALT_NUMBER"); len(saltNumber) != 0 {
		saltNumberFromEnv, saltNumberParseErr := strconv.Atoi(saltNumber)
		if saltNumberParseErr != nil {
			return nil, saltNumberParseErr
		}
		if saltNumberFromEnv > defaultSaltNumber {
			ds.saltNumber = saltNumberFromEnv
		}
	}
	if serverMaxRepeatNum := os.Getenv("SERVER_MAX_REPEAT_NUMBER"); len(serverMaxRepeatNum) != 0 {
		serverMaxRepeatNumEnv, serverMaxRepeatNumErr := strconv.Atoi(serverMaxRepeatNum)
		if serverMaxRepeatNumErr != nil {
			return nil, serverMaxRepeatNumErr
		}
		if serverMaxRepeatNumEnv < defaultServerMaxRepeatNumber && serverMaxRepeatNumEnv > defaultServerMinRepeatNumber {
			ds.serverMaxRepeatNumber = serverMaxRepeatNumEnv
		}
	}
	if serverMinRepeatNum := os.Getenv("SERVER_MIN_REPEAT_NUMBER"); len(serverMinRepeatNum) != 0 {
		serverMinRepeatNumEnv, serverMinRepeatNumErr := strconv.Atoi(serverMinRepeatNum)
		if serverMinRepeatNumErr != nil {
			return nil, serverMinRepeatNumErr
		}
		if serverMinRepeatNumEnv > defaultServerMinRepeatNumber && serverMinRepeatNumEnv < defaultServerMaxRepeatNumber {
			ds.serverMinRepeatNumber = serverMinRepeatNumEnv
		}
	}
	if ds.serverMaxRepeatNumber < ds.serverMinRepeatNumber {
		log.Println("WARN: min server repeat number cannot be greater than max - set defaults")
		ds.serverMaxRepeatNumber = defaultServerMaxRepeatNumber
		ds.serverMinRepeatNumber = defaultServerMinRepeatNumber
	}
	return ds, nil
}

func (ds *DefaultServer) Start() error {
	file, _ := os.Open("wisdom_words.txt")

	reader := bufio.NewScanner(file)

	for reader.Scan() {
		wisdomPhrase := strings.Split(reader.Text(), "~")[0]
		wisdomPhrase = strings.Trim(wisdomPhrase, " ")
		wisdomWords = append(wisdomWords, wisdomPhrase)
	}
	l, listenErr := net.Listen(ds.networkType, ds.address)

	log.Printf("INFO: server network type: %s\n", ds.networkType)
	log.Printf("INFO: server listen: %s\n", ds.address)

	if listenErr != nil {
		return listenErr
	}
	ch := make(chan error, 1)

	go func(ch chan<- error) {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		ch <- fmt.Errorf("INFO: server interrupted: %v", <-sigCh)
	}(ch)

	for {
		select {
		case err := <-ch:
			return err
		default:
			conn, acceptErr := l.Accept()

			if acceptErr != nil {
				return acceptErr
			}
			go ds.handleConn(conn, ch)
		}
	}
}

func (ds *DefaultServer) handleConn(conn net.Conn, ch chan<- error) {
	log.Printf("INFO: connection from %s is accepted\n", conn.RemoteAddr().String())

	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			ch <- fmt.Errorf("ERROR: cannot close connection %s", closeErr)
		}
	}()

	var readLen int
	var readErr, writeErr error
	var clientFirstMessageArr []string
	var clientSecondMessageArr []string
	var preCalculatedProof string

	buff := make([]byte, 1024)

	readLen, readErr = conn.Read(buff)

	if readErr != nil {
		log.Println("ERROR: reading from connection is failed - close")
		return
	}
	clientFirstMessageArr, buff = strings.Split(string(buff[:readLen]), ","), buff[readLen+1:]

	username := clientFirstMessageArr[0]
	nonce := clientFirstMessageArr[1]

	currentUser := users[username]

	updatedNonce := nonce + randomStringGenerator(ds.nonceNumber)
	generatedSalt := randomStringGenerator(ds.saltNumber)
	repeatedNumber := rand.Intn((ds.serverMaxRepeatNumber - ds.serverMinRepeatNumber + 1) + ds.serverMinRepeatNumber)

	preCalculatedProofCh := make(chan string, 1)

	go func(password, generatedSalt, updatedNonce string, repeatedNumber int, preCalculatedProofCh chan<- string) {
		saltNonce := generatedSalt + updatedNonce
		h := hmac.New(sha1.New, []byte(password))
		for i := 0; i < repeatedNumber; i++ {
			h.Write([]byte(saltNonce))
		}
		preCalculatedProofCh <- hex.EncodeToString(h.Sum(nil))
	}(currentUser.password, generatedSalt, updatedNonce, repeatedNumber, preCalculatedProofCh)

	_, writeErr = conn.Write([]byte(fmt.Sprintf("%s,%s,%d",
		updatedNonce,
		generatedSalt,
		repeatedNumber,
	)))
	if writeErr != nil {
		log.Println("ERROR: writing to connection is failed - close")
		return
	}
	readLen, readErr = conn.Read(buff)

	if readErr != nil {
		log.Println("ERROR: reading from connection is failed - close")
		return
	}
	clientSecondMessageArr = strings.Split(string(buff[:readLen]), ",")

	returnedNonce := clientSecondMessageArr[0]
	proof := clientSecondMessageArr[1]

	if returnedNonce != updatedNonce {
		log.Println("ERROR: nonce is changed - close")
		return
	}
	preCalculatedProof = <-preCalculatedProofCh

	close(preCalculatedProofCh)

	if proof != preCalculatedProof {
		log.Println("ERROR: proof is not valid - close")
		return
	}
	_, writeErr = conn.Write([]byte(wisdomWords[rand.Int31()%int32(len(wisdomWords))]))

	if writeErr != nil {
		log.Println("ERROR: cannot write wisdom words")
		return
	}
}

func (ds *DefaultServer) validateResult(rightAnswer, result interface{}) bool {
	intResult, convErr := strconv.Atoi(result.(string))
	intRightAnswer, _ := rightAnswer.(int64)

	if convErr != nil {
		log.Println("Incorrect result type")
		return false
	}
	return int64(intResult) == intRightAnswer
}
