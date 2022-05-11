package benchmap

func NoPreAlloc(size int) {
	data := make(map[int]int)
	for i := 0; i < size; i++ {
		data[i] = 1
	}
}

func PreAlloc(size int) {
	data := make(map[int]int, size)
	for i := 0; i < size; i++ {
		data[i] = 1
	}
}
