package graph

import (
	"reflect"
	"testing"
)

func TestNewAdjacencyList(t *testing.T) {
	type args struct {
		edges []Edge
		n     int
	}
	tests := []struct {
		name string
		args args
		want Graph
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAdjacencyList(tt.args.edges, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAdjacencyList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_adjacencyList_OutgoingEdges(t *testing.T) {
	type args struct {
		v Node
	}
	tests := []struct {
		name string
		l    adjacencyList
		args args
		want []Edge
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.OutgoingEdges(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("adjacencyList.OutgoingEdges() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_adjacencyList_N(t *testing.T) {
	tests := []struct {
		name string
		l    adjacencyList
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.N(); got != tt.want {
				t.Errorf("adjacencyList.N() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_adjacencyList_Reverted(t *testing.T) {
	tests := []struct {
		name string
		l    adjacencyList
		want Graph
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Reverted(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("adjacencyList.Reverted() = %v, want %v", got, tt.want)
			}
		})
	}
}
