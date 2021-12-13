package astar

import (
	"container/heap"
)

func Search(start int,
	neighbors func(from int) []int,
	distance func(from, to int) float64,
	estimate func(from int) float64,
	goal func(at int) bool) (path []int) {

	var ss searchState
	heap.Init(&ss)
	ss.visited = map[int]visitedNode{start: {}}

	node := start
	for !goal(node) {
		for _, next := range neighbors(node) {
			nextDistance := distance(node, next)
			if nextDistance < 0 {
				nextDistance = 0
			}
			nextEstimate := estimate(next)
			if nextEstimate < 0 {
				nextEstimate = 0
			}
			fromStart := ss.visited[node].shortestFromStart + nextDistance
			v, seen := ss.visited[next]
			if seen && fromStart >= v.shortestFromStart {
				continue
			}
			if !v.scheduledForVisit {
				ss.pushNode(next, fromStart+nextEstimate)
			}
			ss.visited[next] = visitedNode{
				shortestFromStart: fromStart,
				previousNode:      node,
				scheduledForVisit: true,
			}
		}
		if len(ss.priorityQueue) < 1 {
			// No path to destination
			return nil
		}
		node = ss.popNode()
		v := ss.visited[node]
		v.scheduledForVisit = false
		ss.visited[node] = v
	}
	if !goal(node) {
		return nil
	}

	path = []int{node}
	for node != start {
		node = ss.visited[node].previousNode
		path = append(path, node)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

type searchState struct {
	priorityQueue []searchNode
	nextNodeOrder uint
	visited       map[int]visitedNode
}

type searchNode struct {
	index             int
	estimatedDistance float64
	fifoOrder         uint
}

type visitedNode struct {
	shortestFromStart float64
	previousNode      int
	scheduledForVisit bool
}

func (ss *searchState) Len() int { return len(ss.priorityQueue) }

func (ss *searchState) Less(i, j int) bool {
	n1 := ss.priorityQueue[i]
	n2 := ss.priorityQueue[j]
	di := n1.estimatedDistance
	dj := n2.estimatedDistance
	if di < dj {
		return true
	}
	// Tied nodes are sorted in FIFO order, improving performance by degrading to DFS
	if di == dj && n1.fifoOrder < n2.fifoOrder {
		return true
	}
	return false
}

func (ss *searchState) Swap(i, j int) {
	ss.priorityQueue[i], ss.priorityQueue[j] = ss.priorityQueue[j], ss.priorityQueue[i]
}

func (ss *searchState) Push(x interface{}) {
	ss.priorityQueue = append(ss.priorityQueue, x.(searchNode))
}

func (ss *searchState) Pop() interface{} {
	n := len(ss.priorityQueue) - 1
	x := ss.priorityQueue[n]
	ss.priorityQueue = ss.priorityQueue[:n]
	return x
}

func (ss *searchState) pushNode(index int, estimatedDistance float64) {
	heap.Push(ss, searchNode{
		index:             index,
		estimatedDistance: estimatedDistance,
		fifoOrder:         ss.nextNodeOrder,
	})
	ss.nextNodeOrder++
}

func (ss *searchState) popNode() (index int) {
	return heap.Pop(ss).(searchNode).index
}
