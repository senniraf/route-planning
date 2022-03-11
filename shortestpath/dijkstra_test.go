package shortestpath_test

import (
	"math/rand"
	"reflect"
	"route-planning/graph"
	"route-planning/shortestpath"
	"testing"
)

var testEdges = []graph.Edge{
	{From: 0, To: 1, Cost: 4},
	{From: 0, To: 7, Cost: 8},
	{From: 1, To: 0, Cost: 4},
	{From: 1, To: 2, Cost: 8},
	{From: 1, To: 7, Cost: 11},
	{From: 2, To: 1, Cost: 8},
	{From: 2, To: 3, Cost: 7},
	{From: 2, To: 5, Cost: 4},
	{From: 2, To: 8, Cost: 2},
	{From: 3, To: 2, Cost: 7},
	{From: 3, To: 4, Cost: 9},
	{From: 3, To: 5, Cost: 14},
	{From: 4, To: 3, Cost: 9},
	{From: 4, To: 5, Cost: 10},
	{From: 5, To: 2, Cost: 4},
	{From: 5, To: 3, Cost: 14},
	{From: 5, To: 4, Cost: 10},
	{From: 5, To: 6, Cost: 2},
	{From: 6, To: 5, Cost: 2},
	{From: 6, To: 7, Cost: 1},
	{From: 6, To: 8, Cost: 6},
	{From: 7, To: 0, Cost: 8},
	{From: 7, To: 1, Cost: 11},
	{From: 7, To: 6, Cost: 1},
	{From: 7, To: 8, Cost: 7},
	{From: 8, To: 2, Cost: 2},
	{From: 8, To: 6, Cost: 6},
	{From: 8, To: 7, Cost: 7},
}

var testGraph = graph.NewAdjacencyList(testEdges, 9)

var expectedCosts = [][]float64{
	{0, 4, 12, 19, 21, 11, 9, 8, 14},
	{4, 0, 8, 15, 22, 12, 12, 11, 10},
	{12, 8, 0, 7, 14, 4, 6, 7, 2},
	{19, 15, 7, 0, 9, 11, 13, 14, 9},
	{21, 22, 14, 9, 0, 10, 12, 13, 16},
	{11, 12, 4, 11, 10, 0, 2, 3, 6},
	{9, 12, 6, 13, 12, 2, 0, 1, 6},
	{8, 11, 7, 14, 13, 3, 1, 0, 7},
	{14, 10, 2, 9, 16, 6, 6, 7, 0},
}

//nolint:govet // These are edges in the form (From, To, Cost)
var expectedPredecessors = [][]graph.Edge{
	{{0, 0, 0}, {0, 1, 4}, {1, 2, 8}, {2, 3, 7}, {5, 4, 10}, {6, 5, 2}, {7, 6, 1}, {0, 7, 8}, {2, 8, 2}},
	{{1, 0, 4}, {0, 0, 0}, {1, 2, 8}, {2, 3, 7}, {5, 4, 10}, {2, 5, 4}, {7, 6, 1}, {1, 7, 11}, {2, 8, 2}},
	{{1, 0, 4}, {2, 1, 8}, {0, 0, 0}, {2, 3, 7}, {5, 4, 10}, {2, 5, 4}, {5, 6, 2}, {6, 7, 1}, {2, 8, 2}},
	{{1, 0, 4}, {2, 1, 8}, {3, 2, 7}, {0, 0, 0}, {3, 4, 9}, {2, 5, 4}, {5, 6, 2}, {6, 7, 1}, {2, 8, 2}},
	{{7, 0, 8}, {2, 1, 8}, {5, 2, 4}, {4, 3, 9}, {0, 0, 0}, {4, 5, 10}, {5, 6, 2}, {6, 7, 1}, {2, 8, 2}},
	{{7, 0, 8}, {2, 1, 8}, {5, 2, 4}, {2, 3, 7}, {5, 4, 10}, {0, 0, 0}, {5, 6, 2}, {6, 7, 1}, {2, 8, 2}},
	{{7, 0, 8}, {7, 1, 11}, {5, 2, 4}, {2, 3, 7}, {5, 4, 10}, {6, 5, 2}, {0, 0, 0}, {6, 7, 1}, {6, 8, 6}},
	{{7, 0, 8}, {7, 1, 11}, {5, 2, 4}, {2, 3, 7}, {5, 4, 10}, {6, 5, 2}, {7, 6, 1}, {0, 0, 0}, {7, 8, 7}},
	{{1, 0, 4}, {2, 1, 8}, {8, 2, 2}, {2, 3, 7}, {5, 4, 10}, {2, 5, 4}, {8, 6, 6}, {8, 7, 7}, {0, 0, 0}},
}

func TestPair(t *testing.T) {
	sut := shortestpath.Dijkstra{Graph: testGraph}

	for v := 0; v < testGraph.N(); v++ {
		for w := 0; w < testGraph.N(); w++ {
			if v == w {
				continue
			}

			cost, path := sut.Pair(graph.Node(v), graph.Node(w))
			if cost != expectedCosts[v][w] {
				t.Errorf("sp(%d, %d): expected %f, got %f", v, w, expectedCosts[v][w], cost)
			}

			for i, e := range path {
				if e != expectedPredecessors[v][e.To] {
					t.Errorf("sp(%d, %d) at pos %d: expected %v, got %v", v, w, i, expectedPredecessors[v][e.To], e)
				}
			}
		}
	}
}

func TestToAll(t *testing.T) {
	sut := shortestpath.Dijkstra{Graph: testGraph}

	for v := 0; v < testGraph.N(); v++ {
		costs, predecessor := sut.ToAll(graph.Node(v))
		if !reflect.DeepEqual(costs, expectedCosts[v]) {
			t.Errorf("sp(%d) costs: expeted: %v, got: %v", v, expectedCosts[v], costs)
		}

		if !reflect.DeepEqual(predecessor, expectedPredecessors[v]) {
			t.Errorf("sp(%d) predecessor: expeted: %v, got: %v", v, expectedPredecessors[v], predecessor)
		}
	}
}

func BenchmarkPair(b *testing.B) {
	b.StopTimer()
	n := 100_000
	m := 4 * n

	sut := shortestpath.Dijkstra{Graph: randomGraph(n, m)}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sut.Pair(graph.Node(rand.Intn(n)), graph.Node(rand.Intn(n)))
	}
}

func BenchmarkToAll(b *testing.B) {
	b.StopTimer()
	n := 100_000
	m := 4 * n

	sut := shortestpath.Dijkstra{Graph: randomGraph(n, m)}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sut.ToAll(graph.Node(rand.Intn(n)))
	}
}

func randomGraph(n, m int) graph.Graph {

	var edges []graph.Edge
	for i := 0; i < m; i++ {
		edges = append(edges, graph.Edge{
			From: graph.Node(rand.Intn(n)),
			To:   graph.Node(rand.Intn(n)),
			Cost: rand.Float64(),
		})
	}

	return graph.NewAdjacencyList(edges, n)
}
