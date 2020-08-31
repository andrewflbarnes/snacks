package udy

import "time"

type SendStrategy interface {
	Wait(currentDataIndex int, totalLength int) <-chan time.Time
}
