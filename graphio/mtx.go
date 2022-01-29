package graphio

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"route-planning/graph"
	"strconv"
	"strings"
)

type mtxInput struct {
	file string
}

func NewMTXInput(file string) GraphInput {
	return &mtxInput{file: file}
}

func (mi mtxInput) LoadGraph() ([]graph.Edge, int, error) {
	f, err := os.Open(mi.file)
	if err != nil {
		return nil, 0, fmt.Errorf("error opening graph file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	if err := parseHeader(scanner); err != nil {
		return nil, 0, err
	}

	n, m, err := parseSize(scanner)
	if err != nil {
		return nil, 0, err
	}

	edges, err := parseData(scanner, m)
	if err != nil {
		return nil, 0, err
	}

	return edges, n, nil
}

func parseHeader(scanner *bufio.Scanner) error {
	if !scanner.Scan() {
		return errors.New("empty file")
	}

	headerLine := scanner.Text()

	if headerLine[0] != '%' {
		return fmt.Errorf("header line expected to begin with '%%', but got '%s'", string(headerLine[0]))
	}

	toParse := headerLine[1:]
	if toParse[0] == '%' {
		toParse = toParse[1:]
	}

	if toParse != "MatrixMarket matrix coordinate pattern symmetric" {
		return fmt.Errorf("unsupported mtx format: %s", toParse)
	}

	return nil
}

func parseSize(scanner *bufio.Scanner) (int, int, error) {
	sizeLine, err := skipComments(scanner)
	if err != nil {
		return 0, 0, err
	}

	fields := strings.Split(sizeLine, " ")

	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing size line '%s': %w", sizeLine, err)
	}

	m, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing size line '%s': %w", sizeLine, err)
	}

	if m != n {
		return 0, 0, fmt.Errorf("error parsing size line '%s': expected m to equal n", sizeLine)
	}

	nonzeros, err := strconv.Atoi(fields[2])
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing size line '%s': %w", sizeLine, err)
	}

	return n, nonzeros, nil
}

func parseData(scanner *bufio.Scanner, m int) ([]graph.Edge, error) {
	edges := make([]graph.Edge, 0, 2*m)

	for scanner.Scan() {
		nodes := strings.Split(scanner.Text(), " ")

		v, err := strconv.Atoi(nodes[0])
		if err != nil {
			return nil, fmt.Errorf("unable to parse '%s': %w", nodes[0], err)
		}
		w, err := strconv.Atoi(nodes[1])
		if err != nil {
			return nil, fmt.Errorf("unable to parse '%s': %w", nodes[1], err)
		}

		edges = append(edges, graph.Edge{
			From: graph.Node(v - 1),
			To:   graph.Node(w - 1),
			Cost: 1.0,
		}, graph.Edge{
			From: graph.Node(w - 1),
			To:   graph.Node(v - 1),
			Cost: 1.0,
		})

	}

	if numEdges := len(edges); numEdges != 2*m {
		return nil, fmt.Errorf("expected mtx file to contain %d data lines but only got %d", m, numEdges/2)
	}

	return edges, nil
}

func skipComments(scanner *bufio.Scanner) (string, error) {
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if line[0] != '%' {
			return line, nil
		}
	}

	return "", errors.New("mtx only contains header and comments")
}
