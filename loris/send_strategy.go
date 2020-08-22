package loris

import (
	"time"
)

type SendStrategy interface {
	GetNextBytes(currentReadIndex int, payload []byte) ([]byte, int)
	Wait(currentReadIndex int, totalLength int)
}

type StubSendStrategy struct{}

func (s StubSendStrategy) GetNextBytes(currentReadIndex int, payload []byte) ([]byte, int) {
	return []byte(payload), len(payload)
}

func (s StubSendStrategy) Wait(currentReadIndex int, totalLength int) {}

type FixedByteSendStrategy struct {
	BytesPerSend int
}

func (s FixedByteSendStrategy) GetNextBytes(currentReadIndex int, payload []byte) ([]byte, int) {
	nextReadIndex := min(currentReadIndex+s.BytesPerSend, len(payload))
	return payload[currentReadIndex:nextReadIndex], nextReadIndex
}

func (s FixedByteSendStrategy) Wait(currentReadIndex int, totalLength int) {
	// Stub implementation - could use another backoff strategy
	time.Sleep(100 * time.Millisecond)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
