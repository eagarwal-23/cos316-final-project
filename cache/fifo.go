package cache

type FIFO struct {
	// Linked HashMap to store key-value pairs and preserve order of insertion
	cachedValues          map[string][]byte // Map containing key-value pairings
	cachedList            []string          // Slice to hold insertion order
	capacity              int               // To hold the capacity of the cache
	currentlyUsedCapacity int               // Currently used capacity of the cache
	stats                 Stats             // Hits and misses for the cache
}

// NewFIFO returns a pointer to a new FIFO with a capacity to store limit bytes
func NewFifo(limit int) *FIFO {
	return &FIFO{cachedValues: make(map[string][]byte), cachedList: make([]string, 0), capacity: limit, currentlyUsedCapacity: 0, stats: Stats{}}
}

// MaxStorage returns the maximum number of bytes this FIFO can store
func (fifo *FIFO) MaxStorage() int {
	return fifo.capacity
}

// RemainingStorage returns the number of unused bytes available in this FIFO
func (fifo *FIFO) RemainingStorage() int {
	return fifo.capacity - fifo.currentlyUsedCapacity
}

// Get returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise.
func (fifo *FIFO) Get(key string) (value []byte, ok bool) {
	valueInterface, ok := fifo.cachedValues[key]

	// Calculating hits and misses
	if !ok {
		fifo.stats.Misses += 1
	} else {
		fifo.stats.Hits += 1
	}
	return valueInterface, ok
}

// Given a list and a string, return the index of the string if found, else -1
func searchIndex(keys []string, key string) int {
	for index, _ := range keys {
		if keys[index] == key {
			return index
		}
	}
	return -1
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (fifo *FIFO) Remove(key string) (value []byte, ok bool) {

	// If value is found
	if value, ok := fifo.cachedValues[key]; ok {
		// Remove from map
		delete(fifo.cachedValues, key)

		// Remove from list
		index := searchIndex(fifo.cachedList, key)
		fifo.cachedList = append(fifo.cachedList[:index], fifo.cachedList[index+1:]...)
		fifo.currentlyUsedCapacity -= len(key) + len(value)
		return value, true
	}

	// If value is not found
	return nil, false
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (fifo *FIFO) Set(key string, value []byte) bool {

	currentObjectSize := len(key) + len(value)
	if currentObjectSize > fifo.capacity {
		return false
	}

	// check if key already
	oldValue, ok := fifo.cachedValues[key]
	if ok {
		// if oldValue and value (new value) are the same length,
		// only replace value in map
		if len(oldValue) == len(value) {
			fifo.cachedValues[key] = value

			// if oldValue was longer than new value, then
			// update currentlyUsedCapacity and map
		} else if len(oldValue) > len(value) {
			fifo.currentlyUsedCapacity -= len(oldValue)
			fifo.cachedValues[key] = value
			fifo.currentlyUsedCapacity += len(value)

			// if oldValue was shorter than new value, then
			// update currentlyUsedCapacity and map
		} else if len(oldValue) < len(value) {
			fifo.currentlyUsedCapacity -= len(oldValue)

			// if not enough capacity after removing old value then evict
			for (fifo.capacity - fifo.currentlyUsedCapacity) < len(value) {
				if successfulEvict := fifo.Evict(); !successfulEvict {
					return false
				}
			}

			fifo.cachedValues[key] = value
			fifo.currentlyUsedCapacity += len(value)
		}
		return true
	}

	// If insufficient capacity, keep evicting elements till
	// enough space
	for fifo.capacity-fifo.currentlyUsedCapacity < currentObjectSize {
		if successfulEvict := fifo.Evict(); !successfulEvict {
			return false
		}
	}

	// Else, sufficient capacity, hence insert
	// Insert in map
	fifo.cachedValues[key] = value

	// Insert in list, appending to the end of the list
	fifo.cachedList = append(fifo.cachedList, key)

	// Increase currentlyUsedCapacity to reflect currentObjectSize
	fifo.currentlyUsedCapacity += currentObjectSize
	return true
}

func (fifo *FIFO) Evict() bool {
	key := (fifo.cachedList)[0]
	// If key is found
	if value, ok := fifo.cachedValues[key]; ok {
		// Delete from map
		delete(fifo.cachedValues, key)

		// Delete from list
		fifo.cachedList = fifo.cachedList[1:]

		// Decrease currentlyUsedCapacity to reflect deletion
		currentObjectSize := len(key) + len(value)
		fifo.currentlyUsedCapacity -= currentObjectSize

		// Eviction successful
		return true
	}
	return false
}

// Len returns the number of bindings in the FIFO.
func (fifo *FIFO) Len() int {
	return len(fifo.cachedValues)
}

// Stats returns statistics about how many search hits and misses have occurred.
func (fifo *FIFO) Stats() *Stats {
	return &fifo.stats
}
