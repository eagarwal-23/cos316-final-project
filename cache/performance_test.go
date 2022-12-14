package cache

// Will have the performance tests

import (
	//"bytes"
	//"fmt"
	//"testing"
	"crypto/rand"
	"testing"
	//"time"
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

func BenchmarkARC_Rand(b *testing.B) {
	capacity := 30
	l := NewArc(capacity)
	

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		trace[i] = getRand(b) % 32768
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < 2*b.N; i++ {
		s := strconv.Itoa(int(trace[i]))
		if i%2 == 0 {
			l.Set(s, []byte(s))
		} else {
			if _, ok := l.Get(s); ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}