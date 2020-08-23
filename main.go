package main

import (
	"github.com/andrewflbarnes/snacks/internal/judy"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	logger      = log.WithFields(log.Fields{})
	embedServer bool
	dest        string
)

func main() {
	if len(os.Args) < 2 {
		logger.Fatal("Subcommand \"judy\" must be used")
	}

	switch os.Args[1] {
	case "judy":
		judy.Judy()
	default:
		logger.Fatal("Subcommand \"judy\" must be used")
	}
}
