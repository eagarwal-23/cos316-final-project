/******************************************************************************
 * arc_test.go
 * Author:
 * Usage:    `go test`  or  `go test -v`
 * Description:
 *    An unit testing suite for arc.go.
 ******************************************************************************/
package cache

import (
	"bytes"
	"fmt"

	//"math/rand"
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

// Check that Get() returns no binding when called on an empty ARC
func TestGetEmptyArc(t *testing.T) {
	capacity := 1024
	keysArray := [4]string{"Hello", "a", "ssup"}

	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	for _, key := range keysArray {
		value, _ := arc.Get(key)
		if value != nil {
			t.Errorf("Returned wrong value for empty ARC. Got %v, Expected %v", value, nil)
			t.FailNow()
		}
	}
}

// Check that Peek() returns no binding when called on an empty ARC
func TestPeekEmptyArc(t *testing.T) {
	capacity := 1024
	keysArray := [4]string{"Hello", "a", "ssup"}

	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	for _, key := range keysArray {
		value, _ := arc.Peek(key)
		if value != nil {
			t.Errorf("Returned wrong value for empty ARC. Got %v, Expected %v", value, nil)
			t.FailNow()
		}
	}
}

// Check various operations on an ARC with a single binding
func TestSingleBindingArc(t *testing.T) {
	capacitiesArray := [3]int{16, 64, 256}
	keysArray := [3]string{"Hello", "Foo", "COS"}
	valuesArray := [3]string{"World", "Bar", "316"}

	for i, _ := range keysArray {
		arc := NewArc(capacitiesArray[i])
		checkCapacity(t, arc, capacitiesArray[i])
		arc.Set(keysArray[i], []byte(valuesArray[i]))
		value, ok := arc.Get(keysArray[i])

		if ok {
			res := bytes.Compare(value, []byte(valuesArray[i]))
			if res != 0 {
				t.Errorf("Returned wrong value for key. Got %v, Expected %v", value, []byte(valuesArray[i]))
				t.FailNow()
			}
		} else {
			t.Errorf("Expected value but did not get one")
			t.FailNow()
		}

		remainingStorage := arc.RemainingStorage()
		expectedremainingStorage := capacitiesArray[i] - (len(keysArray[i]) + len([]byte(valuesArray[i])))
		if remainingStorage != expectedremainingStorage {
			t.Errorf("Returned wrong remaining after  for key. Got %v, Expected %v", remainingStorage, expectedremainingStorage)
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

		expectedRemainingStorageAfter := remainingStorageBefore - (len(key) + len(val))
		if remainingStorageAfter != expectedRemainingStorageAfter {
			t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", remainingStorageAfter, expectedRemainingStorageAfter)
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

// Check that Set() rejects bindings too large for the ARC
func TestSetTooLargeArc(t *testing.T) {
	capacity := 10
	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	ok := arc.Set("123456", []byte("123456"))
	if ok {
		t.Errorf("Failed to reject binding too large for ARC. Set  Got %v, Expected %v", ok, false)
		t.FailNow()
	}
	if arc.RemainingStorage() != 10 {
		t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", arc.RemainingStorage(), 10)
		t.FailNow()
	}
	_, ok = arc.Get("123456")
	if ok {
		t.Errorf("Failed to reject binding too large for ARC. Set  Got %v, Expected %v", ok, false)
		t.FailNow()
	}
}

// Check that Set() only allows zero-size bindings in a zero-capacity ARC

func TestSetZeroArc(t *testing.T) {
	capacity := 0
	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	ok := arc.Set("hello", []byte("world"))
	if ok {
		t.Errorf("Failed to reject binding too large for ARC. Set  Got %v, Expected %v", ok, false)
		t.FailNow()
	}
	ok = arc.Set("foo", []byte("boo"))
	if ok {
		t.Errorf("Failed to reject binding too large for ARC. Set  Got %v, Expected %v", ok, false)
		t.FailNow()
	}
	ok = arc.Set("", []byte(""))
	if !ok {
		t.Errorf("Failed to reject binding too large for ARC. Set  Got %v, Expected %v", ok, true)
		t.FailNow()
	}
}

// Check that the ARC allows the empty string as a valid key
func TestEmptyStringValidArc(t *testing.T) {
	capacity := 1024
	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	ok := arc.Set("", []byte("Value"))
	if arc.RemainingStorage() != 1019 {
		t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", arc.RemainingStorage(), 1019)
		t.FailNow()
	}
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	}
	if arc.Len() != 1 {
		t.Errorf("Len() wrong. Got %v, Expected %v", arc.Len(), 1)
		t.FailNow()
	}
	if arc.MaxStorage() != capacity {
		t.Errorf("MaxStorage wrong. Got %v, Expected %v", arc.MaxStorage(), capacity)
		t.FailNow()
	}

	value, ok := arc.Get("")
	res := bytes.Compare(value, []byte("Value"))
	if res != 0 {
		t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte("Value"))
		t.FailNow()
	}
}

// Check that the ARC allows the empty []byte as a valid value
func TestEmptyValidArc(t *testing.T) {
	capacity := 1024
	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	ok := arc.Set("key", []byte{})
	if arc.RemainingStorage() != 1021 {
		t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", arc.RemainingStorage(), 1021)
		t.FailNow()
	}
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	}
	if arc.Len() != 1 {
		t.Errorf("Len() wrong. Got %v, Expected %v", arc.Len(), 1)
		t.FailNow()
	}
	if arc.MaxStorage() != capacity {
		t.Errorf("MaxStorage wrong. Got %v, Expected %v", arc.MaxStorage(), capacity)
		t.FailNow()
	}

	value, ok := arc.Get("key")
	res := bytes.Compare(value, []byte{})
	if res != 0 {
		t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte{})
		t.FailNow()
	}

}

// Check that the ARC allows the empty []byte as a valid value
func TestNilValidArc(t *testing.T) {
	capacity := 1024
	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	ok := arc.Set("key", nil)
	if arc.RemainingStorage() != 1021 {
		t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", arc.RemainingStorage(), 1021)
		t.FailNow()
	}
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	}
	if arc.Len() != 1 {
		t.Errorf("Len() wrong. Got %v, Expected %v", arc.Len(), 1)
		t.FailNow()
	}
	if arc.MaxStorage() != capacity {
		t.Errorf("MaxStorage wrong. Got %v, Expected %v", arc.MaxStorage(), capacity)
		t.FailNow()
	}

	value, ok := arc.Get("key")
	res := bytes.Compare(value, nil)
	if res != 0 {
		t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, nil)
		t.FailNow()
	}

}

