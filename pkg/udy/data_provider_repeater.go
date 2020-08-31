package udy

import (
	"time"
)

type repeaterSendStrategy struct {
	BytesToSend  []byte
	Repetitions  int
	DelayPerSend time.Duration
}

func NewRepeaterDataProvider(BytesToSend []byte, Repetitions int, DelayPerSend time.Duration) DataProvider {
	return repeaterSendStrategy{
		BytesToSend,
		Repetitions,
		DelayPerSend,
	}
}

func (s repeaterSendStrategy) GetNextBytes(currentDataIndex int, size int) ([]byte, int) {
	currentDataIndex++
	if currentDataIndex > s.Repetitions {
		return nil, currentDataIndex
	}
	return s.BytesToSend, currentDataIndex
}

func (s repeaterSendStrategy) Wait(currentDataIndex int, totalLength int) <-chan time.Time {
	return time.After(s.DelayPerSend)
}
