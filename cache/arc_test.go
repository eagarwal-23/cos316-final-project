/******************************************************************************
 * arc_test.go
 * Author:
 * Usage:    `go test`  or  `go test -v`
 * Description:
 *    An unit testing suite for arc.go.
 ******************************************************************************/
package cache

import (
	//"bytes"
	"fmt"
	"testing"
)

func TestNewArc(t *testing.T) {
	capacityArray := [4]int{16, 32, 64, 128}

	for capacity := range capacityArray {
		arc := NewArc(capacity)
		checkCapacity(t, arc, capacity)

		// Len() = 0 on init
		length := arc.Len()
		if length != 0 {
			t.Errorf("NewARC returned wrong length on init. Got %v, Expected %v", length, 0)
			t.FailNow()
		}

		// MaxStorage() = 64 on init
		maxStorage := arc.MaxStorage()
		if maxStorage != capacity {
			t.Errorf("NewFifo returned wrong maxStorage on init. Got %v, Expected %v", capacity, maxStorage)
			t.FailNow()
		}

		// RemainingStorage() = 64 on init
		remainingStorage := arc.RemainingStorage()
		if remainingStorage != capacity {
			t.Errorf("NewFifo returned wrong remainingStorage on init. Got %v, Expected %v", capacity, remainingStorage)
			t.FailNow()
		}
	}
}

// Add 20 bindings to an ARC, checking each one consumes the right storage
func TestStorageArc(t *testing.T) {
	capacity := 1024
	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	for i := 0; i < 20; i++ {
		remainingStorageBefore := arc.RemainingStorage()
		key := fmt.Sprintf("key%d", i)
		val := []byte(key)
		ok := arc.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}
		remainingStorageAfter := arc.RemainingStorage()

		expectedremainingStorageAfter := remainingStorageBefore - (len(key) + len(val))
		if remainingStorageAfter != expectedremainingStorageAfter {
			t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", remainingStorageAfter, expectedremainingStorageAfter)
			t.FailNow()
		}
	}

}

// Check that Set() adds bindings to a 'full' ARC by evicting old ones
func TestSetFullArc(t *testing.T) {
	capacity := 30
	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	arc.Set("____0", []byte("____0"))
	arc.Set("____1", []byte("____1"))
	arc.Set("____2", []byte("____2"))
	arc.Set("____3", []byte("____3"))
	//len
	len := arc.Len()
	if len != 3 {
		t.Errorf("Len wrong after adding binding to full ARC. Got %v, Expected %v", len, 3)
		t.FailNow()
	}
	arc.Set("____4", []byte("____4"))
	//len
	len = arc.Len()
	if len != 3 {
		t.Errorf("Len wrong after adding binding to full ARC. Got %v, Expected %v", len, 3)
		t.FailNow()
	}
	arc.Set("____5", []byte("____5"))
	//len
	len = arc.Len()
	if len != 3 {
		t.Errorf("Len wrong after adding binding to full ARC. Got %v, Expected %v", len, 3)
		t.FailNow()
	}

}
