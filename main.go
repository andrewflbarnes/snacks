package main

import (
	"andrewflbarnes/snacks/loris"
	"andrewflbarnes/snacks/payloads"
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
	"strconv"
)

var (
	port        = 8989
	logger      = log.WithFields(log.Fields{})
	embedServer = true
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	logger.Info("Starting")

	// Create the HTTP request to send
	tmpl := `POST {{.Endpoint}} HTTP/1.1
Content-Length: {{.Length}}
Content-Type: {{.ContentType}}
User-Agent: snacks

{{.Body}}`
	builder, err := payloads.NewHttp(tmpl)
	if err != nil {
		logger.WithFields(log.Fields{
			"template": tmpl,
			"error":    err,
		}).Fatal("Unable to create builder with template")
	}
	body := `{"ab":"cd"}`
	values := map[string]string{
		"ContentType": "application/json",
		"Body":        body,
		"Length":      strconv.Itoa(len(body)),
		"Endpoint":    "/",
	}

	payload, err := builder.BuildPayload(values)
	if err != nil {
		logger.WithFields(log.Fields{
			"template": tmpl,
			"values":   values,
			"error":    err,
		}).Fatal("Unable to populate payload template")
	}

	logger.WithFields(log.Fields{
		"payload": string(payload),
	}).Debug("Generated payload")

	// Create a new loris instance
	sendStrategy := loris.FixedByteSendStrategy{
		BytesPerSend: 5,
	}
	l := loris.New(sendStrategy, loris.NewlineReceiveStrategy{})
	// l := loris.NewTest()

	if embedServer {
		// Start a server which will receive the payload
		// Temp hacky stuff where we pass in the payload length until reading the Content-Length header or Transfer-Encoding headers is implemented.
		// Or just find a lib to do this
		serverReady := make(chan bool)
		go httpServer(port, serverReady, len(body))
		<-serverReady
	}

	// Create a client connection to end the payload over
	host := "localhost"
	conn, err := clientConnect(host, port)
	if err != nil {
		logger.WithFields(log.Fields{
			"host":  host,
			"port":  port,
			"error": err,
		}).Fatal("Unable to establish client connection")
	}

	logger.WithFields(log.Fields{
		"payload": string(payload),
	}).Info("Sending payload")
	done, response, err := l.Send(payload, conn)
	if err != nil {
		logger.WithFields(log.Fields{
			"payload": string(payload),
			"server":  conn.RemoteAddr().String(),
			"error":   err,
		}).Fatal("Unable send payload to server")
	}

	// Wait for the payload to finish sending
	<-done
	// Wait for the response payload to be received and log
	responsePayload := <-response
	logger.WithFields(log.Fields{
		"response": string(responsePayload),
	}).Info("Received response")
}

func clientConnect(host string, port int) (net.Conn, error) {
	dest := host + ":" + strconv.Itoa(port)
	logger.WithFields(log.Fields{
		"host": host,
		"port": port,
	}).Info("Opening client connection")

	return net.Dial("tcp", dest)
}

func httpServer(port int, ready chan bool, recvLength int) {
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
			go handleConnection(conn, recvLength)
		}
	}
}

func handleConnection(conn net.Conn, recvLength int) {
	fullMsg := ""

	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
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
			body := make([]byte, recvLength)
			if _, err := io.ReadFull(reader, body); err != nil {
				logger.WithFields(log.Fields{
					"addr":   conn.RemoteAddr().String(),
					"error":  err,
					"length": strconv.Itoa(recvLength),
				}).Fatal("Socket read failed")
			}
			logger.WithFields(log.Fields{
				"body": string(body),
			}).Trace("Received body")

			fullMsg += string(body)
			logger.WithFields(log.Fields{
				"payload": fullMsg,
			}).Debug("Received payload")

			response := "Received\n"
			logger.WithFields(log.Fields{
				"response": response,
			}).Debug("Sending response")

			fmt.Fprintf(conn, response)
		}
	}
}
