package loris

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"net"
)

type ReceiveStrategy interface {
	GetResponse(conn net.Conn) []byte
}

type NoReceiveStrategy struct{}

func (s NoReceiveStrategy) GetResponse(conn net.Conn) []byte {
	logger.Debug("No attempt to receive response")

	return []byte{}
}

type NewlineReceiveStrategy struct{}

func (s NewlineReceiveStrategy) GetResponse(conn net.Conn) []byte {
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		logger.WithFields(log.Fields{
			"addr":  conn.RemoteAddr().String(),
			"error": err,
		}).Fatal("Unable to read newline response")
	}

	logger.WithFields(log.Fields{
		"response": status,
	}).Debug("Received newline response")

	return []byte(status)
}
