// Package udy provides mechanisms for launching slow attacks. The name is taken from
// RUDY (r-u-dead-yet) but provides some slightly more generic capabilities so has been
// shortened to UDY
package udy

import (
	"net"
	"net/url"

	log "github.com/sirupsen/logrus"
)

var (
	logger = log.WithFields(log.Fields{})
)

// Udy is the interface defining the API for executing slow attacks.
//
// Execute launches a UDY type attack for this instance on a single connection.
// It returns a channel indicating when the send has complete and a channel
// indicating response data received.
// The operands define the connection to send on, a payload prefix and the number
// of arbitrary bytes/repeats to send to maintain the attack per connection.
// The payload prefix is useful for ensuring the protocol is enforced.
//
// ExecuteContinuous launches a UDY attack for this instance across mutliple connections.
// As this continuously executes it relies on the program ending to stop the attack.
// The operands define the target to attach, a payload prefix and the number of arbitrary
// bytes/repeats to send to maintain the attack per connection.
// The payload prefix is useful for ensuring the protocol is enforced.
//
// Stop gracefully stops the execution in progress. The connections and executions will not
// be fully stopped after a call to this method is made. It will be complete after waiting
// for the amount of time specified as the send delay on the SendStrategy as we need to wait
// for any in progress waits to complete. A call is only required if the containing program
// is expected to continue executing after the attack is complete. Snacks does not make use
// of this.
type Udy interface {
	Execute(conn net.Conn, prefix []byte, size int) chan bool
	ExecuteContinuous(dest *url.URL, prefix []byte, size int)
	Stop()
}

// New returns a new Udy instance with a specific data provider, send strategy, and
// maximum number of connections to attack on
func New(dataProvider DataProvider, sendStrategy SendStrategy, maxConns int) Udy {
	return &defaultUdy{
		dataProvider,
		sendStrategy,
		maxConns,
		0,
		true,
	}
}
