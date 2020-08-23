// Package loris provides mechanisms for launching Slow Loris attacks
package loris

import (
	log "github.com/sirupsen/logrus"
)

var (
	logger = log.WithFields(log.Fields{})
)

// Loris is the interface defining the API for executing Slow Loris attacks
//
// Execute launches a slow loris attack for this Loris instance.
// It returns a channel indicating when the send has complete and a channel
// indicating response data received.
type Loris interface {
	Execute(size int, prefix []byte) (chan bool, chan []byte)
}

func NewLoris(sStrat SendStrategy, rStrat ReceiveStrategy, host string, port int) Loris {
	return defaultLoris{
		sStrat,
		rStrat,
		host,
		port,
	}
}
