package cache

import (
	"container/list"
)

type mapping struct {
	key   string
	value []byte
	Node  *list.Element
}

// An LRU is a fixed-size in-memory cache with least-recently-used eviction
type LRU struct {
	// whatever fields you want here
	cachedValues          map[string]mapping // Map containing key-value pairings
	cachedList            list.List          // Linked list to hold insertion order
	capacity              int                // To hold the capacity of the cache
	currentlyUsedCapacity int                // Currently used capacity of the cache
	stats                 Stats              // Hits and misses for the cache
}

// NewLRU returns a pointer to a new LRU with a capacity to store limit bytes
func NewLru(limit int) *LRU {
	return &LRU{cachedValues: make(map[string]mapping), cachedList: *list.New(), capacity: limit, currentlyUsedCapacity: 0, stats: Stats{}}
}

// MaxStorage returns the maximum number of bytes this LRU can store
func (lru *LRU) MaxStorage() int {
	return lru.capacity
}

// RemainingStorage returns the number of unused bytes available in this LRU
func (lru *LRU) RemainingStorage() int {
	return lru.capacity - lru.currentlyUsedCapacity
}

// Peek returns the value associated with the given key, if it exists.
// This operation does not counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (lru *LRU) Peek(key string) (value []byte, ok bool) {
	currMapping, ok := lru.cachedValues[key]
	return currMapping.value, ok
}

// Get returns the value associated with the given key, if it exists.
// This operation counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (lru *LRU) Get(key string) (value []byte, ok bool) {
	currMapping, ok := lru.cachedValues[key]

	if ok {
		if len(lru.cachedValues) > 1 {
			//if key exists move to front of linkedlist
			lru.cachedList.MoveToFront(currMapping.Node)
		}
		lru.stats.Hits += 1
	} else {
		lru.stats.Misses += 1
	}

	// THIS WORKS

	return currMapping.value, ok
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lru *LRU) Remove(key string) (value []byte, ok bool) {
	currMapping, ok := lru.cachedValues[key]

	if !ok {
		return nil, false
	}

	delete(lru.cachedValues, key)
	lru.cachedList.Remove(currMapping.Node)
	lru.currentlyUsedCapacity -= len(currMapping.key) + len(currMapping.value)

	return currMapping.value, ok
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (lru *LRU) Set(key string, value []byte) bool {
	// check if key exists - simply replace
	currMapping, ok := lru.cachedValues[key]

	currentObjectSize := len(key) + len(value)
	// If objectSize is larger than the whole cache
	if currentObjectSize > lru.capacity {
		return false
	}

	if ok {
		if len(currMapping.value) == len(value) {
			currMapping.value = value
			lru.cachedValues[key] = currMapping
		} else if len(currMapping.value) > len(value) {
			lru.currentlyUsedCapacity -= len(currMapping.value)
			currMapping.value = value
			lru.cachedValues[key] = currMapping
			lru.currentlyUsedCapacity += len(value)
		} else if len(currMapping.value) < len(value) {
			lru.currentlyUsedCapacity -= len(currMapping.value)
			if lru.capacity-lru.currentlyUsedCapacity < len(value) {
				if successfulEvict := lru.Evict(); !successfulEvict {
					return false
				}
			}
			currMapping.value = value
			lru.cachedValues[key] = currMapping
			lru.currentlyUsedCapacity += len(value)
		}
		return true
	}

	if lru.capacity-lru.currentlyUsedCapacity < currentObjectSize {
		if successfulEvict := lru.Evict(); !successfulEvict {
			return false
		}
	}

	currMapping = mapping{key: key, value: value, Node: nil}
	if lru.cachedList.Len() == 0 {
		lru.cachedList.Init()
	}
	mappingPtr := lru.cachedList.PushFront(currMapping)

	currMapping.Node = mappingPtr
	mappingPtr.Value = currMapping
	lru.cachedValues[key] = currMapping

	// Increase currentlyUsedCapacity to reflect currentObjectSize
	lru.currentlyUsedCapacity += currentObjectSize
	return true
}

func (lru *LRU) Empty() {
	lru.cachedValues = make(map[string]mapping)
	lru.cachedList = *list.New()
	lru.currentlyUsedCapacity = 0
}

func (lru *LRU) Evict() bool {
	// Last element of list

	elem := (lru.cachedList).Back()
	currMapping := elem.Value.(mapping)

	// If key is not found
	if _, ok := lru.cachedValues[currMapping.key]; !ok {
		return false
	}

	// If key is found - Delete from map
	delete(lru.cachedValues, currMapping.key)

	// Delete from list
	lru.cachedList.Remove(elem)

	// Decrease currentlyUsedCapacity to reflect deletion
	currentObjectSize := len(currMapping.key) + len(currMapping.value)
	lru.currentlyUsedCapacity -= currentObjectSize

	// Eviction successful
	return true
}

// Len returns the number of bindings in the LRU.
func (lru *LRU) Len() int {
	return len(lru.cachedValues)
}

// Stats returns statistics about how many search hits and misses have occurred.
func (lru *LRU) Stats() *Stats {
	return &lru.stats
}

func (lru *LRU) resizeLRU(length int) {
	for lru.RemainingStorage() < length {
		lru.Evict()
	}
}

/*
SOURCES

https://www.youtube.com/watch?v=S6IfqDXWa10
*/
