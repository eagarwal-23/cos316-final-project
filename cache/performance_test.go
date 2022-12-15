package cache

// Will have the performance tests

import (
	//"bytes"
	//"fmt"
	//"testing"
	"crypto/rand"
	"testing"

	//"time"
	"math"
	"math/big"
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

// Small benchmark test to compare arc
func BenchmarkArc(b *testing.B) {
	capacities := []int{20480} //160, 320, 640, 1280, 2560, 5120, 10240, 20480}

	for i := range capacities {
		capacity := capacities[i]
		arc := NewArc(capacity)

		trace := make([]int64, b.N*2)
		for i := 0; i < b.N*2; i++ {
			trace[i] = getRand(b) % 32768
		}

		b.ResetTimer()

		var hits, misses, total int
		for i := 0; i < 2*b.N; i++ {
			s := strconv.Itoa(int(trace[i]))
			if i%2 == 0 {
				arc.Set(s, []byte(s))
			} else {
				if _, ok := arc.Get(s); ok {
					hits++
				} else {
					misses++
				}
				total = hits + misses
			}
		}
		hit_ratio := float64(hits) / float64(misses)
		b.Logf("ARC(%d): hits: %d | misses: %d | hit-ratio: %f", total, hits, misses, hit_ratio)
	}

}

func BenchmarkLru(b *testing.B) {
	capacities := []int{20480} //160, 320, 640, 1280, 2560, 5120, 10240, 20480}

	for i := range capacities {
		capacity := capacities[i]
		lru := NewLru(capacity)

		trace := make([]int64, b.N*2)
		for i := 0; i < b.N*2; i++ {
			trace[i] = getRand(b) % 327680
		}

		b.ResetTimer()

		var hits, misses, total int
		for i := 0; i < 2*b.N; i++ {
			s := strconv.Itoa(int(trace[i]))
			if i%2 == 0 {
				lru.Set(s, []byte(s))
			} else {
				if _, ok := lru.Get(s); ok {
					hits++
				} else {
					misses++
				}
			}
			total = hits + misses
		}
		hit_ratio := float64(hits) / float64(misses)
		b.Logf("LRU(%d): hits: %d | misses: %d | hit-ratio: %f", total, hits, misses, hit_ratio)
	}
}

// Using trace from webcachism
