/******************************************************************************
 * lru_test.go
 * Author:
 * Usage:    `go test`  or  `go test -v`
 * Description:
 *    An incomplete unit testing suite for lru.go.
 ******************************************************************************/

package cache

import (
	"bytes"
	"fmt"
	"testing"
)

/******************************************************************************/
/*                                Constants                                   */
/******************************************************************************/
// Constants can go here

/******************************************************************************/
/*                                  Tests                                     */
/******************************************************************************/

// Check that NewLru() returns an empty LRU of the correct size
func TestNewLru(t *testing.T) {
	capacityArray := [4]int{16, 32, 64, 128}

	for capacity := range capacityArray {
		lru := NewLru(capacity)
		checkCapacity(t, lru, capacity)

		// Len() = 0 on init
		length := lru.Len()
		if length != 0 {
			t.Errorf("NewLru returned wrong length on init. Got %v, Expected %v", length, 0)
			t.FailNow()
		}

		// MaxStorage() = 64 on init
		maxStorage := lru.MaxStorage()
		if maxStorage != capacity {
			t.Errorf("NewLru returned wrong maxStorage on init. Got %v, Expected %v", capacity, maxStorage)
			t.FailNow()
		}

		// RemainingStorage() = 64 on init
		remainingStorage := lru.RemainingStorage()
		if remainingStorage != capacity {
			t.Errorf("NewLru returned wrong remainingStorage on init. Got %v, Expected %v", capacity, remainingStorage)
			t.FailNow()
		}
	}
}

// Check that Get() returns no binding when called on an empty LRU
func TestGetEmptyLru(t *testing.T) {
	capacity := 1024
	keysArray := [4]string{"Hello", "a", "ssup"}

	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	for _, key := range keysArray {
		value, _ := lru.Get(key)
		if value != nil {
			t.Errorf("Returned wrong value for empty LRU. Got %v, Expected %v", value, nil)
			t.FailNow()
		}
	}
}

// Check that Peek() returns no binding when called on an empty LRU
func TestPeekEmptyLru(t *testing.T) {
	capacity := 1024
	keysArray := [4]string{"Hello", "a", "ssup"}

	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	for _, key := range keysArray {
		value, _ := lru.Peek(key)
		if value != nil {
			t.Errorf("Returned wrong value for empty LRU. Got %v, Expected %v", value, nil)
			t.FailNow()
		}
	}
}

// Check various operations on an LRU with a single binding
func TestSingleBindingLru(t *testing.T) {
	capacitiesArray := [3]int{16, 64, 256}
	keysArray := [3]string{"Hello", "Foo", "COS"}
	valuesArray := [3]string{"World", "Bar", "316"}

	for i, _ := range keysArray {
		lru := NewLru(capacitiesArray[i])
		checkCapacity(t, lru, capacitiesArray[i])
		lru.Set(keysArray[i], []byte(valuesArray[i]))
		value, ok := lru.Get(keysArray[i])

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

		remainingStorage := lru.RemainingStorage()
		expectedremainingStorage := capacitiesArray[i] - (len(keysArray[i]) + len([]byte(valuesArray[i])))
		if remainingStorage != expectedremainingStorage {
			t.Errorf("Returned wrong remaining after  for key. Got %v, Expected %v", remainingStorage, expectedremainingStorage)
			t.FailNow()
		}
	}

}

// Add 20 bindings to an LRU, checking each one consumes the right storage
func TestStorageLru(t *testing.T) {
	capacity := 1024
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	for i := 0; i < 20; i++ {
		remainingStorageBefore := lru.RemainingStorage()
		key := fmt.Sprintf("key%d", i)
		val := []byte(key)
		ok := lru.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}
		remainingStorageAfter := lru.RemainingStorage()

		expectedremainingStorageAfter := remainingStorageBefore - (len(key) + len(val))
		if remainingStorageAfter != expectedremainingStorageAfter {
			t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", remainingStorageAfter, expectedremainingStorageAfter)
			t.FailNow()
		}
	}

	// fmt.Printf("UsedBefore: %v\n", arc.capacity)
	// for i := 200; i < 210; i++ {
	// 	key := fmt.Sprintf("key%d", i)
	// 	arc.Set(key, make([]byte, 0))
	// 	fmt.Printf("CurrentlyUsed: %v\n", arc.currentlyUsedCapacity)
	// }

}

