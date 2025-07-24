package heapsort

func Heapify(arr []int, n int, i int) {
	for {
		largest := i
		left := 2*i + 1
		right := 2*i + 2

		// 找到最大值的索引
		if left < n && arr[left] > arr[largest] {
			largest = left
		}

		if right < n && arr[right] > arr[largest] {
			largest = right
		}

		// 如果最大值就是当前节点，堆化完成
		if largest == i {
			break
		}

		// 交换并继续向下堆化
		arr[i], arr[largest] = arr[largest], arr[i]
		i = largest
	}
}

func BuildHeap(arr []int) {
	n := len(arr)
	for i := n/2 - 1; i >= 0; i-- {
		Heapify(arr, n, i)
	}
}

func Heapsort(arr []int) {
	n := len(arr)
	BuildHeap(arr)

	for i := n - 1; i > 0; i-- {
		arr[0], arr[i] = arr[i], arr[0]
		Heapify(arr, i, 0)
	}
}