// Check that values can be non-ASCII (binary)
func TestBinaryValuesArc(t *testing.T) {
	capacity := 1024
	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	ok := arc.Set("key", []byte("\x00\x01�\x15�"))
	if arc.RemainingStorage() != 1012 {
		t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", arc.RemainingStorage(), 1012)
		t.FailNow()
	}
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	}
	if arc.Len() != 1 {
		t.Errorf("Len() wrong. Got %v, Expected %v", arc.Len(), 1)
		t.FailNow()
	}
	if arc.MaxStorage() != capacity {
		t.Errorf("MaxStorage wrong. Got %v, Expected %v", arc.MaxStorage(), capacity)
		t.FailNow()
	}

	value, ok := arc.Get("key")
	res := bytes.Compare(value, []byte("\x00\x01�\x15�"))
	if res != 0 {
		t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte("\x00\x01�\x15�"))
		t.FailNow()
	}

}

// Check that keys and values can be non-ASCII (Unicode)
func TestUnicodeValuesArc(t *testing.T) {
	capacity := 1024
	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	ok := arc.Set("😂_🚀", []byte("✔_🚗"))
	if arc.RemainingStorage() != 1007 {
		t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", arc.RemainingStorage(), 1007)
		t.FailNow()
	}
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	}
	if arc.Len() != 1 {
		t.Errorf("Len() wrong. Got %v, Expected %v", arc.Len(), 1)
		t.FailNow()
	}
	if arc.MaxStorage() != capacity {
		t.Errorf("MaxStorage wrong. Got %v, Expected %v", arc.MaxStorage(), capacity)
		t.FailNow()
	}

	value, ok := arc.Get("😂_🚀")
	res := bytes.Compare(value, []byte("✔_🚗"))
	if res != 0 {
		t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte("✔_🚗"))
		t.FailNow()
	}

}

// Test that Set() overwrites values when called with an existing key
func TestSetOverwriteArc(t *testing.T) {
	capacity := 1024
	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	ok := arc.Set("key", []byte("old"))
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	} else {
		value, _ := arc.Get("key")
		res := bytes.Compare(value, []byte("old"))
		if res != 0 {
			t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte("old"))
			t.FailNow()
		}
	}

	ok = arc.Set("key", []byte("new"))
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	} else {
		value, _ := arc.Get("key")
		res := bytes.Compare(value, []byte("new"))
		if res != 0 {
			t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte("new"))
			t.FailNow()
		}

	}
}

