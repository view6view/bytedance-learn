package benchmark

import (
	// 字节的fast随机库
	"github.com/bytedance/gopkg/lang/fastrand"
	// 为了保证随机性，其实有一把全局锁，性能会低一点
	"math/rand"
)

var ServerIndex [10]int

func InitServerIndex() {
	for i := 0; i < 10; i++ {
		ServerIndex[i] = i + 100
	}
}

func Select() int {
	return ServerIndex[rand.Intn(10)]
}

func FastSelect() int {
	return ServerIndex[fastrand.Intn(10)]
}
