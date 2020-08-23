package helper

import (
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
)

func DialTcp(host string, port int) (net.Conn, error) {
	dest := host + ":" + strconv.Itoa(port)
	logger.WithFields(log.Fields{
		"host": host,
		"port": port,
	}).Debug("Opening client connection")

	return net.Dial("tcp", dest)
}
