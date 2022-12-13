/******************************************************************************
 * fifo_test.go
 * Author:
 * Usage:    `go test`  or  `go test -v`
 * Description:
 *    An incomplete unit testing suite for fifo.go. You are welcome to change
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

func TestFIFO(t *testing.T) {
	capacity := 64
	fifo := NewFifo(capacity)
	checkCapacity(t, fifo, capacity)
	fmt.Println("Capacity Init = ", fifo.capacity)
	fmt.Println("usedCap Init = ", fifo.currentlyUsedCapacity)

	for i := 0; i < 250; i++ {
		key := fmt.Sprintf("key%d", i)
		val := []byte(key)
		ok := fifo.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}

		res, _ := fifo.Get(key)
		if !bytesEqual(res, val) {
			t.Errorf("Wrong value %s for binding with key: %s", res, key)
			t.FailNow()
		}
	}

	// fmt.Println("usedCap Before inserting = ", fifo.currentlyUsedCapacity)
	// fmt.Println("remainingCap Before inserting = ", fifo.capacity - fifo.currentlyUsedCapacity)
	// fifo.Set("H", []byte("hh"))
	// fmt.Println("usedCap after inserting = ", fifo.currentlyUsedCapacity)
	// fifo.Set("H", []byte("nn"))
	// fmt.Println("usedCap after inserting duplicate same size = ", fifo.currentlyUsedCapacity)
	// fifo.Set("H", []byte("h"))
	// fmt.Println("usedCap after inserting duplicate smaller = ", fifo.currentlyUsedCapacity)
	// fifo.Set("H", []byte("hhh"))
	// fmt.Println("usedCap after inserting duplicate bigger = ", fifo.currentlyUsedCapacity)
	// fifo.Set("Nello", []byte("h"))
	// fmt.Println("Capacity = ", fifo.currentlyUsedCapacity)
	// fifo.Set("Fello", []byte("h"))
	// fmt.Println("Capacity = ", fifo.currentlyUsedCapacity)

	
	fmt.Println("TEST 29")
	fifo2 := NewFifo(20)
	fmt.Println("fifo2 capacity = ", fifo2.capacity)  
	fifo2.Set("abcd", []byte("efgh"))
	fifo2.Set("1234", []byte("efgh"))                                       
	fmt.Println("remainingCap Before inserting = ", fifo2.capacity - fifo2.currentlyUsedCapacity)                                              

	fifo2.Set("1234", []byte("12345678"))     
	str1, _ := fifo2.Get("1234")    
	fmt.Println("inserted", str1)                                      

	str, _ := fifo2.Get("abcd")   
	fmt.Println("retrieved", str)    
}
