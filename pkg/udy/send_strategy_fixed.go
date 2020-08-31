package udy

import "time"

// FixedSendStrategy always waits for the same defined amount of time on every call to Wait.
type FixedSendStrategy struct {
	DelayPerSend time.Duration
}

func (s FixedSendStrategy) Wait(currentDataIndex int, size int) <-chan time.Time {
	return time.After(s.DelayPerSend)
}
