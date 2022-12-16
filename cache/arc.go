package cache

// An ARC is a fixed-size in-memory cache with least-recently-used eviction
type ARC struct {
	p        int // P is the dynamic preference towards t1 or t2
	capacity int // To hold the capacity of the cache

	t1 *LRU // To hold recent cache entries
	t2 *LRU // To hold frequent cache entries, referenced at least twice
	b1 *LRU // To hold ghost entries evicted from the t1 cache
	b2 *LRU // To hold ghost entries evicted from the t2 cache

	currentlyUsedCapacity int   // Currently used capacity of the cache
	stats                 Stats // Hits and misses for the cache
}

// NewArc returns a pointer to a new ARC with a capacity to store limit bytes
func NewArc(limit int) *ARC {
	return &ARC{p: 0, capacity: limit, t1: NewLru(limit), t2: NewLru(limit), b1: NewLru(limit), b2: NewLru(limit), currentlyUsedCapacity: 0, stats: Stats{}}
}

// MaxStorage returns the maximum number of bytes this  can store
func (arc *ARC) MaxStorage() int {
	return arc.capacity
}

// RemainingStorage returns the number of unused bytes available in this ARC
func (arc *ARC) RemainingStorage() int {
	return arc.capacity - arc.currentlyUsedCapacity
}

// Peek returns the value associated with the given key, if it exists.
// This operation does not counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (arc *ARC) Peek(key string) (value []byte, ok bool) {
	currMapping, ok := arc.t1.Peek(key)
	if ok {
		return currMapping, ok
	} else {
		return arc.t2.Peek(key)
	}
}

// Get returns the value associated with the given key, if it exists.
// This operation counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (arc *ARC) Get(key string) (value []byte, ok bool) {
	if val, ok := arc.t1.Get(key); ok {
		arc.t1.Remove(key)
		arc.t2.Set(key, val)
		arc.stats.Hits += 1
		return val, ok
	}

	val, ok := arc.t2.Get(key)
	if ok {
		arc.stats.Hits += 1
	} else {
		arc.stats.Misses += 1
	}

	if _, ok := arc.b1.Peek(key); ok {
		arc.stats.B1Hits += 1
	}

	if _, ok := arc.b2.Peek(key); ok {
		arc.stats.B2Hits += 1
	}

	return val, ok
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (arc *ARC) Remove(key string) (value []byte, ok bool) {
	if val, ok := arc.t1.Remove(key); ok {
		arc.currentlyUsedCapacity -= len(key) + len(val)
		return val, ok
	}

	if val, ok := arc.t2.Remove(key); ok {
		arc.currentlyUsedCapacity -= len(key) + len(val)
		return val, ok
	}

	if val, ok := arc.b1.Remove(key); ok {
		return val, ok
	}

	if val, ok := arc.b2.Remove(key); ok {
		return val, ok
	}

	return nil, false

}

