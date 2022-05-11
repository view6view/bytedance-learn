package benchmap

import "testing"

func BenchmarkNoPreAlloc(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NoPreAlloc(1000)
	}
}

func BenchmarkPreAlloc(b *testing.B) {
	for n := 0; n < b.N; n++ {
		PreAlloc(1000)
	}
}
