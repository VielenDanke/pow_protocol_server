package main

import (
	"bufio"
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

type Server interface {
	handleConn(conn net.Conn, ch chan<- error)
	Start() error
	preparePOWAlgorithm()
	validatePOWResult(rightAnswer, result interface{}) bool
}

var wisdomWords = make([]string, 0)
var defaultMaxServerDiffLvl = 70
var defaultMinServerDiffLvl = 30

type DefaultServer struct {
	maxFibLvl   int
	minFibLvl   int
	networkType string
	address     string
	fibMap      map[int]int64
}

func NewDefaultServer(networkType, address string) (Server, error) {
	ds := &DefaultServer{networkType: networkType, address: address, fibMap: make(map[int]int64)}
	if diffLvl := os.Getenv("SERVER_MAX_DIFFICULTY_LEVEL"); len(diffLvl) != 0 {
		diffLvlInt, diffLvlConfErr := strconv.Atoi(diffLvl)
		if diffLvlConfErr != nil {
			return nil, diffLvlConfErr
		}
		maxDiffLvl := 1 << 31
		if diffLvlInt > maxDiffLvl {
			log.Printf("INFO: difficulty max level is too big. Set to %d\n", maxDiffLvl)
			diffLvlInt = maxDiffLvl
		}
		ds.maxFibLvl = diffLvlInt
	} else {
		log.Println("INFO: set default max difficulty level")
		ds.maxFibLvl = defaultMaxServerDiffLvl
	}
	if diffLvl := os.Getenv("SERVER_MIN_DIFFICULTY_LEVEL"); len(diffLvl) != 0 {
		diffLvlInt, diffLvlConfErr := strconv.Atoi(diffLvl)
		if diffLvlConfErr != nil {
			return nil, diffLvlConfErr
		}
		minDiffLvl := 1
		if diffLvlInt < minDiffLvl {
			log.Printf("INFO: difficulty min level is too low. Set to %d\n", minDiffLvl)
			diffLvlInt = minDiffLvl
		}
		ds.minFibLvl = diffLvlInt
	} else {
		log.Println("INFO: set default min difficulty level")
		ds.minFibLvl = defaultMinServerDiffLvl
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
	ds.preparePOWAlgorithm()

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

	randInt := rand.Intn(ds.maxFibLvl-ds.minFibLvl+1) + ds.minFibLvl

	buff := make([]byte, 1024)

	_, writeErr := conn.Write([]byte(strconv.Itoa(int(randInt))))

	if writeErr != nil {
		log.Printf("WARN: cannot write to client %s\n", writeErr)
		return
	}
	readLen, readErr := conn.Read(buff)
	if readErr != nil {
		log.Printf("WARN: error reading result %s\n", readErr)
		return
	}
	readResultBuffStr := string(buff[:readLen])

	isValidResponse := ds.validatePOWResult(ds.fibMap[randInt], readResultBuffStr)

	if isValidResponse {
		idx := randInt % len(wisdomWords)
		_, writeWisdomErr := conn.Write([]byte(wisdomWords[idx]))
		if writeWisdomErr != nil {
			log.Printf("WARN: cannot write wisdom word %s", writeWisdomErr)
			return
		}
	} else {
		_, writeFailedErr := conn.Write([]byte(""))
		if writeFailedErr != nil {
			log.Printf("WARN: cannot write failed signal %s", writeFailedErr)
			return
		}
	}
	if closeErr := conn.Close(); closeErr != nil {
		ch <- fmt.Errorf("ERROR: cannot close connection %s", closeErr)
	}
}

func (ds *DefaultServer) preparePOWAlgorithm() {
	for i := ds.minFibLvl; i < ds.maxFibLvl; i++ {
		ds.fibMap[i] = int64(fib(i))
	}
}

func (ds *DefaultServer) validatePOWResult(rightAnswer, result interface{}) bool {
	intResult, convErr := strconv.Atoi(result.(string))
	intRightAnswer, _ := rightAnswer.(int64)

	if convErr != nil {
		log.Println("Incorrect result type")
		return false
	}
	return int64(intResult) == intRightAnswer
}
