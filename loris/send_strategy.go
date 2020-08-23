package loris

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	sendBuffer = []byte{}
	sendMux    = sync.Mutex{}
)

func initSendBuffer(size int) {
	more := size - len(sendBuffer)
	if more < 1 {
		return
	}

	sendMux.Lock()
	defer sendMux.Unlock()

	currentSize := len(sendBuffer)
	more = size - currentSize
	if more > 0 {
		logger.WithFields(log.Fields{
			"currentSize":   currentSize,
			"requestedSize": size,
			"additional":    more,
		}).Trace("Increasing send buffer size")

		trail := bytes.Repeat([]byte{'a'}, more)
		sendBuffer = append(sendBuffer, trail...)
	}
}

type SendStrategy interface {
	GetNextBytes(currentReadIndex int, payload []byte, size int) ([]byte, int)
	Wait(currentReadIndex int, totalLength int)
}

type StubSendStrategy struct{}

func (s StubSendStrategy) GetNextBytes(currentReadIndex int, payload []byte, size int) ([]byte, int) {
	// ignore the arbitrary trailing size
	return []byte(payload), len(payload)
}

func (s StubSendStrategy) Wait(currentReadIndex int, totalLength int) {}

type fixedByteSendStrategy struct {
	BytesPerSend int
	DelayPerSend int
}

func NewFixedByteSendStrategy(BytesPerSend int, DelayPerSend int) SendStrategy {
	initSendBuffer(BytesPerSend)
	return fixedByteSendStrategy{
		BytesPerSend,
		DelayPerSend,
	}
}

func (s fixedByteSendStrategy) GetNextBytes(currentReadIndex int, payload []byte, size int) ([]byte, int) {
	nextReadIndex := min(currentReadIndex+s.BytesPerSend, size+len(payload))
	payloadLength := len(payload)

	logger.WithFields(log.Fields{
		"sendBytes": s.BytesPerSend,
		"payload":   payloadLength,
		"iCurrent":  currentReadIndex,
		"iNext":     nextReadIndex,
	}).Trace("Generate next bytes")

	if currentReadIndex >= payloadLength {
		arbSize := nextReadIndex - currentReadIndex
		return sendBuffer[:arbSize], nextReadIndex
	}

	if nextReadIndex <= payloadLength {
		return payload[currentReadIndex:nextReadIndex], nextReadIndex
	}

	trail := payload[currentReadIndex:]
	arbSize := min(s.BytesPerSend-len(trail), size)
	arb := sendBuffer[:arbSize]
	logger.WithFields(log.Fields{
		"sendBytes": s.BytesPerSend,
		"iCurrent":  currentReadIndex,
		"iNext":     nextReadIndex,
		"trail":     string(trail),
		"arb":       string(arb),
		"arbSize":   arbSize,
	}).Trace("Next bytes split")
	return append(trail, arb...), nextReadIndex
}

func (s fixedByteSendStrategy) Wait(currentReadIndex int, totalLength int) {
	// Stub implementation - could use another backoff strategy
	time.Sleep(time.Duration(s.DelayPerSend) * time.Millisecond)
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
