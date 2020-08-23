package run

import (
	"andrewflbarnes/snacks/http"
	"andrewflbarnes/snacks/loris"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

var (
	logger = log.WithFields(log.Fields{})

	flagsLoris   = flag.NewFlagSet("loris", flag.ExitOnError)
	flagLogJson  = flagsLoris.Bool("j", false, "Enables JSON logging")
	flagLogDebug = flagsLoris.Bool("v", false, "Enables debug logging")
	flagLogTrace = flagsLoris.Bool("vv", false, "Enables trace logging")
	flagTest     = flagsLoris.Bool("test", false, "Runs an embedded server to connect to")
	flagHost     = flagsLoris.String("host", "localhost", "The host to send the payload to")
	flagPath     = flagsLoris.String("path", "/", "The path to send the request to")
	flagPort     = flagsLoris.Int("port", 80, "The port to send the payload to")
	flagSize     = flagsLoris.Int("size", 1_000_000, "The size of the request payload to send")
	flagDelay    = flagsLoris.Int("sd", 50, "The delay in ms between each send")
	flagBytes    = flagsLoris.Int("sb", 5, "The number of bytes to send in each send")
)

func Loris() {
	flagsLoris.Parse(os.Args[2:])
	flag.Parse()

	host := *flagHost
	port := *flagPort

	if *flagLogTrace {
		log.SetLevel(log.TraceLevel)
	} else if *flagLogDebug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if *flagLogJson {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}

	test := *flagTest

	log.SetOutput(os.Stdout)

	logger.Info("Starting")

	// Create a new loris instance
	sendStrategy := loris.NewFixedByteSendStrategy(*flagBytes, *flagDelay)
	receiveStrategy := loris.NoReceiveStrategy{}
	l := loris.NewLoris(sendStrategy, receiveStrategy, host, port)

	if test {
		// Start a server which will receive the payload
		// Temp hacky stuff where we pass in the payload length until reading the Content-Length header or Transfer-Encoding headers is implemented.
		// Or just find a lib to do this
		serverReady := make(chan bool)
		go httpServer(port, serverReady)
		<-serverReady
	}

	// Create the HTTP request to send
	headers := map[string]string{
		"Content-Type":   "application/json",
		"Accept":         "*/*",
		"Content-Length": strconv.Itoa(*flagSize),
		"Host":           host + ":" + strconv.Itoa(port),
	}
	payload := getHttpPayload(*flagPath, http.Post, headers)

	logger.WithFields(log.Fields{
		"payload": string(payload),
	}).Debug("Generated payload")

	done, response := l.Execute(*flagSize, payload)

	// Wait for the payload to finish sending
	<-done
	logger.Info("Finished sending payload")

	// Wait for the response payload to be received and log
	responsePayload := <-response
	logger.WithFields(log.Fields{
		"response": string(responsePayload),
	}).Info("Received response")
}

func getHttpPayload(endpoint string, verb http.HttpVerb, headers map[string]string) []byte {
	builder := http.HttpRequestBuilder{
		Proto:    http.Http11,
		Verb:     verb,
		Endpoint: endpoint,
		Body:     "",
		Headers:  headers,
	}

	return builder.GetPayloadBytes()
}
