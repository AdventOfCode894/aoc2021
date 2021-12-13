package astar

import (
	"math"
	"testing"
)

type point2D struct {
	x float64
	y float64
}

func (p1 point2D) distance(p2 point2D) float64 {
	dx := p2.x - p1.x
	dy := p2.y - p1.y
	return math.Sqrt(dx*dx + dy*dy)
}

func Test2DPoints(t *testing.T) {
	type testCase struct {
		points       []point2D
		neighbors    [][]int
		expectedPath []int
		start        int
		end          int
	}
	table := []testCase{
		{
			points: []point2D{
				{0, 0},
				{1, 0},
				{2, 0},
				{0, 1},
				{0, 2},
				{1, 3},
				{2, 2},
				{2, 1},
			},
			neighbors: [][]int{
				{1, 3},
				{0, 2},
				{1, 7},
				{0, 4},
				{3, 5},
				{4, 6},
				{5, 7},
				{2, 6},
			},
			expectedPath: []int{0, 1, 2, 7, 6},
			start:        0,
			end:          6,
		},
	}
	for i, example := range table {
		path := Search(example.start,
			func(from int) []int {
				return example.neighbors[from]
			}, func(from, to int) float64 {
				return example.points[from].distance(example.points[to])
			}, func(from int) float64 {
				return example.points[from].distance(example.points[example.end])
			}, func(at int) bool {
				return at == example.end
			})
		if len(path) != len(example.expectedPath) {
			t.Errorf("test %d failed: path had wrong length: got %d (expected: %d)", i, len(path), len(example.expectedPath))
		}
		for j := range path {
			if path[j] != example.expectedPath[j] {
				t.Errorf("test %d failed: path had wrong node at step %d: got %d (expected: %d)", i, j, path[j], example.expectedPath[j])
			}
		}
	}
}
