package udy

import "time"

type FixedSendStrategy struct {
	DelayPerSend time.Duration
}

func (s FixedSendStrategy) Wait(currentDataIndex int, totalLength int) <-chan time.Time {
	return time.After(s.DelayPerSend)
}
