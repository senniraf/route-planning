package shortestpath

import (
	"math"
	"route-planning/graph"
)

type Dijkstra struct {
	Graph graph.Graph
}

func (d Dijkstra) PairShortestPath(s, t graph.Node) (float64, []graph.Edge) {
	predecessor := make([]graph.Edge, d.Graph.N())

	cost := make([]float64, d.Graph.N())
	fill(cost, math.Inf(0))
	cost[s] = 0.0

	pq := newPriorityQueue()
	pq.Push(pqElement{node: s, cost: 0.0})

	for pq.Len() != 0 {
		v := pq.Pop().node

		if v == t {
			break
		}

		for _, e := range d.Graph.OutgoingEdges(v) {
			w := e.To

			if newCost := cost[v] + e.Cost; cost[w] > newCost {
				cost[w] = newCost
				predecessor[w] = e

				if i := pq.Index(w); i < pq.Len() {
					pq.Remove(i)
				}
				pq.Push(pqElement{node: w, cost: cost[w]})
			}
		}
	}

	if cost[t] == math.Inf(0) {
		return math.Inf(0), []graph.Edge{}
	}

	path := []graph.Edge{}
	v := t
	for v != s {
		e := predecessor[v]
		path = prependEdge(path, e)
		v = e.From
	}

	return cost[t], path
}

type BidirectDijkstra struct {
	ForwardGraph  graph.Graph
	BackwardGraph graph.Graph
}

func (bd *BidirectDijkstra) PairShortestPath(s, t graph.Node) (float64, []graph.Edge) {
	forward := newBidirectDijkstraPart(bd.ForwardGraph, s, t)
	backward := newBidirectDijkstraPart(bd.BackwardGraph, t, s)

	upperBound := math.Inf(0)
	var meetingNode graph.Node
	for forward.pq.Len() > 0 || backward.pq.Len() > 0 {

		if upperBound < backward.pq.Top().cost+forward.pq.Top().cost {
			break
		}

		if forward.pq.Top().cost <= backward.pq.Top().cost {
			if cont := forward.advance(*backward, &meetingNode, &upperBound); !cont {
				break
			}
		} else {
			if cont := backward.advance(*forward, &meetingNode, &upperBound); !cont {
				break
			}
		}
	}

	if upperBound == math.Inf(0) {
		return math.Inf(0), []graph.Edge{}
	}

	path := []graph.Edge{}

	if meetingNode == t {
		v := t
		for v != s {
			e := forward.predecessor[v]
			path = prependEdge(path, e)
			v = e.From
		}
		return forward.cost[t], path
	}

	if meetingNode == s {
		v := s
		for v != t {
			e := backward.predecessor[v]
			path = append(path, e.Reverted())
			v = e.From
		}
		return backward.cost[t], path
	}

	v := meetingNode
	for v != s {
		e := forward.predecessor[v]
		path = prependEdge(path, e)
		v = e.From
	}

	v = meetingNode
	for v != t {
		e := backward.predecessor[v]
		path = append(path, e.Reverted())
		v = e.From
	}

	return backward.cost[meetingNode] + forward.cost[meetingNode], path
}

type bidirectDijkstraPart struct {
	g           graph.Graph
	pq          *priorityQueue
	predecessor []graph.Edge
	cost        []float64

	target graph.Node
}

func newBidirectDijkstraPart(g graph.Graph, source, target graph.Node) *bidirectDijkstraPart {
	predecessor := make([]graph.Edge, g.N())

	cost := make([]float64, g.N())
	fill(cost, math.Inf(0))
	cost[source] = 0.0

	pq := newPriorityQueue()
	pq.Push(pqElement{node: source, cost: 0.0})

	return &bidirectDijkstraPart{
		g:           g,
		pq:          pq,
		predecessor: predecessor,
		cost:        cost,
		target:      target,
	}
}

func (bdp *bidirectDijkstraPart) advance(otherPart bidirectDijkstraPart, meetingNode *graph.Node, upperBound *float64) bool {
	v := bdp.pq.Pop().node

	if v == bdp.target {
		*meetingNode = bdp.target
		*upperBound = bdp.cost[v]
		return false
	}

	for _, e := range bdp.g.OutgoingEdges(v) {
		w := e.To

		if newCost := bdp.cost[v] + e.Cost; bdp.cost[w] > newCost {
			bdp.cost[w] = newCost
			bdp.predecessor[w] = e

			if i := bdp.pq.Index(w); i < bdp.pq.Len() {
				bdp.pq.Remove(i)
			}

			bdp.pq.Push(pqElement{node: w, cost: newCost})
			if newUpperBound := newCost + otherPart.cost[w]; *upperBound > newUpperBound {
				*meetingNode = w
				*upperBound = newUpperBound
			}
		}
	}

	return true
}

func prependEdge(edges []graph.Edge, e graph.Edge) []graph.Edge {
	out := append(edges, graph.Edge{})
	copy(out[1:], out)
	out[0] = e
	return out
}

func fill(s []float64, v float64) {
	for i := range s {
		s[i] = v
	}
}
