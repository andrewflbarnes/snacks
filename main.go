package main

import (
	"os"

	"github.com/andrewflbarnes/snacks/internal/judy"
	"github.com/andrewflbarnes/snacks/internal/loris"
	log "github.com/sirupsen/logrus"
)

var (
	logger = log.WithFields(log.Fields{})
)

func main() {
	if len(os.Args) < 2 {
		logger.Fatal("Subcommand [judy, loris] must be used")
	}

	switch os.Args[1] {
	case "judy":
		judy.Judy()
	case "loris":
		loris.Loris()
	default:
		logger.Fatal("Subcommand [judy, loris] must be used")
	}
}
