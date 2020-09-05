package snacks

import (
	"net"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

type defaultSnacks struct {
	Provider    DataProvider
	Sender      SendStrategy
	MaxConns    int
	connections int
	running     bool
}

func (l *defaultSnacks) Stop() {
	l.running = false
}

func (l *defaultSnacks) Execute(conn net.Conn, prefix []byte, size int) chan bool {
	return l.executeOnConnection(conn, prefix, size)
}

func (l *defaultSnacks) ExecuteContinuous(dest *url.URL, prefix []byte, size int) {
	target := dest.Host

	go l.track()

	for {
		// TODO make configurable
		<-time.After(100 * time.Millisecond)

		if !l.running {
			return
		}

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

func (l *defaultSnacks) track() {
	for {
		logger.Infof("Managing %d connections", l.connections)
		time.Sleep(5 * time.Second)
		if !l.running {
			return
		}
	}
}

func (l *defaultSnacks) executeOnConnection(conn net.Conn, prefix []byte, size int) chan bool {
	closed := make(chan bool)

	go l.send(conn, prefix, size, closed)

	return closed
}

func (l *defaultSnacks) send(conn net.Conn, prefix []byte, endMark int, closed chan bool) {
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

	// Note: marks are specific to the DataProvider implementation (see the docs). They may relate to the number
	// of bytes sent of the number of repeitions complete, for example.
	// The current mark should be updated by the data provider only.
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

		if !l.running {
			return
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

func (l defaultSnacks) monitorConnection(conn net.Conn, done chan<- bool) {
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
