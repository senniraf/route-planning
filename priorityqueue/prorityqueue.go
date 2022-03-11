package priorityqueue

import "route-planning/graph"

type Element struct {
	Node graph.Node
	Cost float64
}

type PriorityQueue interface {
	Len() int
	Push(e Element)
	Pop() Element
	Top() Element
}
