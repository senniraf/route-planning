package graph

type listEntry struct {
	node Node
	cost float64
}

type adjacencyList [][]listEntry

func NewAdjacencyList(edges []Edge, n int) Graph {
	list := make(adjacencyList, n)

	for _, edge := range edges {
		neighbours := list[edge.From]
		list[edge.From] = append(neighbours, listEntry{
			node: edge.To,
			cost: edge.Cost,
		})
	}

	return list
}

func (l adjacencyList) OutgoingEdges(v Node) []Edge {
	neighbours := l[v]
	edges := make([]Edge, len(neighbours))
	for i, w := range neighbours {
		edges[i] = Edge{
			From: v,
			To:   w.node,
			Cost: w.cost,
		}
	}

	return edges
}

func (l adjacencyList) N() int {
	return len(l)
}

func (l adjacencyList) Reverted() Graph {
	reverted := make(adjacencyList, l.N())

	for i := 0; i < l.N(); i++ {
		v := Node(i)
		for _, outgoing := range l.OutgoingEdges(v) {
			neighbours := reverted[outgoing.To]
			reverted[outgoing.To] = append(neighbours, listEntry{
				node: outgoing.From,
				cost: outgoing.Cost,
			})
		}
	}

	return reverted
}
