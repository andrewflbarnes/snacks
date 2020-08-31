package udy

type repeaterDataProvider struct {
	BytesToSend []byte
	Repetitions int
}

func NewRepeaterDataProvider(BytesToSend []byte, Repetitions int) DataProvider {
	return repeaterDataProvider{
		BytesToSend,
		Repetitions,
	}
}

func (s repeaterDataProvider) GetNextBytes(currentDataIndex int, size int) ([]byte, int) {
	currentDataIndex++
	if currentDataIndex > s.Repetitions {
		return nil, currentDataIndex
	}
	return s.BytesToSend, currentDataIndex
}
