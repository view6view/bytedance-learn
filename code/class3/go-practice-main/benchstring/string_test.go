package benchstring

import "testing"

func BenchmarkPlus(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Plus(1000, "string")
	}
}

func BenchmarkStrBuilder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		StrBuilder(1000, "string")
	}
}

func BenchmarkByteBuffer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ByteBuffer(1000, "string")
	}
}

func BenchmarkPreStrBuilder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		PreStrBuilder(1000, "string")
	}
}

func BenchmarkPreByteBuffer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		PreByteBuffer(1000, "string")
	}
}
