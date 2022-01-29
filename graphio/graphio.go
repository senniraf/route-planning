package graphio

import "route-planning/graph"

type GraphInput interface {
	LoadGraph() ([]graph.Edge, int, error)
}

type GraphOutput interface {
	PrintNode(graph.Node) string
	PrintEdge(graph.Edge) string
}
