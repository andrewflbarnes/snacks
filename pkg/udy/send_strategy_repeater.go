package udy

import (
	"time"
)

type repeaterSendStrategy struct {
	BytesToSend  []byte
	Repetitions  int
	DelayPerSend time.Duration
}

func NewRepeaterSendStrategy(BytesToSend []byte, Repetitions int, DelayPerSend time.Duration) SendStrategy {
	return repeaterSendStrategy{
		BytesToSend,
		Repetitions,
		DelayPerSend,
	}
}

func (s repeaterSendStrategy) GetNextBytes(currentReadIndex int, size int) ([]byte, int) {
	currentReadIndex++
	if currentReadIndex > s.Repetitions {
		return nil, currentReadIndex
	}
	return s.BytesToSend, currentReadIndex
}

func (s repeaterSendStrategy) Wait(currentReadIndex int, totalLength int) <-chan time.Time {
	return time.After(s.DelayPerSend)
}
