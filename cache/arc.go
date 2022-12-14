package cache

// import (
// 	"container/list"
// )

// // An ARC is a fixed-size in-memory cache with least-recently-used eviction
// type ARC struct {
// 	p					  int				 // P is the dynamic preference towards T1 or T2
// 	capacity              int                // To hold the capacity of the cache
// 	t1				      LRU				 // 			 
// 	currentlyUsedCapacity int                // Currently used capacity of the cache
// 	stats                 Stats              // Hits and misses for the cache
// }

// // NewArc returns a pointer to a new ARC with a capacity to store limit bytes
// func NewArc(limit int) *ARC {
// 	return &ARC{cachedValues: make(map[string]mapping), cachedList: *list.New(), capacity: limit, currentlyUsedCapacity: 0, stats: Stats{}}
// }

// // MaxStorage returns the maximum number of bytes this  can store
// func (arc *ARC) MaxStorage() int {
// 	return arc.capacity
// }

// // RemainingStorage returns the number of unused bytes available in this ARC
// func (arc *ARC) RemainingStorage() int {
// 	return arc.capacity - arc.currentlyUsedCapacity
// }

// // Peek returns the value associated with the given key, if it exists.
// // This operation does not counts as a "use" for that key-value pair
// // ok is true if a value was found and false otherwise.
// func (arc *ARC) Peek(key string) (value []byte, ok bool) {
// 	currMapping, ok := arc.cachedValues[key]
// 	return currMapping.value, ok
// }

// // Get returns the value associated with the given key, if it exists.
// // This operation counts as a "use" for that key-value pair
// // ok is true if a value was found and false otherwise.
// func (arc *ARC) Get(key string) (value []byte, ok bool) {
// 	currMapping, ok := arc.cachedValues[key]

// 	if ok {
// 		if len(arc.cachedValues) > 1 {
// 			//if key exists move to front of linkedlist
// 			arc.cachedList.MoveToFront(currMapping.Node)
// 		}
// 		arc.stats.Hits += 1
// 	} else {
// 		arc.stats.Misses += 1
// 	}

// 	// THIS WORKS

// 	return currMapping.value, ok
// }

// // Remove removes and returns the value associated with the given key, if it exists.
// // ok is true if a value was found and false otherwise
// func (arc *ARC) Remove(key string) (value []byte, ok bool) {
// 	currMapping, ok := arc.cachedValues[key]

// 	if !ok {
// 		return nil, false
// 	}

// 	delete(arc.cachedValues, key)
// 	arc.cachedList.Remove(currMapping.Node)
// 	arc.currentlyUsedCapacity -= len(currMapping.key) + len(currMapping.value)

// 	return currMapping.value, ok
// }

// // Set associates the given value with the given key, possibly evicting values
// // to make room. Returns true if the binding was added successfully, else false.
// func (arc *ARC) Set(key string, value []byte) bool {
// 	// check if key exists - simply replace
// 	currMapping, ok := arc.cachedValues[key]

// 	currentObjectSize := len(key) + len(value)
// 	// If objectSize is larger than the whole cache
// 	if currentObjectSize > arc.capacity {
// 		return false
// 	}

// 	if ok {
// 		if len(currMapping.value) == len(value) {
// 			currMapping.value = value
// 			arc.cachedValues[key] = currMapping
// 		} else if len(currMapping.value) > len(value) {
// 			arc.currentlyUsedCapacity -= len(currMapping.value)
// 			currMapping.value = value
// 			arc.cachedValues[key] = currMapping
// 			arc.currentlyUsedCapacity += len(value)
// 		} else if len(currMapping.value) < len(value) {
// 			arc.currentlyUsedCapacity -= len(currMapping.value)
// 			if arc.capacity-arc.currentlyUsedCapacity < len(value) {
// 				if successfulEvict := arc.Evict(); !successfulEvict {
// 					return false
// 				}
// 			}
// 			currMapping.value = value
// 			arc.cachedValues[key] = currMapping
// 			arc.currentlyUsedCapacity += len(value)
// 		}
// 		return true
// 	}

// 	if arc.capacity-arc.currentlyUsedCapacity < currentObjectSize {
// 		if successfulEvict := arc.Evict(); !successfulEvict {
// 			return false
// 		}
// 	}

// 	currMapping = mapping{key: key, value: value, Node: nil}
// 	if arc.cachedList.Len() == 0 {
// 		arc.cachedList.Init()
// 	}
// 	mappingPtr := arc.cachedList.PushFront(currMapping)

// 	currMapping.Node = mappingPtr
// 	mappingPtr.Value = currMapping
// 	arc.cachedValues[key] = currMapping

// 	// Increase currentlyUsedCapacity to reflect currentObjectSize
// 	arc.currentlyUsedCapacity += currentObjectSize
// 	return true
// }

// func (arc *ARC) Empty() {
// 	arc.cachedValues = make(map[string]mapping)
// 	arc.cachedList = *list.New()
// 	arc.currentlyUsedCapacity = 0
// }

// func (arc *ARC) Evict() bool {
// 	// Last element of list

// 	elem := (arc.cachedList).Back()
// 	currMapping := elem.Value.(mapping)

// 	// If key is not found
// 	if _, ok := arc.cachedValues[currMapping.key]; !ok {
// 		return false
// 	}

// 	// If key is found - Delete from map
// 	delete(arc.cachedValues, currMapping.key)

// 	// Delete from list
// 	arc.cachedList.Remove(elem)

// 	// Decrease currentlyUsedCapacity to reflect deletion
// 	currentObjectSize := len(currMapping.key) + len(currMapping.value)
// 	arc.currentlyUsedCapacity -= currentObjectSize

// 	// Eviction successful
// 	return true
// }

// // Len returns the number of bindings in the ARC.
// func (arc *ARC) Len() int {
// 	return len(arc.cachedValues)
// }

// // Stats returns statistics about how many search hits and misses have occurred.
// func (arc *ARC) Stats() *Stats {
// 	return &arc.stats
// }

// func (arc *ARC) resizeARC(length int) {
// 	for arc.RemainingStorage() < length {
// 		arc.Evict()
// 	}
// }

// /*
// SOURCES

// https://www.youtube.com/watch?v=S6IfqDXWa10
// */
