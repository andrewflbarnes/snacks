package udy

import (
	"time"
)

type SendStrategy interface {
	GetNextBytes(currentReadIndex int, size int) ([]byte, int)
	Wait(currentReadIndex int, totalLength int) <-chan time.Time
}
