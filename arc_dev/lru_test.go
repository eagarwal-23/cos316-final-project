/******************************************************************************
 * lru_test.go
 * Author:
 * Usage:    `go test`  or  `go test -v`
 * Description:
 *    An incomplete unit testing suite for lru.go. You are welcome to change
 *    anything in this file however you would like. You are strongly encouraged
 *    to create additional tests for your implementation, as the ones provided
 *    here are extremely basic, and intended only to demonstrate how to test
 *    your program.
 ******************************************************************************/

package cache

import (
	"bytes"
	"fmt"
	"testing"
)

/******************************************************************************/
//METHODS TO TEST
// MaxStorage() int - DONE - TESTED AT BEGGINING AFTER INIT

// RemainingStorage() int

// Get(key string) (value []byte, ok bool) - TESTED

// Remove(key string) (value []byte, ok bool)

// Set(key string, value []byte) bool

// Peek(key string) (value []byte, ok bool)

// Empty()

// Len() int

// Stats() *Stats

/******************************************************************************/
/*                                Constants                                   */
/******************************************************************************/
// Constants can go here

/******************************************************************************/
/*                                  Tests                                     */
/******************************************************************************/

//Check that NewFIFO() returns an empty FIFO of the correct size
func TestNewLru(t *testing.T) {
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

// Check that Get() returns no binding when called on an empty FIFO
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

// Check that Peek() returns no binding when called on an empty FIFO
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
	capacity := 16

	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)
	lru.Set("Hello",[]byte("World")) 
	value, ok := lru.Get("Hello")

	if ok {
		res := bytes.Compare(value, []byte("World"))
		if res != 0 {
			t.Errorf("Returned wrong value for key. Got %v, Expected %v", value, []byte("World"))
			t.FailNow()
		}
	} else {
		t.Errorf("Expected value but did not get one")
			t.FailNow()
	}

	maxStorage := lru.MaxStorage()
	if maxStorage == 6 {
		t.Errorf("Returned wrong maxStorage after  for key. Got %v, Expected %v", maxStorage, )
		t.FailNow()
	}
}

// NewLru(16) &LRU{}                                             ...ok
// Set("Hello",'World') true                                               ...ok
// Get("Hello") cache_hit:<'World'>                                ...ok
// MaxStorage() 16                                                 ...ok
// RemainingStorage() 6                                                  ...ok
//      Len() 1                                                  ...ok
// NewLru(64) &LRU{}                                             ...ok
// Set("Foo",'Bar') true                                               ...ok
// Get("Foo") cache_hit:<'Bar'>                                  ...ok
// MaxStorage() 64                                                 ...ok
// RemainingStorage() 58                                                 ...ok
//      Len() 1#01                                               ...ok
// NewLru(256) &LRU{}                                             ...ok
// Set("COS",'316') true                                               ...ok
// Get("COS") cache_hit:<'316'>                                  ...ok
// MaxStorage() 256                                                ...ok
// RemainingStorage() 250                                                ...ok
//      Len() 1#02                                               ...ok

// Check various operations on an LRU with a single binding


func TestLRU(t *testing.T) {
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
	}

	fmt.Println("TEST 52")
	fifo2 := NewLru(100)
	fifo2.Set("____0", []byte("____0"))
	fifo2.Set("____1", []byte("____1"))
	fifo2.Set("____2", []byte("____2"))
	fifo2.Set("____3", []byte("____3"))
	fifo2.Set("____4", []byte("____4"))
	fifo2.Set("____5", []byte("____5"))
	fifo2.Set("____6", []byte("____6"))
	fifo2.Set("____7", []byte("____7"))
	fifo2.Set("____8", []byte("____8"))
	fifo2.Set("____9", []byte("____9"))
	fifo2.Set("____10", []byte("____a"))
	fmt.Println("length", fifo2.Len())
	fmt.Println("remainingcapacity", fifo2.RemainingStorage())
	fmt.Println(fifo2.Get("____0"))

	// fmt.Println("fifo2 capacity = ", fifo2.capacity)
	// fifo2.Set("12345", []byte("12345"))
	// fmt.Println("remainingCap Before inserting = ", fifo2.capacity-fifo2.currentlyUsedCapacity)

	// fifo2.Set("123", []byte("123"))
	// fmt.Println("length", fifo2.Len())

	// fmt.Println("remainingCap After inserting = ", fifo2.capacity-fifo2.currentlyUsedCapacity)

	// fmt.Println("GET 0 FROM END OF LIST")
	// ind := 0
	// key := fmt.Sprintf("key%d", ind)
	// val := []byte(key)
	// res, _ := lru.Get(key)
	// if !bytesEqual(res, val) {
	// 	t.Errorf("Wrong value %s for binding with key: %s", res, key)
	// 	t.FailNow()
	// }

	// fmt.Println("GET 3 FROM MIDDLE OF LIST")
	// ind = 3
	// key = fmt.Sprintf("key%d", ind)
	// val = []byte(key)
	// res, _ = lru.Get(key)
	// if !bytesEqual(res, val) {
	// 	t.Errorf("Wrong value %s for binding with key: %s", res, key)
	// 	t.FailNow()
	// }

	// fmt.Println("IAN TESTING CAPACITY")
	// fmt.Println("usedCap Before inserting = ", lru.currentlyUsedCapacity)
	// fmt.Println("remainingCap Before inserting = ", lru.capacity-lru.currentlyUsedCapacity)
	// lru.Set("H", []byte("hh"))
	// fmt.Println("usedCap after inserting = ", lru.currentlyUsedCapacity)
	// lru.Set("H", []byte("nn"))
	// fmt.Println("usedCap after inserting duplicate same size = ", lru.currentlyUsedCapacity)
	// lru.Set("H", []byte("h"))
	// fmt.Println("usedCap after inserting duplicate smaller = ", lru.currentlyUsedCapacity)
	// lru.Set("H", []byte("hhh"))
	// fmt.Println("usedCap after inserting duplicate bigger = ", lru.currentlyUsedCapacity)
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
