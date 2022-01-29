package graph

type Graph interface {
	OutgoingEdges(Node) []Edge
	N() int
	Reverted() Graph
}

type Node int

type Edge struct {
	From Node
	To   Node
	Cost float64
}

func (e Edge) Reverted() Edge {
	return Edge{
		From: e.To,
		To:   e.From,
		Cost: e.Cost,
	}
}
