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
	SStrat SendStrategy
	RStrat ReceiveStrategy
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
		segment, readIndex = l.SStrat.GetNextBytes(readIndex, payload)
		logger.WithFields(log.Fields{
			"segment":       string(segment),
			"payloadLength": length,
			"readIndex":     readIndex,
		}).Trace("Sending segment")

		l.SStrat.Wait(readIndex, length)

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

	response <- l.RStrat.GetResponse(conn)
}

func NewSendOnly(sStrat SendStrategy) Loris {
	return DefaultLoris{sStrat, NoReceiveStrategy{}}
}

func NewTest() Loris {
	return DefaultLoris{StubSendStrategy{}, NewlineReceiveStrategy{}}
}
