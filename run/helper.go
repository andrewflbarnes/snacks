package run

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"strconv"
)

func httpServer(port int, ready chan bool) {
	strPort := strconv.Itoa(port)
	host := "localhost"
	logger.WithFields(log.Fields{
		"host": host,
		"port": port,
	}).Info("Starting server")
	ln, err := net.Listen("tcp", host+":"+strPort)

	if err != nil {
		logger.WithFields(log.Fields{
			"host":  host,
			"port":  port,
			"error": err,
		}).Fatal("Unable to start server")
	}

	ready <- true

	for {
		conn, err := ln.Accept()
		logger.WithFields(log.Fields{
			"host": host,
			"port": port,
		}).Info("Accepted connection")
		if err != nil {
			logger.WithFields(log.Fields{
				"host":  host,
				"port":  port,
				"error": err,
			}).Warn("Failed to accept connection")
		} else {
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
				}).Fatal("Socket read failed")
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
				}).Fatal("Socket newline read failed")
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
