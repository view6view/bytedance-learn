package sort

// InsertSort 插入排序
func InsertSort(arr []int) {
	for i, num := range arr {
		insertIdx := i - 1
		for insertIdx > -1 && arr[insertIdx] > num {
			arr[insertIdx+1] = arr[insertIdx]
			insertIdx--
		}
		arr[insertIdx+1] = num
	}
}

// QuickSort 快速排序,以最后一个元素为pivot
func QuickSort(arr []int) {
	var partition = func(begin, end int) int {
		pointer := begin
		pivot := arr[end]
		var temp int
		for i := begin; i < end; i++ {
			if arr[i] <= pivot {
				if i != pointer {
					temp = arr[i]
					arr[i] = arr[pointer]
					arr[pointer] = temp
				}
				pointer++
			}
		}
		temp = arr[pointer]
		arr[pointer] = arr[end]
		arr[end] = temp
		return pointer
	}

	var quickSort func(begin, end int)
	quickSort = func(begin, end int) {
		if begin < end {
			position := partition(begin, end)
			quickSort(begin, position-1)
			quickSort(position+1, end)
		}
	}

	quickSort(0, len(arr)-1)
}

// HeapSort 堆排序
func HeapSort(arr []int) {

	var siftDown = func(i, end int) {
		// 取出当前元素
		temp := arr[i]
		// 从i节点的左子节点开始，也就是2*i+1处开始
		for j := i*2 + 1; j < end; j = j*2 + 1 {
			// 如果左子节点小于右子节点，j指向右子节点，因为需要构建大顶堆
			if j+1 < end && arr[j] < arr[j+1] {
				j++
			}
			// 如果子节点大于父节点，将子节点的值赋值父节点，不用交换
			if arr[j] > temp {
				arr[i] = arr[j]
				i = j
			} else {
				break
			}
		}
		arr[i] = temp
	}

	len := len(arr)
	// 构建大顶堆
	for i := len/2 - 1; i >= 0; i-- {
		// 从第一个非叶子节点从下到上，从右至左调整结构
		siftDown(i, len)
	}
	// 调整堆结构，交换堆顶元素与末尾元素
	var temp int
	for i := len - 1; i > 0; i-- {
		// 将堆顶元素与末尾元素进行交换
		temp = arr[0]
		arr[0] = arr[i]
		arr[i] = temp
		// 重新对堆进行调整
		siftDown(0, i)
	}
}
