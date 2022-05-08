package concurrence

func CalSquare() {
	src := make(chan int)
	dest := make(chan int, 3)
	// 生产者
	go func() {
		defer close(src)
		for i := 0; i < 10; i++ {
			src <- i
		}
	}()
	// 消费者 -> 生产者
	go func() {
		defer close(dest)
		for i := range src {
			dest <- i * i
		}
	}()
	for i := range dest {
		//复杂操作
		println(i)
	}
}
