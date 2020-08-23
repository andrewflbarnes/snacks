package main

import (
	"andrewflbarnes/snacks/run"
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
		logger.Fatal("Subcommand \"loris\" must be used")
	}

	switch os.Args[1] {
	case "loris":
		run.Loris()
	default:
		logger.Fatal("Subcommand \"loris\" must be used")
	}
}
