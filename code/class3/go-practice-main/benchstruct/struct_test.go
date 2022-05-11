package benchstruct

import "testing"

func BenchmarkEmptyStructMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		EmptyStructMap(10000)
	}
}

func BenchmarkBoolMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		BoolMap(10000)
	}
}
