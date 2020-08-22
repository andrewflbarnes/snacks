package loris

import (
	log "github.com/sirupsen/logrus"
	"net"
)

var (
	logger = log.WithFields(log.Fields{})
)

type Loris interface {
	Send(payload []byte, conn net.Conn) (chan bool, chan []byte, error)
}

type DefaultLoris struct {
	sStrat SendStrategy
	rStrat ReceiveStrategy
}

func (l DefaultLoris) Send(payload []byte, conn net.Conn) (chan bool, chan []byte, error) {
	sent := make(chan bool)
	response := make(chan []byte)

	go l.send(payload, conn, sent, response)

	return sent, response, nil
}

func (l DefaultLoris) send(payload []byte, conn net.Conn, sent chan bool, response chan []byte) {
	readIndex := 0
	length := len(payload)
	var segment []byte

	for readIndex < length {
		segment, readIndex = l.sStrat.GetNextBytes(readIndex, payload)
		logger.WithFields(log.Fields{
			"segment":       string(segment),
			"payloadLength": length,
			"readIndex":     readIndex,
		}).Trace("Sending segment")

		l.sStrat.Wait(readIndex, length)

		if _, err := conn.Write(segment); err != nil {
			logger.WithFields(log.Fields{
				"segment":       segment,
				"payloadLength": length,
				"readIndex":     readIndex,
				"error":         err,
			}).Fatal("Failed while writing payload segment")
		}
	}

	sent <- true
	logger.Debug("Payload sent")

	response <- l.rStrat.GetResponse(conn)
}

func New(sStrat SendStrategy, rStrat ReceiveStrategy) Loris {
	return DefaultLoris{
		sStrat,
		rStrat,
	}
}

func NewSendOnly(sStrat SendStrategy) Loris {
	return New(sStrat, NoReceiveStrategy{})
}

func NewTest() Loris {
	return New(StubSendStrategy{}, NewlineReceiveStrategy{})
}
