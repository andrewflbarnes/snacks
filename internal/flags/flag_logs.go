package flags

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

func InitLogFlags(flagSet *flag.FlagSet) LogFlags {
	return logFlagsImpl{
		json:  flagSet.Bool("j", false, "Enables JSON logging"),
		debug: flagSet.Bool("v", false, "Enables debug logging"),
		trace: flagSet.Bool("vv", false, "Enables trace logging"),
	}
}

type LogFlags interface {
	IsJson() bool
	IsDebug() bool
	IsTrace() bool
	Apply()
}

type logFlagsImpl struct {
	json  *bool
	debug *bool
	trace *bool
}

func (l logFlagsImpl) IsJson() bool {
	return *l.json
}

func (l logFlagsImpl) IsDebug() bool {
	return *l.debug
}

func (l logFlagsImpl) IsTrace() bool {
	return *l.trace
}

func (l logFlagsImpl) Apply() {
	if l.IsTrace() {
		log.SetLevel(log.TraceLevel)
	} else if l.IsDebug() {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if l.IsJson() {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}

	log.SetOutput(os.Stdout)
}
