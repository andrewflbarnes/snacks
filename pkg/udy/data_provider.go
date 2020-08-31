package udy

type DataProvider interface {
	GetNextBytes(currentDataIndex int, size int) ([]byte, int)
}
