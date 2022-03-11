package shortestpath

import "route-planning/graph"

type Algorithm interface {
	Pair(s, t graph.Node) (float64, []graph.Edge)
}
