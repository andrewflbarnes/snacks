package udy

import (
	"time"
)

type DataProvider interface {
	GetNextBytes(currentDataIndex int, size int) ([]byte, int)
	Wait(currentReadIndex int, totalLength int) <-chan time.Time
}
