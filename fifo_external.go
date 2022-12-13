// package cache

// import (
// 	"github.com/emirpasic/gods/maps/linkedhashmap"
// )

// /*
// cursor parking lot
// 🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳
// 🌳                                                     🌳
// 🌳     🐝           🦝                 🌲               🌳
// 🌳                                                     🌳
// 🌳              🦋                🦩         🐕         🌳
// 🌳                                                     🌳
// 🌳      🐐            🌲                     🐧         🌳
// 🌳                                                     🌳
// 🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳🌳

// */

// // An FIFO is a fixed-size in-memory cache with first-in first-out eviction
// type OLD_FIFO struct {
// 	// Linked HashMap to store key-value pairs and preserve order of insertion
// 	cachedValues          *linkedhashmap.Map
// 	capacity              int
// 	currentlyUsedCapacity int
// 	stats                 Stats
// }

// // NewFIFO returns a pointer to a new FIFO with a capacity to store limit bytes
// func NewFifo(limit int) *FIFO {
// 	return &FIFO{cachedValues: linkedhashmap.New(), capacity: limit, currentlyUsedCapacity: limit, stats: Stats{}}
// }

// // MaxStorage returns the maximum number of bytes this FIFO can store
// func (fifo *FIFO) MaxStorage() int {
// 	return fifo.capacity
// }

// // RemainingStorage returns the number of unused bytes available in this FIFO
// func (fifo *FIFO) RemainingStorage() int {
// 	return fifo.capacity - fifo.currentlyUsedCapacity
// }

// // Get returns the value associated with the given key, if it exists.
// // ok is true if a value was found and false otherwise.
// func (fifo *FIFO) Get(key string) (value []byte, ok bool) {
// 	valueInterface, ok := fifo.cachedValues.Get(key)
// 	// value = value1.([]byte)
// 	if !ok {
// 		fifo.stats.Misses += 1
// 	} else {
// 		fifo.stats.Hits += 1
// 	}
// 	return valueInterface.([]byte), ok
// }

// // Remove removes and returns the value associated with the given key, if it exists.
// // ok is true if a value was found and false otherwise
// func (fifo *FIFO) Remove(key string) (value []byte, ok bool) {
// 	valueInterface, ok := fifo.cachedValues.Get(key)
// 	if !ok {
// 		return nil, false
// 	}
// 	fifo.cachedValues.Remove(key)
// 	return valueInterface.([]byte), ok
// }

// // Set associates the given value with the given key, possibly evicting values
// // to make room. Returns true if the binding was added successfully, else false.
// func (fifo *FIFO) Set(key string, value []byte) bool {

// 	currentObjectSize := len(key) + len(value)

// 	// If insufficient capacity, keep evicting elements till
// 	// enough space
// 	for fifo.capacity-fifo.currentlyUsedCapacity < currentObjectSize {
// 		fifo.Evict()
// 	}

// 	// Else, sufficient capacity, hence insert
// 	fifo.cachedValues.Put(key, value)
// 	fifo.currentlyUsedCapacity = +currentObjectSize

// 	return false
// }

// func (fifo *FIFO) Evict() {
// 	iterate := fifo.cachedValues.Iterator()
// 	var currentObjectSize int
// 	if iterate.First() {
// 		key := iterate.Key().(string)
// 		value, _ := fifo.Remove(key)
// 		currentObjectSize = len(key) + len(value)
// 	}
// 	fifo.currentlyUsedCapacity = -currentObjectSize
// }

// // Len returns the number of bindings in the FIFO.
// func (fifo *FIFO) Len() int {
// 	return fifo.cachedValues.Size()
// }

// // Stats returns statistics about how many search hits and misses have occurred.
// func (fifo *FIFO) Stats() *Stats {
// 	return &fifo.stats
// }
