package judy

import (
	"time"
)

type SendStrategy interface {
	GetNextBytes(currentReadIndex int, payload []byte, size int) ([]byte, int)
	Wait(currentReadIndex int, totalLength int) <-chan time.Time
	// TODO graceful close
}