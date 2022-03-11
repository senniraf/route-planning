package shortestpath

import (
	"math"
	"route-planning/graph"
	"route-planning/priorityqueue"
)

type Dijkstra struct {
	Graph graph.Graph
}

type StoppingCriterion func(priorityqueue.Element, DijkstraState) bool

type DijkstraState struct {
	Cost        []float64
	Predecessor []graph.Edge
	PQ          priorityqueue.PriorityQueue
}

func (d Dijkstra) Pair(s, t graph.Node) (float64, []graph.Edge) {

	stop := func(element priorityqueue.Element, state DijkstraState) bool {
		return element.Node == t
	}

	cost, predecessor := d.Run(s, stop)

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

func (d Dijkstra) ToAll(s graph.Node) ([]float64, []graph.Edge) {
	return d.Run(s, nil)
}

func (d Dijkstra) Run(s graph.Node, stop StoppingCriterion) ([]float64, []graph.Edge) {
	if stop == nil {
		stop = func(el priorityqueue.Element, ds DijkstraState) bool {
			return false
		}
	}

	state := DijkstraState{
		Predecessor: make([]graph.Edge, d.Graph.N()),
		Cost:        make([]float64, d.Graph.N()),
		PQ:          priorityqueue.NewMinHeap(),
	}

	fill(state.Cost, math.Inf(0))
	state.Cost[s] = 0.0

	state.PQ.Push(priorityqueue.Element{Node: s, Cost: 0.0})

	for state.PQ.Len() != 0 {
		element := state.PQ.Pop()

		if stop(element, state) {
			break
		}

		v := element.Node

		if element.Cost > state.Cost[v] {
			continue
		}

		for _, e := range d.Graph.OutgoingEdges(v) {
			w := e.To

			if newCost := state.Cost[v] + e.Cost; state.Cost[w] > newCost {
				state.Cost[w] = newCost
				state.Predecessor[w] = e

				state.PQ.Push(priorityqueue.Element{Node: w, Cost: state.Cost[w]})
			}
		}
	}

	return state.Cost, state.Predecessor
}

type BidirectDijkstra struct {
	ForwardGraph  graph.Graph
	BackwardGraph graph.Graph
}

func (bd *BidirectDijkstra) Pair(s, t graph.Node) (float64, []graph.Edge) {
	forward := newBidirectDijkstraPart(bd.ForwardGraph, s, t)
	backward := newBidirectDijkstraPart(bd.BackwardGraph, t, s)

	upperBound := math.Inf(0)
	var meetingNode graph.Node
	for forward.pq.Len() > 0 || backward.pq.Len() > 0 {

		if upperBound < backward.pq.Top().Cost+forward.pq.Top().Cost {
			break
		}

		if forward.pq.Top().Cost <= backward.pq.Top().Cost {
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
	pq          priorityqueue.PriorityQueue
	predecessor []graph.Edge
	cost        []float64

	target graph.Node
}

func newBidirectDijkstraPart(g graph.Graph, source, target graph.Node) *bidirectDijkstraPart {
	predecessor := make([]graph.Edge, g.N())

	cost := make([]float64, g.N())
	fill(cost, math.Inf(0))
	cost[source] = 0.0

	pq := priorityqueue.NewMinHeap()
	pq.Push(priorityqueue.Element{Node: source, Cost: 0.0})

	return &bidirectDijkstraPart{
		g:           g,
		pq:          pq,
		predecessor: predecessor,
		cost:        cost,
		target:      target,
	}
}

func (bdp *bidirectDijkstraPart) advance(otherPart bidirectDijkstraPart, meetingNode *graph.Node, upperBound *float64) bool {
	element := bdp.pq.Pop()
	v := element.Node

	if v == bdp.target {
		*meetingNode = bdp.target
		*upperBound = bdp.cost[v]
		return false
	}

	if element.Cost > bdp.cost[v] {
		return true
	}

	for _, e := range bdp.g.OutgoingEdges(v) {
		w := e.To

		if newCost := bdp.cost[v] + e.Cost; bdp.cost[w] > newCost {
			bdp.cost[w] = newCost
			bdp.predecessor[w] = e

			bdp.pq.Push(priorityqueue.Element{Node: w, Cost: newCost})
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
