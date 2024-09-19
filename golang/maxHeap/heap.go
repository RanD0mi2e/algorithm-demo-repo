package maxHeap

type MaxHeap struct {
	Pq []int
	N  int
}

func NewMaxHeap(size int) *MaxHeap {
	return &MaxHeap{
		Pq: make([]int, size),
		N:  0,
	}
}

func (mh *MaxHeap) Insert(val int) {
	mh.N++
	mh.Pq[mh.N] = val
	mh.Swim(mh.N)
}

func (mh *MaxHeap) DeleteMax() int {
	max := mh.Pq[1]
	mh.Exchange(1, mh.N)
	mh.Pq[mh.N] = 0
	mh.N--
	mh.Sink(1)
	return max
}

func (mh *MaxHeap) IsEmpty() bool {
	return mh.N == 0
}

func (mh *MaxHeap) Size() int {
	return mh.N
}

// 交换元素
func (mh *MaxHeap) Exchange(i, j int) {
	// 位置互换0
	mh.Pq[i], mh.Pq[j] = mh.Pq[j], mh.Pq[i]
}

// 上浮
func (mh *MaxHeap) Swim(i int) {
	for i > 1 && mh.less(i/2, i) {
		mh.Exchange(i/2, i)
		i = i / 2
	}
}

// 下沉
func (mh *MaxHeap) Sink(i int) {
	for 2*i <= mh.N {
		// 左节点
		j := 2 * i
		// 右节点
		if j < mh.N && mh.less(j, j+1) {
			j++
		}
		mh.Exchange(i, j)
		i = j
	}
}

func (mh *MaxHeap) less(i, j int) bool {
	if mh.Pq[i] < mh.Pq[j] {
		return true
	}
	return false
}

// 索引 i 节点的子节点为 2i + 1, 2i + 2
// 索引 i 节点的父节点为 (i - 1 / 2) 向下取整

// 往数组最后位插入新元素 i：
// heapify_up 递归向上查找是否存在 h[i] > h[parent(i)], 存在则交换两个元素，直到终止条件 i = 0

// 直接取根元素就是最大值

// 根元素和最后一个叶子节点互换，互换后删掉最后一个叶子节点
// heapify_down 递归向下查找是否存在 h[i] < h[chile[i]], 存在则交换两个元素，直到终止条件 2i + 1(或者2i + 2) > len(h) - 1
// 此时当前节点已经是叶子节点了

// 堆排序：
// 1.建立一个定长数组
// 2.往定长数组填入数据时建立起最大堆
// 3.不断取出最大值并和定长数组最后一位互换
// 4.互换完毕后，可以得到从小到大排序的数组
