package snacks

// DataProvider defines an interface for providing data to be included in
// payloads for attacks orchestrated by snacks.
//
// GetNextBytes returns the next sequence of bytes based on the current data index
// and the total "size" of the data required. Note that is the data provider which
// is responsible for incrementing the data index (returned as the second argument
// in the result). It is thus the data provider which defines the context of the
// size and index.
// For example we may simply use size as the maximum number of bytes required and
// the data index as how many bytes have been returned so far. Alternatively we
// could use size to represent the number of repetitions of data and the current
// index as how many repetitions of the data have been returned so far.
type DataProvider interface {
	GetNextBytes(currentDataIndex int, size int) ([]byte, int)
}
