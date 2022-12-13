// package cache

// import (
// 	"container/list"
// 	"fmt"
// )

// type mapping struct {
// 	val  []byte
// 	self *string // for testing purposes
// 	prev *string
// 	next *string
// }

// // An LRU is a fixed-size in-memory cache with least-recently-used eviction
// type LRU struct {
// 	// whatever fields you want here
// 	head                  *string
// 	cachedValues          map[string]mapping
// 	cachedList            list.List
// 	capacity              int
// 	currentlyUsedCapacity int
// 	stats                 Stats
// }

// // NewLRU returns a pointer to a new LRU with a capacity to store limit bytes
// func NewLru(limit int) *LRU {
// 	return &LRU{head: nil, cachedValues: make(map[string]mapping), cachedList: *list.New(), capacity: limit, currentlyUsedCapacity: 0, stats: Stats{}}
// }

// // MaxStorage returns the maximum number of bytes this LRU can store
// func (lru *LRU) MaxStorage() int {
// 	return lru.capacity
// }

// // RemainingStorage returns the number of unused bytes available in this LRU
// func (lru *LRU) RemainingStorage() int {
// 	return lru.capacity - lru.currentlyUsedCapacity
// }

// // Get returns the value associated with the given key, if it exists.
// // This operation counts as a "use" for that key-value pair
// // ok is true if a value was found and false otherwise.
// func (lru *LRU) Get(key string) (value []byte, ok bool) {
// 	currMapping, ok := lru.cachedValues[key]

// 	if ok {
// 		if len(lru.cachedValues) > 1 {
// 			// remove node from current position

// 			var currPrev *string = nil
// 			var currNext *string = nil

// 			// If key is not already at the head
// 			if currMapping.prev != nil {
// 				// Key is in the middle of the list
// 				if currMapping.next != nil {
// 					currPrev = currMapping.prev // address of prev key
// 					currNext = currMapping.next // address of next key

// 					prevMapping := lru.cachedValues[*currPrev]
// 					nextMapping := lru.cachedValues[*currNext]

// 					lru.cachedValues[*currPrev] = mapping{val: prevMapping.val, self: currPrev, prev: prevMapping.prev, next: currNext}
// 					lru.cachedValues[*currNext] = mapping{val: nextMapping.val, self: currNext, prev: currPrev, next: nextMapping.next}
// 				} else { // Key is at the end of the list
// 					currPrev = currMapping.prev // address of prev key
// 					prevMapping := lru.cachedValues[*currPrev]
// 					lru.cachedValues[*currPrev] = mapping{val: prevMapping.val, self: currPrev, prev: prevMapping.prev, next: nil}
// 				}

// 				// Place key at the head of the list
// 				prevHead := lru.head
// 				lru.head = &key
// 				lru.cachedValues[key] = mapping{val: currMapping.val, self: &key, prev: nil, next: prevHead}
// 			}
// 		}

// 		lru.stats.Hits += 1

// 	} else {
// 		lru.stats.Misses += 1
// 	}

// 	fmt.Println("KEY GET: ", key)
// 	fmt.Println("AFTER GET: ", lru.cachedValues)
// 	// THIS WORKS

// 	return currMapping.val, ok
// }

// // Remove removes and returns the value associated with the given key, if it exists.
// // ok is true if a value was found and false otherwise
// func (lru *LRU) Remove(key string) (value []byte, ok bool) {
// 	return nil, false
// }

// // Set associates the given value with the given key, possibly evicting values
// // to make room. Returns true if the binding was added successfully, else false.
// func (lru *LRU) Set(key string, value []byte) bool {
// 	currentObjectSize := len(key) + len(value)
// 	if currentObjectSize > lru.capacity {
// 		return false
// 	}

// 	for lru.capacity-lru.currentlyUsedCapacity < currentObjectSize {
// 		if successfulEvict := lru.Evict(); !successfulEvict {
// 			return false
// 		}
// 	}

// 	// If empty cache, then prev and next are nil
// 	if lru.head == nil {
// 		lru.head = &key
// 		lru.cachedValues[key] = mapping{val: value, self: &key, prev: nil, next: nil}
// 	} else {
// 		// If cache is not empty
// 		prevHead := lru.head                                                                                                    // Address of current head of list
// 		lru.head = &key                                                                                                         // Head of LRU is address of new key
// 		lru.cachedValues[key] = mapping{val: value, self: &key, prev: nil, next: prevHead}                                      // Create new mapping for added node
// 		prevMapping := lru.cachedValues[*prevHead]                                                                              // Mapping for previous head
// 		lru.cachedValues[*prevHead] = mapping{val: prevMapping.val, self: prevMapping.self, prev: &key, next: prevMapping.next} // Updating prev for previous head to be address of added key
// 	}

// 	// Increase currentlyUsedCapacity to reflect currentObjectSize
// 	lru.currentlyUsedCapacity += currentObjectSize

// 	fmt.Println("AFTER SET: ", lru.cachedValues)
// 	return true
// }

// func (lru *LRU) Evict() bool {

// 	// // Last element of list
// 	// elem := (lru.cachedList).Back()
// 	// key := elem.Value.(string)

// 	// // If key is found
// 	// if value, ok := lru.cachedValues[key]; ok {
// 	// 	// Delete from map
// 	// 	delete(lru.cachedValues, key)

// 	// 	// Delete from list
// 	// 	lru.cachedList.Remove(elem)

// 	// 	// Decrease currentlyUsedCapacity to reflect deletion
// 	// 	currentObjectSize := len(key) + len(value)
// 	// 	lru.currentlyUsedCapacity -= currentObjectSize

// 	// 	// Eviction successful
// 	// 	return true
// 	// }
// 	// return false
// 	return false
// }

// // Len returns the number of bindings in the LRU.
// func (lru *LRU) Len() int {
// 	return len(lru.cachedValues)
// }

// // Stats returns statistics about how many search hits and misses have occurred.
// func (lru *LRU) Stats() *Stats {
// 	return &Stats{}
// }

// /*
// SOURCES

// https://www.youtube.com/watch?v=S6IfqDXWa10
// */
