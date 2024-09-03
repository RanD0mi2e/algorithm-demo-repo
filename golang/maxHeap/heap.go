package maxHeap

type intHeap []int

// 索引 i 节点的子节点为 2i + 1, 2i + 2
// 索引 i 节点的父节点为 (i - 1 / 2) 向下取整

// 往数组最后位插入新元素 i：
// heapify_up 递归向上查找是否存在 h[i] > h[parent(i)], 存在则交换两个元素，直到终止条件 i = 0
func (h *intHeap) Insert() {}

// 直接取根元素就是最大值
func (h *intHeap) FindMax() {}

// 根元素和最后一个叶子节点互换，互换后删掉最后一个叶子节点
// heapify_down 递归向下查找是否存在 h[i] < h[chile[i]], 存在则交换两个元素，直到终止条件 2i + 1(或者2i + 2) > len(h) - 1
// 此时当前节点已经是叶子节点了
func (h *intHeap) DeleteMax() {}
func (h *intHeap) Build()     {}

// 堆排序：
// 1.建立一个定长数组
// 2.往定长数组填入数据时建立起最大堆
// 3.不断取出最大值并和定长数组最后一位互换
// 4.互换完毕后，可以得到从小到大排序的数组
