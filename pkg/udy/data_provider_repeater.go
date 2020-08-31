package udy

type RepeaterDataProvider struct {
	BytesToSend []byte
	Repetitions int
}

func (s RepeaterDataProvider) GetNextBytes(currentDataIndex int, size int) ([]byte, int) {
	currentDataIndex++
	if currentDataIndex > s.Repetitions {
		return nil, currentDataIndex
	}
	return s.BytesToSend, currentDataIndex
}
