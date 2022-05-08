package util

import (
	"math"
)

// GetUniId 获取唯一id
func GetUniId(ids []int) int {
	//// 获取最后一个数下标
	//idx := len(ids) - 1
	//// 对id进行排序
	//sort.Ints(ids)
	//return ids[idx] + 1
	max := math.MinInt
	for _, id := range ids {
		max = Max(max, id)
	}
	return max + 1
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
