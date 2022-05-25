package sort

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSort(t *testing.T) {
	var isSort = func(arr []int) bool {
		for i := 0; i < len(arr)-1; i++ {
			if arr[i] > arr[i+1] {
				return false
			}
		}
		return true
	}

	arr1 := []int{3, 7, 9, 2, 1, 5, 6, 34, 64, 765, 23, 547, 12, 534, 57, 12, 12, 325}
	len := len(arr1)
	arr2 := make([]int, len)
	copy(arr2, arr1)
	arr3 := make([]int, len)
	copy(arr3, arr1)

	assert.Equal(t, false, isSort(arr1))
	InsertSort(arr1)
	assert.Equal(t, true, isSort(arr1))

	assert.Equal(t, false, isSort(arr2))
	QuickSort(arr2)
	assert.Equal(t, true, isSort(arr2))

	assert.Equal(t, false, isSort(arr3))
	HeapSort(arr3)
	assert.Equal(t, true, isSort(arr3))
}