// Check that Set() adds bindings to a 'full' LRU by evicting old ones
func TestSetFullLru(t *testing.T) {
	capacity := 30
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	lru.Set("____0", []byte("____0"))
	lru.Set("____1", []byte("____1"))
	lru.Set("____2", []byte("____2"))
	lru.Set("____3", []byte("____3"))
	//len
	len := lru.Len()
	if len != 3 {
		t.Errorf("Len wrong after adding binding to full LRU. Got %v, Expected %v", len, 3)
		t.FailNow()
	}
	lru.Set("____4", []byte("____4"))
	//len
	len = lru.Len()
	if len != 3 {
		t.Errorf("Len wrong after adding binding to full LRU. Got %v, Expected %v", len, 3)
		t.FailNow()
	}
	lru.Set("____5", []byte("____5"))
	//len
	len = lru.Len()
	if len != 3 {
		t.Errorf("Len wrong after adding binding to full LRU. Got %v, Expected %v", len, 3)
		t.FailNow()
	}

}

// Check that Set() rejects bindings too large for the LRU
func TestSetTooLargeLru(t *testing.T) {
	capacity := 10
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	ok := lru.Set("123456", []byte("123456"))
	if ok {
		t.Errorf("Failed to reject binding too large for LRU. Set  Got %v, Expected %v", ok, false)
		t.FailNow()
	}
	if lru.RemainingStorage() != 10 {
		t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", lru.RemainingStorage(), 10)
		t.FailNow()
	}
	_, ok = lru.Get("123456")
	if ok {
		t.Errorf("Failed to reject binding too large for LRU. Set  Got %v, Expected %v", ok, false)
		t.FailNow()
	}
}

// Check that Set() only allows zero-size bindings in a zero-capacity LRU

func TestSetZeroLru(t *testing.T) {
	capacity := 0
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	ok := lru.Set("hello", []byte("world"))
	if ok {
		t.Errorf("Failed to reject binding too large for LRU. Set  Got %v, Expected %v", ok, false)
		t.FailNow()
	}
	ok = lru.Set("foo", []byte("boo"))
	if ok {
		t.Errorf("Failed to reject binding too large for LRU. Set  Got %v, Expected %v", ok, false)
		t.FailNow()
	}
	ok = lru.Set("", []byte(""))
	if !ok {
		t.Errorf("Failed to reject binding too large for LRU. Set  Got %v, Expected %v", ok, true)
		t.FailNow()
	}
}

// Check that the LRU allows the empty string as a valid key
func TestEmptyStringValidLru(t *testing.T) {
	capacity := 1024
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	ok := lru.Set("", []byte("Value"))
	if lru.RemainingStorage() != 1019 {
		t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", lru.RemainingStorage(), 1019)
		t.FailNow()
	}
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	}
	if lru.Len() != 1 {
		t.Errorf("Len() wrong. Got %v, Expected %v", lru.Len(), 1)
		t.FailNow()
	}
	if lru.MaxStorage() != capacity {
		t.Errorf("MaxStorage wrong. Got %v, Expected %v", lru.MaxStorage(), capacity)
		t.FailNow()
	}

	value, ok := lru.Get("")
	res := bytes.Compare(value, []byte("Value"))
	if res != 0 {
		t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte("Value"))
		t.FailNow()
	}
}

// Check that the LRU allows the empty []byte as a valid value
func TestEmptyValidLru(t *testing.T) {
	capacity := 1024
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	ok := lru.Set("key", []byte{})
	if lru.RemainingStorage() != 1021 {
		t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", lru.RemainingStorage(), 1021)
		t.FailNow()
	}
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	}
	if lru.Len() != 1 {
		t.Errorf("Len() wrong. Got %v, Expected %v", lru.Len(), 1)
		t.FailNow()
	}
	if lru.MaxStorage() != capacity {
		t.Errorf("MaxStorage wrong. Got %v, Expected %v", lru.MaxStorage(), capacity)
		t.FailNow()
	}

	value, ok := lru.Get("key")
	res := bytes.Compare(value, []byte{})
	if res != 0 {
		t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte{})
		t.FailNow()
	}

}

