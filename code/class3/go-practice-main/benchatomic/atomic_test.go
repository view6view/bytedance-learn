package benchatomic

import (
	"testing"
)

func BenchmarkAtomicAddOne(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var counter = atomicCounter{}
		AtomicAddOne(&counter)
	}
}

func BenchmarkMutexAddOne(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var counter = mutexCounter{}
		MutexAddOne(&counter)
	}
}