// Test that Set() overwrites values when called with an existing key
func TestSetOverwriteStorageArc(t *testing.T) {
	capacity := 1024
	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	ok := arc.Set("key", []byte("old"))
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	} else {
		value, _ := arc.Get("key")
		res := bytes.Compare(value, []byte("old"))
		if res != 0 {
			t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte("old"))
			t.FailNow()
		}

		if arc.RemainingStorage() != 1018 {
			t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", arc.RemainingStorage(), 1018)
			t.FailNow()
		}
	}

	ok = arc.Set("key", []byte("nw"))
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	} else {
		value, _ := arc.Get("key")
		res := bytes.Compare(value, []byte("nw"))
		if res != 0 {
			t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte("nw"))
			t.FailNow()
		}
		if arc.RemainingStorage() != 1019 {
			t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", arc.RemainingStorage(), 1019)
			t.FailNow()
		}

	}
}

// Check that Remove() prevents Get() from retrieving a binding
func TestRemovePreventGetArc(t *testing.T) {
	capacity := 1024
	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	ok := arc.Set("key", []byte("value"))
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	} else {
		value, _ := arc.Get("key")
		res := bytes.Compare(value, []byte("value"))
		if res != 0 {
			t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte("value"))
			t.FailNow()
		}
	}

	_, ok = arc.Remove("key")
	if !ok {
		t.Errorf("Failed to remove  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	} else {
		_, ok := arc.Get("key")
		if ok {
			t.Errorf("Fetched a removed value. Set  Got %v, Expected %v", ok, false)
			t.FailNow()
		}
	}

}

// Check that Remove() correctly updates available storage
func TestRemoveStorageArc(t *testing.T) {
	capacity := 1024
	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	arc.Set("__0", []byte("__0"))
	arc.Set("__1", []byte("__1"))
	arc.Set("__2", []byte("__2"))
	arc.Set("__3", []byte("__3"))
	arc.Remove("__0")

	//len
	len := arc.Len()
	if len != 3 {
		t.Errorf("Len wrong after removing. Got %v, Expected %v", len, 3)
		t.FailNow()
	}
	if arc.RemainingStorage() != 1006 {
		t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", arc.RemainingStorage(), 1006)
		t.FailNow()
	}
	arc.Remove("__1")
	len = arc.Len()
	//len
	if len != 2 {
		t.Errorf("Len wrong after removing. Got %v, Expected %v", len, 2)
		t.FailNow()
	}
	if arc.RemainingStorage() != 1012 {
		t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", arc.RemainingStorage(), 1012)
		t.FailNow()
	}
}

