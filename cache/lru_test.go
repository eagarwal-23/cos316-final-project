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
