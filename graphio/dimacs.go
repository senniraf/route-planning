package graphio

import (
	"bufio"
	"fmt"
	"os"
	"route-planning/graph"
	"strconv"
	"strings"
)

type dimacsInput struct {
	file string
}

func NewDIMANCSInput(file string) GraphInput {
	return &dimacsInput{file: file}
}

func (di dimacsInput) LoadGraph() ([]graph.Edge, int, error) {
	f, err := os.Open(di.file)
	if err != nil {
		return nil, 0, fmt.Errorf("error opening graph file %s: %w", di.file, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	n, m, err := parseProblemLine(scanner)
	if err != nil {
		return nil, 0, fmt.Errorf("DIMACS parsing failed: \n\t%w", err)
	}

	edges, err := parseArcDescriptorLines(scanner, m)
	if err != nil {
		return nil, 0, fmt.Errorf("DIMACS parsing failed: \n\t%w", err)
	}

	return edges, n, nil
}

func parseProblemLine(scanner *bufio.Scanner) (int, int, error) {
	problemLine := nextRelevantLine(scanner)
	if problemLine[0] != 'p' {
		return 0, 0, fmt.Errorf("expected problem line but got %q", problemLine)
	}

	fields := strings.Split(problemLine, " ")
	if len(fields) != 4 || fields[1] != "sp" {
		return 0, 0, fmt.Errorf("expected problem line of format 'p sp {n} {m}', but got %s", problemLine)
	}

	n, err := strconv.Atoi(fields[2])
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing {n} in problem line: %w", err)
	}

	m, err := strconv.Atoi(fields[3])
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing {m} in problem line: %w", err)
	}

	return n, m, nil
}

func parseArcDescriptorLines(scanner *bufio.Scanner, m int) ([]graph.Edge, error) {
	line := nextRelevantLine(scanner)

	edges := make([]graph.Edge, 0, 2*m)

	for line != "" {

		if line[0] != 'a' {
			return nil, fmt.Errorf("expected arc descriptor line but got %q", line)
		}

		fields := strings.Split(line, " ")

		if len(fields) != 4 {
			return nil, fmt.Errorf("expected arc descriptor line to have format 'a {u} {v} {w}' but got %q", line)
		}

		u, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("error parsing {u} of arc line %q: %w", line, err)
		}

		v, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, fmt.Errorf("error parsing {v} of arc line %q: %w", line, err)
		}

		w, err := strconv.Atoi(fields[3])
		if err != nil {
			return nil, fmt.Errorf("error parsing {w} of arc line %q: %w", line, err)
		}

		edges = append(edges, graph.Edge{
			From: graph.Node(u - 1),
			To:   graph.Node(v - 1),
			Cost: float64(w),
		})

		line = nextRelevantLine(scanner)
	}

	return edges, nil
}

// nextRelevantLine returns the next relevant line while skipping empty lines and comments.
func nextRelevantLine(scanner *bufio.Scanner) string {
	for scanner.Scan() {
		line := scanner.Text()
		if !checkEmpty(line) && !(line[0] == 'c') {
			return line
		}
	}
	return ""
}

func checkEmpty(line string) bool {
	l := strings.Replace(line, " ", "", -1)
	return l == ""
}
