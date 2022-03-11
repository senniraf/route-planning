package priorityqueue_test

import (
	"math/rand"
	"route-planning/graph"
	"route-planning/priorityqueue"
	"testing"
)

var testData = []priorityqueue.Element{
	{Node: 0, Cost: 17.0},
	{Node: 1, Cost: 3.0},
	{Node: 2, Cost: 2.0},
	{Node: 3, Cost: 14.0},
	{Node: 5, Cost: 22.0},
	{Node: 4, Cost: 12.0},
}

var order = []int{2, 1, 5, 3, 0, 4}

func TestPush(t *testing.T) {
	n := 17
	rand.Seed(42)

	sut := priorityqueue.NewMinHeap()

	for i := 0; i < n; i++ {
		sut.Push(priorityqueue.Element{
			Node: graph.Node(rand.Int()),
			Cost: rand.Float64(),
		})
	}

	if sut.Len() != n {
		t.Fatalf("Len expected to be %d, is %d", n, sut.Len())
	}
}

func TestPop(t *testing.T) {
	sut := priorityqueue.NewMinHeap()

	for _, e := range testData {
		sut.Push(e)
	}

	for i, idx := range order {
		expected := testData[idx]

		got := sut.Pop()
		if expected != got {
			t.Errorf("Pop %d: expected %+v, got %+v", i, expected, got)
		}

		if expectedLen := len(testData) - (i + 1); sut.Len() != expectedLen {
			t.Errorf("Pop %d: expected len %d, got %d", i, expectedLen, sut.Len())
		}
	}
}

func TestTop(t *testing.T) {
	sut := priorityqueue.NewMinHeap()

	for _, e := range testData {
		sut.Push(e)
	}

	for i, idx := range order {
		expected := testData[idx]

		top := sut.Top()

		if expected != top {
			t.Errorf("Top %d: expected %+v, got %+v", i, expected, top)
		}

		if expectedLen := len(testData) - i; sut.Len() != expectedLen {
			t.Errorf("Top %d: expected len %d, got %d", i, expectedLen, sut.Len())
		}

		popped := sut.Pop()
		if expected != popped {
			t.Errorf("Pop %d: expected %+v, got %+v", i, expected, popped)
		}

		if expectedLen := len(testData) - (i + 1); sut.Len() != expectedLen {
			t.Errorf("Pop %d: expected len %d, got %d", i, expectedLen, sut.Len())
		}
	}
}
