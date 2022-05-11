package benchslice

func NoPreAlloc(size int) {
	data := make([]int, 0)
	for k := 0; k < size; k++ {
		data = append(data, k)
	}
}

func PreAlloc(size int) {
	data := make([]int, 0, size)
	for k := 0; k < size; k++ {
		data = append(data, k)
	}
}

func GetLastBySlice(origin []int) []int {
	return origin[len(origin)-2:]
}

func GetLastByCopy(origin []int) []int {
	result := make([]int, 2)
	copy(result, origin[len(origin)-2:])
	return result
}
