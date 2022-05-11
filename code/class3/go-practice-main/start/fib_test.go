package start

import "testing"

// BenchmarkFib10 run 'go test -bench=. -benchmem' to get the benchmark result
func BenchmarkFib10(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		Fib(10)
	}
}
