// Package udy provides mechanisms for launching slow attacks. The name is taken from
// RUDY (r-u-dead-yet) but provides some slightly more generic capabilities so has been
// shortened to UDY
package udy

import (
	"net"

	log "github.com/sirupsen/logrus"
)

var (
	logger = log.WithFields(log.Fields{})
)

// Udy is the interface defining the API for executing slow attacks
//
// Execute launches a UDY type attack for this instance on a single connection
// It returns a channel indicating when the send has complete and a channel
// indicating response data received.
// The operands define the connection to send on, a payload prefix and the number
// of arbitrary bytes to send to maintain the attack. The payload prefix is useful
// for ensuring the protocol is enforced
//
// ExecuteContinuous launches a UDY attack for this instance across mutliple.
// As this continuously executes Stop must be explicitly called to end it.
// The operands define the host and port to connect to, a payload prefix and the number
// of arbitrary bytes to send to maintain the attack. The payload prefix is useful
// for ensuring the protocol is enforced
//
// Stop stops the current execution in progress and any returned channel from an
// Execute method will complete
type Udy interface {
	Execute(conn net.Conn, prefix []byte, size int) chan bool
	// TODO use URI?
	ExecuteContinuous(host string, port int, prefix []byte, size int)
}

// NewUdy returns a new Udy instance with the requested send strategy
func NewUdy(sStrat SendStrategy, maxConns int) Udy {
	return &defaultUdy{
		sStrat,
		maxConns,
		0,
	}
}
