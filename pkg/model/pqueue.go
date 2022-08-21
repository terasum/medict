package model

type PriorityQueue struct {
	targetW *WordItem
	list    []*WordItem
}

func NewPQ(word *WordItem) *PriorityQueue {
	return &PriorityQueue{
		targetW: word,
		list:    make([]*WordItem, 0),
	}
}

func (pq PriorityQueue) Len() int { return len(pq.list) }
func (pq PriorityQueue) Less(i, j int) bool {
	a := pq.list[i]
	b := pq.list[j]
	return a.Distance(pq.targetW) < b.Distance(pq.targetW)
}

// We just implement the pre-defined function in interface of heap.

func (pq *PriorityQueue) Pop() interface{} {
	n := len(pq.list)
	item := pq.list[n-1]
	pq.list = pq.list[0 : n-1]
	return item
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*WordItem)
	pq.list = append(pq.list, item)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq.list[i], pq.list[j] = pq.list[j], pq.list[i]
}
