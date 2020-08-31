package udy

import (
	"bytes"
	"sync"

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

type fixedByteDataProvider struct {
	BytesPerSend int
}

func FixedByteDataProvider(BytesPerSend int) DataProvider {
	initSendBuffer(BytesPerSend)
	return fixedByteDataProvider{
		BytesPerSend,
	}
}

func (s fixedByteDataProvider) GetNextBytes(currentDataIndex int, size int) ([]byte, int) {
	nextDataIndex := maths.Min(currentDataIndex+s.BytesPerSend, size)

	logger.WithFields(log.Fields{
		"sendBytes": s.BytesPerSend,
		"iCurrent":  currentDataIndex,
		"iNext":     nextDataIndex,
	}).Trace("Generate next bytes")

	arbSize := nextDataIndex - currentDataIndex
	return sharedSendBuffer[:arbSize], nextDataIndex
}
