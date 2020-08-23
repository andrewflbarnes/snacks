package loris

import (
	log "github.com/sirupsen/logrus"
	"net"
)

type defaultLoris struct {
	SStrat SendStrategy
	RStrat ReceiveStrategy
	// TODO use URI object?
	Host string
	Port int
}

func (l defaultLoris) Execute(size int, prefix []byte) (chan bool, chan []byte) {

	host := l.Host
	port := l.Port

	// Create a client connection to end the payload over
	conn, err := clientConnect(host, port)
	if err != nil {
		logger.WithFields(log.Fields{
			"host":  host,
			"port":  port,
			"error": err,
		}).Fatal("Unable to establish client connection")
	}

	logger.WithFields(log.Fields{
		"prefix":   string(prefix),
		"dataSize": size,
	}).Info("Sending payload")

	sent := make(chan bool)
	response := make(chan []byte)

	go l.send(prefix, size, conn, sent, response)

	if err != nil {
		logger.WithFields(log.Fields{
			"payload":  string(prefix),
			"dataSize": size,
			"server":   conn.RemoteAddr().String(),
			"error":    err,
		}).Fatal("Unable send payload to server")
	}

	return sent, response
}

func (l defaultLoris) send(payload []byte, size int, conn net.Conn, sent chan bool, response chan []byte) {
	readIndex := 0
	length := len(payload) + size
	var segment []byte

	for readIndex < length {
		segment, readIndex = l.SStrat.GetNextBytes(readIndex, payload, size)
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
