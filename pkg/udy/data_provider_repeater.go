package udy

// RepeaterDataProvider implements DataProvider providing the specified byte slice repeatedly.
// The current index is used to track how many repetitions have been returned so far and the size
// is the number of repetitions to return in total.
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
