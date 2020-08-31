package udy

import (
	"bytes"
	"sync"
	"time"

	"github.com/andrewflbarnes/snacks/pkg/maths"
	log "github.com/sirupsen/logrus"
)

var (
	sharedSendBuffer = []byte{}
	sendMux          = sync.Mutex{}
)

func initSendBuffer(size int) {
	more := size - len(sharedSendBuffer)
	if more < 1 {
		return
	}

	sendMux.Lock()
	defer sendMux.Unlock()

	currentSize := len(sharedSendBuffer)
	more = size - currentSize
	if more > 0 {
		logger.WithFields(log.Fields{
			"currentSize":   currentSize,
			"requestedSize": size,
			"additional":    more,
		}).Trace("Increasing send buffer size")

		trail := bytes.Repeat([]byte{'a'}, more)
		sharedSendBuffer = append(sharedSendBuffer, trail...)
	}
}

type fixedByteSendStrategy struct {
	BytesPerSend int
	DelayPerSend time.Duration
}

func NewFixedByteSendStrategy(BytesPerSend int, DelayPerSend time.Duration) SendStrategy {
	initSendBuffer(BytesPerSend)
	return fixedByteSendStrategy{
		BytesPerSend,
		DelayPerSend,
	}
}

func (s fixedByteSendStrategy) GetNextBytes(currentReadIndex int, size int) ([]byte, int) {
	nextReadIndex := maths.Min(currentReadIndex+s.BytesPerSend, size)

	logger.WithFields(log.Fields{
		"sendBytes": s.BytesPerSend,
		"iCurrent":  currentReadIndex,
		"iNext":     nextReadIndex,
	}).Trace("Generate next bytes")

	arbSize := nextReadIndex - currentReadIndex
	return sharedSendBuffer[:arbSize], nextReadIndex
}

func (s fixedByteSendStrategy) Wait(currentReadIndex int, totalLength int) <-chan time.Time {
	return time.After(s.DelayPerSend)
}
