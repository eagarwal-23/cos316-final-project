/******************************************************************************
 * lru_test.go
 * Author:
 * Usage:    `go test`  or  `go test -v`
 * Description:
 *    An unit testing suite for arc.go.
 ******************************************************************************/
 package cache

 import (
	//"bytes"
	//"fmt"
	"testing"
)


 func TestNewArc(t *testing.T) {
	capacityArray := [4]int{16, 32, 64, 128}

	for capacity := range capacityArray {
		lru := NewLru(capacity)
		checkCapacity(t, lru, capacity)

		// Len() = 0 on init
		length := lru.Len()
		if length != 0 {
			t.Errorf("NewFifo returned wrong length on init. Got %v, Expected %v", 0, length)
			t.FailNow()
		}

		// MaxStorage() = 64 on init
		maxStorage := lru.MaxStorage()
		if maxStorage != capacity {
			t.Errorf("NewFifo returned wrong maxStorage on init. Got %v, Expected %v", capacity, maxStorage)
			t.FailNow()
		}

		// RemainingStorage() = 64 on init
		remainingStorage := lru.RemainingStorage()
		if remainingStorage != capacity {
			t.Errorf("NewFifo returned wrong remainingStorage on init. Got %v, Expected %v", capacity, remainingStorage)
			t.FailNow()
		}
	}
}