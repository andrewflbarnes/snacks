// Package judy provides mechanisms for launching JSON RUDY attacks. The only real difference is that RUDY
// is typically described as being executed against more traditional form data e.g. x-www-form-urlencoded
package judy

import (
	"net"

	log "github.com/sirupsen/logrus"
)

var (
	logger = log.WithFields(log.Fields{})
)

// Judy is the interface defining the API for executing JSON RUDY attacks
//
// Execute launches a judy attack for this Judy instance on a single connection
// It returns a channel indicating when the send has complete and a channel
// indicating response data received.
// The operands define the connection to send on, a payload prefix and the number
// of arbitrary bytes to send to maintain the attack. The payload prefix is useful
// for ensuring the protocol is enforced
//
// ExecuteContinuous launches a judy attack for this Judy instance across mutliple.
// As this continuously executes Stop must be explicitly called to end it.
// The operands define the host and port to connect to, a payload prefix and the number
// of arbitrary bytes to send to maintain the attack. The payload prefix is useful
// for ensuring the protocol is enforced
//
// Stop stops the current execution in progress and any returned channel from an
// Execute method will complete
type Judy interface {
	Execute(conn net.Conn, prefix []byte, size int) chan bool
	// TODO use URI?
	ExecuteContinuous(host string, port int, prefix []byte, size int)
	// TODO Add graceful close method so we don't depend on program close
	Stop()
}

// NewJudy returns a new Judy instance with the requested send strategy
func NewJudy(sStrat SendStrategy, maxConns int) Judy {
	return &defaultJudy{
		sStrat,
		maxConns,
		0,
		false,
	}
}
