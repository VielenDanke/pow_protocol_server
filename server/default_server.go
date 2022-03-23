package server

import (
	"bufio"
	"context"
	"fmt"
	"github.com/vielendanke/pow_protocol_server/powalgo"
	"github.com/vielendanke/pow_protocol_server/user"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// For user storage we should use one of SQL or NoSQL storages instead
// but for now to make things simpler using map[string]*user
var users = make(map[string]*user.User)
var wisdomWords = make([]string, 0)
var defaultNonceNumber = 18
var defaultSaltNumber = 16
var defaultServerMaxRepeatNumber = 4096
var defaultServerMinRepeatNumber = 2048
var defaultAddress = ":8080"
var defaultNetworkType = "tcp"

type DefaultServer struct {
	serverMinRepeatNumber int
	serverMaxRepeatNumber int
	nonceNumber           int
	saltNumber            int
	networkType           string
	address               string
}

func NewDefaultServer(networkType, address string, opts ...serverOption) (Server, error) {
	ds := &DefaultServer{
		networkType:           networkType,
		address:               address,
		serverMinRepeatNumber: defaultServerMinRepeatNumber,
		serverMaxRepeatNumber: defaultServerMaxRepeatNumber,
		nonceNumber:           defaultNonceNumber,
		saltNumber:            defaultSaltNumber,
	}
	for _, v := range opts {
		v(ds)
	}
	if ds.serverMaxRepeatNumber < ds.serverMinRepeatNumber {
		log.Println("WARN: min server repeat number cannot be greater than max - set defaults")
		ds.serverMaxRepeatNumber = defaultServerMaxRepeatNumber
		ds.serverMinRepeatNumber = defaultServerMinRepeatNumber
	}
	if len(ds.networkType) == 0 {
		ds.networkType = defaultNetworkType
	}
	if len(ds.address) == 0 {
		ds.address = defaultAddress
	}
	// for config we could use separate custom interface with methods Read() Parse()
	// for simplicity - open and parse file here
	wisdomWordsFilename := os.Getenv("WISDOM_WORDS_FILE_PATH")
	usersFiles := os.Getenv("USER_FILE_PATH")

	file, fileErr := os.Open(wisdomWordsFilename)

	if fileErr != nil {
		log.Printf("ERROR: wisdom words config wisdomFile error %s\n", fileErr)
		return nil, fileErr
	}
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		wisdomPhrase := strings.Split(fileScanner.Text(), "~")[0]
		wisdomPhrase = strings.Trim(wisdomPhrase, " ")
		wisdomWords = append(wisdomWords, wisdomPhrase)
	}

	file, fileErr = os.Open(usersFiles)

	if fileErr != nil {
		log.Printf("ERROR: user config userFile error %s\n", fileErr)
		return nil, fileErr
	}

	fileScanner = bufio.NewScanner(file)

	for fileScanner.Scan() {
		userProperty := strings.Split(fileScanner.Text(), "=")
		users[userProperty[0]] = user.NewUser(userProperty[1])
	}
	return ds, nil
}

func (ds *DefaultServer) Start(ctx context.Context) error {
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
		case <-ctx.Done():
			log.Println("INFO: context canceled")
			return nil
		case err := <-ch:
			return err
		default:
			conn, acceptErr := l.Accept()

			if acceptErr != nil {
				return acceptErr
			}
			go ds.HandleConn(conn, ch)
		}
	}
}

func (ds *DefaultServer) HandleConn(conn net.Conn, ch chan<- error) {
	var readLen int
	var readErr, writeErr error
	var clientFirstMessageArr []string
	var clientSecondMessageArr []string
	var preCalculatedProof string

	log.Printf("INFO: connection from %s is accepted\n", conn.RemoteAddr().String())

	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			ch <- fmt.Errorf("ERROR: cannot close connection %s", closeErr)
		}
	}()

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

	updatedNonce := nonce + powalgo.RandomStringGenerator(ds.nonceNumber)
	generatedSalt := powalgo.RandomStringGenerator(ds.saltNumber)
	repeatedNumber := rand.Intn((ds.serverMaxRepeatNumber - ds.serverMinRepeatNumber + 1) + ds.serverMinRepeatNumber)

	preCalculatedProof = powalgo.HMACGenerator(currentUser.GetPassword(), generatedSalt, updatedNonce, repeatedNumber)

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
	clientSecondMessageArr, buff = strings.Split(string(buff[:readLen]), ","), buff[readLen+1:]

	returnedNonce := clientSecondMessageArr[0]
	proof := clientSecondMessageArr[1]

	if returnedNonce != updatedNonce {
		log.Println("ERROR: nonce is changed - close")
		_, writeErr = conn.Write([]byte{})
		if writeErr != nil {
			log.Println("ERROR: cannot write to response - close")
		}
		return
	}
	if proof != preCalculatedProof {
		log.Println("ERROR: proof is not valid - close")
		_, writeErr = conn.Write([]byte{})
		if writeErr != nil {
			log.Println("ERROR: cannot write to response - close")
		}
		return
	}
	_, writeErr = conn.Write([]byte(wisdomWords[rand.Int31()%int32(len(wisdomWords))]))

	if writeErr != nil {
		log.Println("ERROR: cannot write wisdom words")
		return
	}
}

func (ds *DefaultServer) SetSaltNumber(saltNumber int) {
	ds.saltNumber = saltNumber
}

func (ds *DefaultServer) SetNonceNumber(nonceNumber int) {
	ds.nonceNumber = nonceNumber
}

func (ds *DefaultServer) SetMaxRepeatNumber(maxNumber int) {
	ds.serverMaxRepeatNumber = maxNumber
}

func (ds *DefaultServer) SetMinRepeatNumber(minNumber int) {
	ds.serverMinRepeatNumber = minNumber
}

func (ds *DefaultServer) GetMaxRepeatNumber() int {
	return ds.serverMaxRepeatNumber
}

func (ds *DefaultServer) GetNonceNumber() int {
	return ds.nonceNumber
}

func (ds *DefaultServer) GetSaltNumber() int {
	return ds.saltNumber
}

func (ds *DefaultServer) GetMinRepeatNumber() int {
	return ds.serverMinRepeatNumber
}

func (ds *DefaultServer) GetNetworkType() string {
	return ds.networkType
}

func (ds *DefaultServer) GetAddress() string {
	return ds.address
}
