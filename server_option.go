package main

import "log"

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

func WithServerMaxRepeatNumber(serverMaxRepeatNumber int) ServerOption {
	return func(s Server) {
		if serverMaxRepeatNumber < defaultServerMaxRepeatNumber && serverMaxRepeatNumber > defaultServerMinRepeatNumber {
			s.SetMaxRepeatNumber(serverMaxRepeatNumber)
		} else {
			log.Printf("WARN: server max repeat number cannot be less than %d or greater than %d\n",
				defaultServerMinRepeatNumber, defaultServerMaxRepeatNumber)
		}
	}
}

func WithServerMinRepeatNumber(serverMinRepeatNumber int) ServerOption {
	return func(s Server) {
		if serverMinRepeatNumber > defaultServerMinRepeatNumber && serverMinRepeatNumber < defaultServerMaxRepeatNumber {
			s.SetMinRepeatNumber(serverMinRepeatNumber)
		} else {
			log.Printf("WARN: server min repeat number cannot be less than %d or greater than %d\n",
				defaultServerMinRepeatNumber, defaultServerMaxRepeatNumber)
		}
	}
}
