package priorityqueue

import (
	"container/heap"
)

type minHeap struct {
	elements *elements
}

func NewMinHeap() PriorityQueue {
	return &minHeap{elements: &elements{}}
}

func (pq minHeap) Len() int {
	return pq.elements.Len()
}

func (pq *minHeap) Push(e Element) {
	heap.Push(pq.elements, e)
}

func (pq *minHeap) Pop() Element {
	return heap.Pop(pq.elements).(Element)
}

func (pq *minHeap) Top() Element {
	return (*pq.elements)[0]
}

type elements []Element

func (es elements) Less(i, j int) bool {
	return es[i].Cost < es[j].Cost
}

func (es elements) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

func (es elements) Len() int {
	return len(es)
}

func (es *elements) Push(x interface{}) {
	*es = append(*es, x.(Element))
}

func (es *elements) Pop() interface{} {
	old := *es
	n := len(old)
	x := old[n-1]
	*es = old[0 : n-1]
	return x
}
