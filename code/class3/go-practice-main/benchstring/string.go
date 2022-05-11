package benchstring

import (
	"bytes"
	"strings"
)

func Plus(n int, str string) string {
	s := ""
	for i := 0; i < n; i++ {
		s += str
	}
	return s
}

func StrBuilder(n int, str string) string {
	var builder strings.Builder
	for i := 0; i < n; i++ {
		builder.WriteString(str)
	}
	return builder.String()
}

func ByteBuffer(n int, str string) string {
	buf := new(bytes.Buffer)
	for i := 0; i < n; i++ {
		buf.WriteString(str)
	}
	return buf.String()
}

func PreStrBuilder(n int, str string) string {
	var builder strings.Builder
	builder.Grow(n * len(str))
	for i := 0; i < n; i++ {
		builder.WriteString(str)
	}
	return builder.String()
}

func PreByteBuffer(n int, str string) string {
	buf := new(bytes.Buffer)
	buf.Grow(n * len(str))
	for i := 0; i < n; i++ {
		buf.WriteString(str)
	}
	return buf.String()
}
