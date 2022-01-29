package main

import (
	"fmt"
	"os"
	"route-planning/graph"
	"route-planning/graphio"
	"route-planning/shortestpath"
	"time"

	"github.com/jessevdk/go-flags"
)

type FileArg struct {
	File string
}

type DijkstraCommand struct {
	Source int `short:"s" long:"source" required:"true" description:"source node"`
	Target int `short:"t" long:"target" required:"true" description:"target node"`

	Bidirectional bool `short:"b" long:"bidirect" description:"use bidirectional mode"`

	FileArg FileArg `positional-args:"true" required:"true"`
}

var cli struct {
	Dijkstra DijkstraCommand `command:"dijkstra"`

	Format  string `short:"f" long:"format" description:"the input format" choice:"mtx" choice:"dimacs" default:"dimacs"`
	Verbose bool   `short:"v" long:"verbose" description:"display additional information"`
}

func main() {
	p := flags.NewParser(&cli, flags.Default)
	if _, err := p.Parse(); err != nil {
		if flags.WroteHelp(err) {
			return
		}
		p.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	var in graphio.GraphInput
	switch cli.Format {
	case "dimacs":
		in = graphio.NewDIMANCSInput(cli.Dijkstra.FileArg.File)
	case "mtx":
		in = graphio.NewMTXInput(cli.Dijkstra.FileArg.File)
	}

	edges, n, err := in.LoadGraph()
	if err != nil {
		fmt.Printf("Error loading graph from input: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded %d edges and %d nodes\n", len(edges), n)

	g := graph.NewAdjacencyList(edges, n)

	var algo shortestpath.Algorithm
	algo = &shortestpath.Dijkstra{Graph: g}

	if cli.Dijkstra.Bidirectional {
		algo = &shortestpath.BidirectDijkstra{
			ForwardGraph:  g,
			BackwardGraph: g.Reverted(),
		}
	}

	s, t := graph.Node(cli.Dijkstra.Source-1), graph.Node(cli.Dijkstra.Target-1)

	start := time.Now()
	c, path := algo.PairShortestPath(s, t)
	took := time.Since(start)

	fmt.Printf("Shortest path algorithm took: %v\n", took)

	if len(path) == 0 {
		fmt.Printf("No path from %d to %d\n", s+1, t+1)
		os.Exit(1)
	}

	fmt.Printf("Cost: %f, Hops: %d\n", c, len(path))

	if cli.Verbose {
		pathNodes := []graph.Node{path[0].From + 1}
		for _, e := range path {
			pathNodes = append(pathNodes, e.To+1)
		}
		fmt.Printf("%v\n", pathNodes)
	}

}