// Check that Stats() returns correct values when there are mixed cache hits and misses
func TestStatsArc(t *testing.T) {
	capacity := 1024
	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	arc.Set("____1", []byte("____1"))
	_, ok := arc.Get("____1")

	if !ok {
		t.Errorf("Failed to fetch  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	}

	_, ok = arc.Get("miss")

	if ok {
		t.Errorf("Fetched absent binding binding. Got %v, Expected %v", ok, false)
		t.FailNow()
	}

	hits, misses := arc.Stats().Hits, arc.stats.Misses

	if hits != 1 {
		t.Errorf("Hits wrong. Got %v, Expected %v", hits, 1)
		t.FailNow()
	}

	if misses != 1 {
		t.Errorf("Misses wrong. Got %v, Expected %v", misses, 1)
		t.FailNow()
	}

}

// Check that Remove() works as expected on bindings whose values have been overwritten
func TestRemoveOverwrittenArc(t *testing.T) {
	capacity := 1024
	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	arc.Set("key", []byte("old"))
	arc.Set("key", []byte("newval"))
	arc.Remove("key")
	_, ok := arc.Get("key")
	if ok {
		t.Errorf("Failed to remove binding with key: 'key'. Got %v, Expected %v", ok, false)
		t.FailNow()
	}

}

// Check that Remove() has no effect when called on an empty ARC
func TestRemoveEmptyArc(t *testing.T) {
	capacity := 1024
	arc := NewArc(capacity)
	checkCapacity(t, arc, capacity)

	_, ok := arc.Remove("key")
	if ok {
		t.Errorf("Removed empty binding. Got %v, Expected %v", ok, false)
		t.FailNow()
	}
	_, ok = arc.Remove("foo")
	if ok {
		t.Errorf("Removed empty binding. Got %v, Expected %v", ok, false)
		t.FailNow()
	}
	_, ok = arc.Remove("bar")
	if ok {
		t.Errorf("Removed empty binding. Got %v, Expected %v", ok, false)
		t.FailNow()
	}
}

// Attempt to Remove() a binding that has already been removed.
func TestRemoveRemovedArc(t *testing.T) {
	capacity := 1024
	arc := NewLru(capacity)
	checkCapacity(t, arc, capacity)

	arc.Set("key", []byte("value"))
	_, ok := arc.Remove("key")
	if !ok {
		t.Errorf("Failed to remove binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	}
	_, ok = arc.Remove("key")
	if ok {
		t.Errorf("Removed removed binding. Got %v, Expected %v", ok, false)
		t.FailNow()
	}

}

// Ensures that peek does not update hits and misses as it does not count as an access.
func TestPeekStatsArc(t *testing.T) {
	capacity := 64
	arc := NewLru(capacity)
	checkCapacity(t, arc, capacity)

	for i := 0; i < 15; i++ {
		key := fmt.Sprintf("key%d", i)
		val := []byte(key)
		ok := arc.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}

		res, _ := arc.Get(key)
		if !bytesEqual(res, val) {
			t.Errorf("Wrong value %s for binding with key: %s", res, key)
			t.FailNow()
		}

		// Peek value
		hits_before_access := arc.stats.Hits
		misses_before_access := arc.stats.Misses
		res2, _ := arc.Peek(key)
		hits_after_access := arc.stats.Hits
		misses_after_access := arc.stats.Misses

		if !bytesEqual(res, val) {
			t.Errorf("Wrong value %s for binding with key: %s", res2, key)
			t.FailNow()
		}
		//check that peek did not interfere with stats - order, hits, missess
		if hits_before_access != hits_after_access {
			t.Errorf("Wrong value %s for binding with key: %s", res2, key)
			t.FailNow()
		}
		if misses_before_access != misses_after_access {
			t.Errorf("Wrong value %s for binding with key: %s", res2, key)
			t.FailNow()
		}
	}
}

// Ensures that peek does not change order of recent accesses.
func TestPeekRecencyArc(t *testing.T) {
	capacity := 3
	arc := NewLru(capacity)
	checkCapacity(t, arc, capacity)

	arc.Set("a", []byte(""))
	arc.Set("b", []byte(""))
	arc.Peek("a")

	arc.Set("c", []byte(""))

	_, ok := arc.Get("a")

	if !ok {
		t.Errorf("should not have updated recent-ness of a. Got %v, Expected %v", ok, false)
		t.FailNow()
	}

}

// Ensures that keys on recentLRU that are re-accessed are properly moved to frequentLRU
func TestGetMoveRecentToFrequentArc(t *testing.T) {
	capacity := 256
	arc := NewArc(capacity)

	numItemsAdded := 0
	itemsAdded := make([]string, 0)

	// for i < capacity, try adding key-val pair until capacity is full
	for i := 0; i < capacity; i++ {
		key := fmt.Sprintf("key%d", i)
		if arc.capacity > arc.currentlyUsedCapacity+len(key) {
			arc.Set(key, make([]byte, 0))
			itemsAdded = append(itemsAdded, key)
			numItemsAdded++
		}
	}

	if arc.t1.Len() != numItemsAdded {
		t.Errorf("Recently-used cache t1 has wrong length.  Got %v, Expected %v", arc.t1.Len(), numItemsAdded)
		t.FailNow()
	}

	if arc.t2.Len() != 0 {
		t.Errorf("Recently-used cache t2 has wrong length.  Got %v, Expected %v", arc.t2.Len(), 0)
		t.FailNow()
	}

	// Get should upgrade to t2
	for _, key := range itemsAdded {
		if _, ok := arc.Get(key); !ok {
			t.Errorf("Missing key in cache: %v", key)
			t.FailNow()
		}
	}

	if arc.t1.Len() != 0 {
		t.Errorf("Recently-used cache t1 has wrong length.  Got %v, Expected %v", arc.t1.Len(), 0)
		t.FailNow()
	}

	if arc.t2.Len() != numItemsAdded {
		t.Errorf("Recently-used cache t2 has wrong length.  Got %v, Expected %v", arc.t2.Len(), numItemsAdded)
		t.FailNow()
	}

	// Get should get from t2
	for _, key := range itemsAdded {
		if _, ok := arc.Get(key); !ok {
			t.Errorf("Missing key in cache: %v", key)
			t.FailNow()
		}
	}

	if arc.t1.Len() != 0 {
		t.Errorf("Recently-used cache t1 has wrong length.  Got %v, Expected %v", arc.t1.Len(), 0)
		t.FailNow()
	}

	if arc.t2.Len() != numItemsAdded {
		t.Errorf("Recently-used cache t2 has wrong length.  Got %v, Expected %v", arc.t2.Len(), numItemsAdded)
		t.FailNow()
	}
}

func TestAddMoveRecentToFrequentArc(t *testing.T) {
	capacity := 256
	arc := NewArc(capacity)

	numItemsAdded := 0
	itemsAdded := make([]string, 0)

	// for i < capacity, try adding key-val pair until capacity is full
	for i := 0; i < capacity; i++ {
		key := fmt.Sprintf("key%d", i)
		if arc.capacity > arc.currentlyUsedCapacity+len(key) {
			arc.Set(key, make([]byte, 0))
			itemsAdded = append(itemsAdded, key)
			numItemsAdded++
		}
	}

	if arc.t1.Len() != numItemsAdded {
		t.Errorf("Recently-used cache t1 has wrong length.  Got %v, Expected %v", arc.t1.Len(), numItemsAdded)
		t.FailNow()
	}

	if arc.t2.Len() != 0 {
		t.Errorf("Recently-used cache t2 has wrong length.  Got %v, Expected %v", arc.t2.Len(), 0)
		t.FailNow()
	}

	// Add should upgrade to t2
	for _, key := range itemsAdded {
		arc.Set(key, make([]byte, 0))
	}

	if arc.t1.Len() != 0 {
		t.Errorf("Recently-used cache t1 has wrong length.  Got %v, Expected %v", arc.t1.Len(), 0)
		t.FailNow()
	}

	if arc.t2.Len() != numItemsAdded {
		t.Errorf("Recently-used cache t2 has wrong length.  Got %v, Expected %v", arc.t2.Len(), numItemsAdded)
		t.FailNow()
	}

	// For next set, should remain in t2
	for _, key := range itemsAdded {
		arc.Set(key, make([]byte, 0))
	}

	if arc.t1.Len() != 0 {
		t.Errorf("Recently-used cache t1 has wrong length.  Got %v, Expected %v", arc.t1.Len(), 0)
		t.FailNow()
	}

	if arc.t2.Len() != numItemsAdded {
		t.Errorf("Recently-used cache t2 has wrong length.  Got %v, Expected %v", arc.t2.Len(), numItemsAdded)
		t.FailNow()
	}
}

func TestAdaptiveArc(t *testing.T) {
	capacity := 10240
	arc := NewArc(capacity)

	numItemsAdded := 0
	itemsAdded := make([]string, 0)

	// for i <= capacity, try adding key-val pair until capacity is full
	for i := 0; i <= capacity; i++ {
		key := fmt.Sprintf("key%d", i)
		if arc.capacity > arc.currentlyUsedCapacity+len(key) {
			arc.Set(key, make([]byte, 0))
			itemsAdded = append(itemsAdded, key)
			numItemsAdded++
		}
	}

	//fmt.Printf("%v\n", numItemsAdded)
	//fmt.Printf("%v\n", arc.currentlyUsedCapacity)

	if arc.t1.Len() != numItemsAdded {
		t.Errorf("Recently-used cache t1 has wrong length.  Got %v, Expected %v", arc.t1.Len(), numItemsAdded)
		t.FailNow()
	}

	// Randomly move 50 items to T2
	movedNumber := numItemsAdded / 3
	for i := 0; i < movedNumber; i++ {
		key := fmt.Sprintf("key%d", i)
		arc.Get(key)
	}
	if arc.capacity != capacity {
		t.Errorf("HAHA - you played yourself")
	}

	// make sure the items are successfully moved to t2
	if arc.t2.Len() != movedNumber {
		t.Errorf("Recently-used cache t2 has wrong length.  Got %v, Expected %v", arc.t2.Len(), movedNumber)
		t.FailNow()
	}
	if arc.capacity != capacity {
		t.Errorf("HAHA - you played yourself")
	}

	for i := 200; i < 100000; i++ {
		key := fmt.Sprintf("key%d", i)
		arc.Set(key, make([]byte, 0))
	}

}