func (arc *ARC) updateCapacity() {
	arc.currentlyUsedCapacity = arc.t1.currentlyUsedCapacity + arc.t2.currentlyUsedCapacity
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (arc *ARC) Set(key string, value []byte) bool {
	currObjectSize := len(key) + len(value)
	if currObjectSize > arc.capacity {
		return false
	}

	// If the key is in recently-used cache t1, then promote it to t2
	if _, ok := arc.t1.Peek(key); ok {
		arc.t1.Remove(key)
		arc.t2.Set(key, value)
		arc.updateCapacity()
		return true
	}

	// If key is already in frequently-used cache t2, then update its corresponding value
	if _, ok := arc.t2.Peek(key); ok {
		arc.t2.Set(key, value)
		arc.updateCapacity()
		return true
	}

	// If key is part of ghost entries recently-evicted from recently-used list,
	// adjust dynamic preference towards t1 v t2 in favour of t1, because client's
	// usage shows preference for recently-used entries
	if _, ok := arc.b1.Peek(key); ok {
		change := 1
		if arc.b2.currentlyUsedCapacity > arc.b1.currentlyUsedCapacity {
			change = (arc.b2.currentlyUsedCapacity) / (arc.b1.currentlyUsedCapacity)
		}
		updated_p := arc.p + change
		if updated_p > (arc.capacity) {
			arc.p = arc.capacity
		} else {
			arc.p = updated_p
		}
		for arc.capacity < arc.currentlyUsedCapacity+currObjectSize {
			arc.Replace(key)
		}

		// move key-value pair from b1 into t2
		arc.b1.Remove(key)
		arc.t2.Set(key, value)
		arc.updateCapacity()
	}

	// If key is part of ghost entries recently-evicted from frequently-used list,
	// adjust dynamic preference towards t1 v t2 in favour of t2, because client's
	// usage shows preference for frequently-used entries
	if _, ok := arc.b2.Peek(key); ok {
		// fmt.Println("I AM SETTING: D")
		change := 1
		if arc.b2.currentlyUsedCapacity > arc.b1.currentlyUsedCapacity {
			change = (arc.b2.currentlyUsedCapacity) / (arc.b1.currentlyUsedCapacity)
		}
		updated_p := arc.p - change
		if updated_p < 0 {
			arc.p = 0
		} else {
			arc.p = updated_p
		}
		for arc.capacity < arc.currentlyUsedCapacity+currObjectSize {
			arc.Replace(key)
		}

		// move key-value pair from b2 into t2
		arc.b2.Remove(key)
		arc.t2.Set(key, value)
		arc.updateCapacity()
	}

	// if not in any of the four
	// if b1.Len() + t1.Len() = arc.capacity, then if b1 is not empty,
	// delete from b1 & move key-value pair from t1 to b1

	if arc.t1.currentlyUsedCapacity+arc.b1.currentlyUsedCapacity == arc.capacity {
		if arc.t1.currentlyUsedCapacity < arc.capacity {
			arc.b1.Evict()
			for arc.capacity < arc.currentlyUsedCapacity+currObjectSize {
				arc.Replace(key)
			}
		} else {
			arc.t1.Evict()
		}
	}
	if arc.t1.currentlyUsedCapacity+arc.b1.currentlyUsedCapacity < arc.capacity {
		total_cap := arc.t1.currentlyUsedCapacity + arc.t2.currentlyUsedCapacity + arc.b1.currentlyUsedCapacity + arc.b2.currentlyUsedCapacity
		if total_cap >= arc.capacity {
			if total_cap == 2*arc.capacity {
				arc.b2.Evict()
				for arc.capacity < arc.currentlyUsedCapacity+currObjectSize {
					arc.Replace(key)
				}
			}
		}
	}

	arc.t1.Set(key, value)
	arc.updateCapacity()
	return true
}

func (arc *ARC) Replace(key string) {
	_, key_in_b2 := arc.b1.Peek(key)
	t1_length := arc.t1.currentlyUsedCapacity
	if (t1_length > 0) && (t1_length > arc.p || (t1_length == arc.p && key_in_b2)) {
		key, ok := arc.t1.Evict()
		emptyVal := make([]byte, 0)
		if ok {
			arc.b1.Set(key, emptyVal)
		}
	} else {
		key, ok := arc.t2.Evict()
		emptyVal := make([]byte, 0)
		if ok {
			arc.b2.Set(key, emptyVal)
		}
	}
	arc.updateCapacity()
}

func (arc *ARC) Empty() {
	arc.t1.Empty()
	arc.t2.Empty()
	arc.b1.Empty()
	arc.b2.Empty()
	arc.updateCapacity()
}

// Len returns the number of bindings in the ARC.
func (arc *ARC) Len() int {
	return arc.t1.Len() + arc.t2.Len()
}

// Stats returns statistics about how many search hits and misses have occurred.
func (arc *ARC) Stats() *Stats {
	return &arc.stats
}

/*
SOURCES

https://www.youtube.com/watch?v=S6IfqDXWa10
*/
