package main

type diagnosticTreeNode struct {
	isValue      bool
	subtreeCount uint
	child        [2]*diagnosticTreeNode
}

func (n *diagnosticTreeNode) Insert(x uint64, bitLen int) {
	n.subtreeCount++

	if bitLen < 1 {
		n.isValue = true
		return
	}

	bitLen--
	b := (x >> bitLen) & 1
	x = x & ((1 << bitLen) - 1)

	if n.child[b] == nil {
		n.child[b] = new(diagnosticTreeNode)
	}
	n.child[b].Insert(x, bitLen)
}

func (n *diagnosticTreeNode) FindRating(mostPopularPrefix bool) uint64 {
	return n.findRatingRecurse(mostPopularPrefix, 0)
}

func (n *diagnosticTreeNode) findRatingRecurse(mostPopularPrefix bool, partial uint64) uint64 {
	if n.isValue {
		return partial
	}
	var mostPopular, leastPopular int
	if n.child[0] == nil {
		mostPopular = 1
		leastPopular = 1
	} else if n.child[1] == nil {
		mostPopular = 0
		leastPopular = 0
	} else {
		mostPopular = 1
		leastPopular = 0
		// In case of a tie, prefer "1" for most popular search or "0" for least popular search
		if n.child[0].subtreeCount > n.child[1].subtreeCount {
			mostPopular = 0
			leastPopular = 1
		}
	}

	var idx int
	if mostPopularPrefix {
		idx = mostPopular
	} else {
		idx = leastPopular
	}
	partial = (partial << 1) | uint64(idx)
	return n.child[idx].findRatingRecurse(mostPopularPrefix, partial)
}
