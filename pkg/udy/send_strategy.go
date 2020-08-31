package udy

import "time"

// SendStrategy defines a data for how data is sent in attacks orchestrated by udy.
//
// Wait should return a channel which is received from at some point in the future,
// typically by returning some kind of time.After(duration) channel. The current data
// index and size are the same values as defined on the DataProvider and may be used
// by the implementation if desired to vary the wait time.
type SendStrategy interface {
	Wait(currentDataIndex int, size int) <-chan time.Time
}
