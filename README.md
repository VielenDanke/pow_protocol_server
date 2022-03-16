# Proof Of Work challenge-response protocol server

### Protocol description
Client -> Server (message in format {username},{nonce})  
Server -> Client (message in format {nonce + server_nonce},{repeat_number})   
Client -> Server (message in format {nonce + server_nonce},{proof_of_work calculated with nonce + server_nonce})   
Server -> Client (message in format {nonce},{proof_of_work calculated using client nonce}).   

### Environment variables:  
SERVER_NETWORK_TYPE - type of network for server (default tcp);  
WISDOM_WORDS_FILE_PATH - path to wisdom words file;
USER_FILE_PATH - path to users key-value storage file;
SERVER_ADDRESS - address of server (default :8080);  

### Restrictions
Server could not be started without USER_FILE_PATH and WISDOM_WORDS_FILE_PATH env variables.   
By default, user file is expected to be in format key=value, each pair - new line.   
````
user=password
user1=password1
````
By default, wisdom words file is expected to be in format:
````
If you want to achieve greatness stop asking for permission. ~Anonymous
Things work out best for those who make the best of how things work out. ~John Wooden
````
Be aware of restrictions for salt, nonce symbols amount and min/max repeat number:
````go
var defaultNonceNumber = 18 // not less than
var defaultSaltNumber = 16 // not less than
var defaultServerMaxRepeatNumber = 4096 // not more than
var defaultServerMinRepeatNumber = 2048 // not less than
````

### Options
To supply headers above or configuration for nonce, salt symbols amount, min and max repeat number - 
use options (type ServerOption func(s Server)) in server_option.go file.   
Some examples below:   
````go
package main

type ServerOption func(s Server)

func WithServerNonceNumber(serverNonceNumber int) ServerOption {
	return func(s Server) {
		if serverNonceNumber < defaultNonceNumber {
			log.Printf("WARN: server nonce number cannot be less than %d\n", defaultNonceNumber)
		} else {
			s.SetNonceNumber(serverNonceNumber)
		}
	}
}

func WithServerSaltNumber(serverSaltNumber int) ServerOption {
	return func(s Server) {
		if serverSaltNumber < defaultSaltNumber {
			log.Printf("WARN: server salt number cannot be less than %d\n", defaultSaltNumber)
		} else {
			s.SetSaltNumber(serverSaltNumber)
		}
	}
}
````