// Package loris provides mechanisms for launching Slow Loris attacks
package loris

import (
	log "github.com/sirupsen/logrus"
	"net"
)

var (
	logger = log.WithFields(log.Fields{})
)

// Loris is the interface defining the API for executing Slow Loris attacks
//
// Execute launches a slow loris attack for this Loris instance.
// It returns a channel indicating when the send has complete and a channel
// indicating response data received.
// The operands define the connection to send on, a payload prefix and the number
// of arbitrary bytes to send to maintain the attack. The payload prefix is useful
// for ensuring the protocol is enforced
type Loris interface {
	Execute(conn net.Conn, prefix []byte, size int) chan bool
}

// NewLoris returns a new Loris instance with the requested send and receive
// strategies.
func NewLoris(sStrat SendStrategy) Loris {
	return defaultLoris{
		sStrat,
	}
}
