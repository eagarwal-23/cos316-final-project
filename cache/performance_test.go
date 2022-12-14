package cache

// Will have the performance tests

import (
	//"bytes"
	//"fmt"
	//"testing"
	"crypto/rand"
	"testing"
	"time"
	"math/big"
	"math"
	"strconv"
)

// Evaluate: Latency, Hit ratio
// Use on really large dataset(a lot of misses)

func getRand(tb testing.TB) int64 {
	out, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		tb.Fatal(err)
	}
	return out.Int64()
}

// Small benchmark test to compare lru and arc
func BenchmarkARC_Rand(b *testing.B) {
	capacity := 30
	arc := NewArc(capacity)
	lru := NewLru(capacity)
	

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		trace[i] = getRand(b) % 32768
	}

	b.ResetTimer()

	var hits_arc, misses_arc int
	start_arc := time.Now()
	for i := 0; i < 2*b.N; i++ {
		s := strconv.Itoa(int(trace[i]))
		if i%2 == 0 {
		  arc.Set(s, []byte(s))
		} else {
			if _, ok := arc.Get(s); ok {
				hits_arc++
			} else {
				misses_arc++
			}
		}
	}

	arc_hit_ratio := float64(hits_arc)/float64(misses_arc)
	duration_arc := time.Since(start_arc)
	
	start_lru := time.Now()
	var hits_lru, misses_lru int
	for i := 0; i < 2*b.N; i++ {
		s := strconv.Itoa(int(trace[i]))
		if i%2 == 0 {
		  lru.Set(s, []byte(s))
		} else {
			if _, ok := lru.Get(s); ok {
				hits_lru++
			} else {
				misses_lru++
			}
		}
	}

	lru_hit_ratio := float64(hits_lru)/float64(misses_lru)
	duration_lru := time.Since(start_lru)
	
	//b.Logf("hit: %d miss: %d ratio: %f", )
	b.Logf("LRU hits: %d | LRU misses: %d | LRU ratio: %f, LRU time: %d| ARC hits: %d | ARC misses: %d | ARC ratio: %f, ARC time: %v ", hits_lru, misses_lru, lru_hit_ratio, duration_lru, hits_arc, misses_arc, arc_hit_ratio, duration_arc)
}