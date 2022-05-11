package benchslice

import (
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func BenchmarkNoPreAlloc(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NoPreAlloc(100)
	}
}

func BenchmarkPreAlloc(b *testing.B) {
	for n := 0; n < b.N; n++ {
		PreAlloc(100)
	}
}

func generateWithCap(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func printMem(t *testing.T) {
	t.Helper()
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	t.Logf("%.2f MB", float64(rtm.Alloc)/1024./1024.)
}

func testGetLast(t *testing.T, f func([]int) []int) {
	result := make([][]int, 0)
	for k := 0; k < 100; k++ {
		origin := generateWithCap(128 * 1024) // 1M
		result = append(result, f(origin))
	}
	printMem(t)
	_ = result
}

func TestLastBySlice(t *testing.T) {
	testGetLast(t, GetLastBySlice)
}

func TestLastByCopy(t *testing.T) {
	testGetLast(t, GetLastByCopy)
}
