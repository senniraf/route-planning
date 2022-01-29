package shortestpath

import (
	"container/heap"
	"route-planning/graph"
)

type priorityQueue struct {
	elements *pqElements
	top      pqElement
	topped   bool
}

type pqElements []pqElement

type pqElement struct {
	node graph.Node
	cost float64
}

func newPriorityQueue() *priorityQueue {
	return &priorityQueue{elements: &pqElements{}}
}

func (pq priorityQueue) Index(n graph.Node) int {
	for i, element := range *pq.elements {
		if element.node == n {
			return i
		}
	}

	return pq.Len()
}

func (pq priorityQueue) Len() int {
	return pq.elements.Len()
}

func (pq *priorityQueue) Push(e pqElement) {
	if pq.topped {
		pq.topped = false
		heap.Push(pq.elements, pq.top)
	}
	heap.Push(pq.elements, e)
}

func (pq *priorityQueue) Pop() pqElement {
	if pq.topped {
		pq.topped = false
		return pq.top
	}
	return heap.Pop(pq.elements).(pqElement)
}

func (pq *priorityQueue) Top() pqElement {
	if pq.topped {
		return pq.top
	}
	pq.top = pq.Pop()
	pq.topped = true
	return pq.top
}

func (pq *priorityQueue) Remove(i int) {
	heap.Remove(pq.elements, i)
}

func (es pqElements) Less(i, j int) bool {
	return es[i].cost < es[j].cost
}

func (es pqElements) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

func (es pqElements) Len() int {
	return len(es)
}

func (es *pqElements) Push(x interface{}) {
	*es = append(*es, x.(pqElement))
}

func (es *pqElements) Pop() interface{} {
	old := *es
	n := len(old)
	x := old[n-1]
	*es = old[0 : n-1]
	return x
}
