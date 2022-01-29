package shortestpath

import "route-planning/graph"

type Algorithm interface {
	PairShortestPath(s, t graph.Node) (float64, []graph.Edge)
}