// Check that the LRU allows the empty []byte as a valid value
func TestNilValidLru(t *testing.T) {
	capacity := 1024
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	ok := lru.Set("key", nil)
	if lru.RemainingStorage() != 1021 {
		t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", lru.RemainingStorage(), 1021)
		t.FailNow()
	}
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	}
	if lru.Len() != 1 {
		t.Errorf("Len() wrong. Got %v, Expected %v", lru.Len(), 1)
		t.FailNow()
	}
	if lru.MaxStorage() != capacity {
		t.Errorf("MaxStorage wrong. Got %v, Expected %v", lru.MaxStorage(), capacity)
		t.FailNow()
	}

	value, ok := lru.Get("key")
	res := bytes.Compare(value, nil)
	if res != 0 {
		t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, nil)
		t.FailNow()
	}

}

// Check that values can be non-ASCII (binary)
func TestBinaryValuesLru(t *testing.T) {
	capacity := 1024
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	ok := lru.Set("key", []byte("\x00\x01???\x15???"))
	if lru.RemainingStorage() != 1012 {
		t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", lru.RemainingStorage(), 1012)
		t.FailNow()
	}
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	}
	if lru.Len() != 1 {
		t.Errorf("Len() wrong. Got %v, Expected %v", lru.Len(), 1)
		t.FailNow()
	}
	if lru.MaxStorage() != capacity {
		t.Errorf("MaxStorage wrong. Got %v, Expected %v", lru.MaxStorage(), capacity)
		t.FailNow()
	}

	value, ok := lru.Get("key")
	res := bytes.Compare(value, []byte("\x00\x01???\x15???"))
	if res != 0 {
		t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte("\x00\x01???\x15???"))
		t.FailNow()
	}

}

// Check that keys and values can be non-ASCII (Unicode)
func TestUnicodeValuesLru(t *testing.T) {
	capacity := 1024
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	ok := lru.Set("????_????", []byte("???_????"))
	if lru.RemainingStorage() != 1007 {
		t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", lru.RemainingStorage(), 1007)
		t.FailNow()
	}
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	}
	if lru.Len() != 1 {
		t.Errorf("Len() wrong. Got %v, Expected %v", lru.Len(), 1)
		t.FailNow()
	}
	if lru.MaxStorage() != capacity {
		t.Errorf("MaxStorage wrong. Got %v, Expected %v", lru.MaxStorage(), capacity)
		t.FailNow()
	}

	value, ok := lru.Get("????_????")
	res := bytes.Compare(value, []byte("???_????"))
	if res != 0 {
		t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte("???_????"))
		t.FailNow()
	}

}

// Test that Set() overwrites values when called with an existing key
func TestSetOverwriteLru(t *testing.T) {
	capacity := 1024
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	ok := lru.Set("key", []byte("old"))
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	} else {
		value, _ := lru.Get("key")
		res := bytes.Compare(value, []byte("old"))
		if res != 0 {
			t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte("old"))
			t.FailNow()
		}
	}

	ok = lru.Set("key", []byte("new"))
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	} else {
		value, _ := lru.Get("key")
		res := bytes.Compare(value, []byte("new"))
		if res != 0 {
			t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte("new"))
			t.FailNow()
		}

	}
}

// Test that Set() overwrites values when called with an existing key
func TestSetOverwriteStorageLru(t *testing.T) {
	capacity := 1024
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	ok := lru.Set("key", []byte("old"))
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	} else {
		value, _ := lru.Get("key")
		res := bytes.Compare(value, []byte("old"))
		if res != 0 {
			t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte("old"))
			t.FailNow()
		}

		if lru.RemainingStorage() != 1018 {
			t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", lru.RemainingStorage(), 1018)
			t.FailNow()
		}
	}

	ok = lru.Set("key", []byte("nw"))
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	} else {
		value, _ := lru.Get("key")
		res := bytes.Compare(value, []byte("nw"))
		if res != 0 {
			t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte("nw"))
			t.FailNow()
		}
		if lru.RemainingStorage() != 1019 {
			t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", lru.RemainingStorage(), 1019)
			t.FailNow()
		}

	}
}

