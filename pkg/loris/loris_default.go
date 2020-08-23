package loris

import (
	log "github.com/sirupsen/logrus"
	"net"
)

type defaultLoris struct {
	SStrat SendStrategy
}

func (l defaultLoris) Execute(conn net.Conn, prefix []byte, size int) chan bool {
	logger.WithFields(log.Fields{
		"prefix":   string(prefix),
		"dataSize": size,
	}).Info("Sending payload")

	closed := make(chan bool)

	go l.send(conn, prefix, size, closed)

	return closed
}

func (l defaultLoris) send(conn net.Conn, prefix []byte, size int, closed chan bool) {
	defer func() { closed <- true }()

	received := make(chan bool)
	go l.monitor(conn, received)

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

func (l defaultLoris) monitor(conn net.Conn, closed chan<- bool) {
	defer func() { closed <- true }()

	read := make([]byte, 1)
	_, err := conn.Read(read)
	if err != nil {
		logger.WithFields(log.Fields{
			"error":  err,
			"remote": conn.RemoteAddr().String(),
			"local":  conn.LocalAddr().String(),
		}).Error("Error while monitoring connection")
	} else {
		logger.WithFields(log.Fields{
			"remote": conn.RemoteAddr().String(),
			"local":  conn.LocalAddr().String(),
		}).Warn("Received data while monitoring connection")
	}
}
