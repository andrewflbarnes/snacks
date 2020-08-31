package helper

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// HTTPServer starts an HTTP server listening on the specificed port. HTTP is "supported" only in
// that the server will read up to each newline until a blank line is read at which point it will
// attempt to read the body in segments of 20 bytes. Neither chunk encoding, content length or any
// other methods are supported to accurately read the payload.
func HTTPServer(port string, ready chan bool) {
	local := fmt.Sprintf("localhost:%s", port)
	logger.WithFields(log.Fields{
		"local": local,
	}).Info("Starting server")
	ln, err := net.Listen("tcp", local)

	if err != nil {
		logger.WithFields(log.Fields{
			"local": local,
			"error": err,
		}).Fatal("Unable to start server")
	}

	ready <- true

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.WithFields(log.Fields{
				"port":  port,
				"error": err,
			}).Warn("Failed to accept connection")
		} else {
			logger.WithFields(log.Fields{
				"port":   port,
				"remote": conn.RemoteAddr().String(),
			}).Debug("Accepted connection")
			go handleConnection(conn)
		}
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	segmentSize := 20
	fullMsg := ""
	readBody := false
	body := make([]byte, segmentSize)
	reader := bufio.NewReader(conn)

	for {
		if readBody {
			if _, err := io.ReadFull(reader, body); err != nil {
				logger.WithFields(log.Fields{
					"addr":   conn.RemoteAddr().String(),
					"error":  err,
					"length": strconv.Itoa(segmentSize),
				}).Warn("Socket read failed")
				return
			}
			logger.WithFields(log.Fields{
				"segment": string(body),
			}).Trace("Received body segment")
		} else {
			msg, err := reader.ReadString('\n')
			if err != nil {
				logger.WithFields(log.Fields{
					"addr":  conn.RemoteAddr().String(),
					"error": err,
				}).Error("Socket newline read failed")
				return
			}
			logger.WithFields(log.Fields{
				"line": msg,
			}).Trace("Received line")

			fullMsg += msg
			if msg == "\n" {
				readBody = true
			}
		}
	}
}