// Check that Remove() prevents Get() from retrieving a binding
func TestRemovePreventGetLru(t *testing.T) {
	capacity := 1024
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	ok := lru.Set("key", []byte("value"))
	if !ok {
		t.Errorf("Failed to add  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	} else {
		value, _ := lru.Get("key")
		res := bytes.Compare(value, []byte("value"))
		if res != 0 {
			t.Errorf("Fetched wrong value. Set  Got %v, Expected %v", value, []byte("value"))
			t.FailNow()
		}
	}

	_, ok = lru.Remove("key")
	if !ok {
		t.Errorf("Failed to remove  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	} else {
		_, ok := lru.Get("key")
		if ok {
			t.Errorf("Fetched a removed value. Set  Got %v, Expected %v", ok, false)
			t.FailNow()
		}
	}

}

// Check that Remove() correctly updates available storage
func TestRemoveStorageLru(t *testing.T) {
	capacity := 1024
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	lru.Set("__0", []byte("__0"))
	lru.Set("__1", []byte("__1"))
	lru.Set("__2", []byte("__2"))
	lru.Set("__3", []byte("__3"))
	lru.Remove("__0")

	//len
	len := lru.Len()
	if len != 3 {
		t.Errorf("Len wrong after removing. Got %v, Expected %v", len, 3)
		t.FailNow()
	}
	if lru.RemainingStorage() != 1006 {
		t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", lru.RemainingStorage(), 1006)
		t.FailNow()
	}
	lru.Remove("__1")
	len = lru.Len()
	//len
	if len != 2 {
		t.Errorf("Len wrong after removing. Got %v, Expected %v", len, 2)
		t.FailNow()
	}
	if lru.RemainingStorage() != 1012 {
		t.Errorf("RemainingStorage wrong after adding binding. Got %v, Expected %v", lru.RemainingStorage(), 1012)
		t.FailNow()
	}
}

// Check that Stats() returns correct values when there are mixed cache hits and misses
func TestStatsLru(t *testing.T) {
	capacity := 1024
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	lru.Set("____1", []byte("____1"))
	_, ok := lru.Get("____1")

	if !ok {
		t.Errorf("Failed to fetch  binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	}

	_, ok = lru.Get("miss")

	if ok {
		t.Errorf("Fetched absent binding binding. Got %v, Expected %v", ok, false)
		t.FailNow()
	}

	hits, misses := lru.Stats().Hits, lru.stats.Misses

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
func TestRemoveOverwrittenLru(t *testing.T) {
	capacity := 1024
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	lru.Set("key", []byte("old"))
	lru.Set("key", []byte("newval"))
	lru.Remove("key")
	_, ok := lru.Get("key")
	if ok {
		t.Errorf("Failed to remove binding with key: 'key'. Got %v, Expected %v", ok, false)
		t.FailNow()
	}

}

// Check that Remove() has no effect when called on an empty LRU
func TestRemoveEmptyLru(t *testing.T) {
	capacity := 1024
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	_, ok := lru.Remove("key")
	if ok {
		t.Errorf("Removed empty binding. Got %v, Expected %v", ok, false)
		t.FailNow()
	}
	_, ok = lru.Remove("foo")
	if ok {
		t.Errorf("Removed empty binding. Got %v, Expected %v", ok, false)
		t.FailNow()
	}
	_, ok = lru.Remove("bar")
	if ok {
		t.Errorf("Removed empty binding. Got %v, Expected %v", ok, false)
		t.FailNow()
	}
}

// Attempt to Remove() a binding that has already been removed
func TestRemoveRemovedLru(t *testing.T) {
	capacity := 1024
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	lru.Set("key", []byte("value"))
	_, ok := lru.Remove("key")
	if !ok {
		t.Errorf("Failed to remove binding. Got %v, Expected %v", ok, true)
		t.FailNow()
	}
	_, ok = lru.Remove("key")
	if ok {
		t.Errorf("Removed removed binding. Got %v, Expected %v", ok, false)
		t.FailNow()
	}

}

func TestLRU_Peek(t *testing.T) {
	capacity := 64
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	for i := 0; i < 15; i++ {
		key := fmt.Sprintf("key%d", i)
		val := []byte(key)
		ok := lru.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}

		res, _ := lru.Get(key)
		if !bytesEqual(res, val) {
			t.Errorf("Wrong value %s for binding with key: %s", res, key)
			t.FailNow()
		}

		// Peek value
		hits_before_access := lru.stats.Hits
		misses_before_access := lru.stats.Misses
		res2, _ := lru.Peek(key)
		hits_after_access := lru.stats.Hits
		misses_after_access := lru.stats.Misses

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
