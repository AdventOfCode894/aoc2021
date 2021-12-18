package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/AdventOfCode894/aoc2021/internal/aocio"
	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

type snailfishNumberType int

const (
	snailfishRegular snailfishNumberType = iota
	snailfishPair
)

var ErrInvalidNumber = errors.New("invalid snailfish number")

type snailfishNumber struct {
	numType snailfishNumberType
	regular uint
	left    *snailfishNumber
	right   *snailfishNumber
}

func (n *snailfishNumber) Add(a *snailfishNumber, b *snailfishNumber) *snailfishNumber {
	n.numType = snailfishPair
	n.left = a.clone()
	n.right = b.clone()
	for {
		if _, _, more := n.doReduce(0, false); more {
			continue
		}
		if _, _, more := n.doReduce(0, true); more {
			continue
		}
		break
	}
	return n
}

func (n *snailfishNumber) clone() *snailfishNumber {
	if n.numType == snailfishRegular {
		return &snailfishNumber{numType: snailfishRegular, regular: n.regular}
	} else {
		return &snailfishNumber{
			numType: snailfishPair,
			left:    n.left.clone(),
			right:   n.right.clone(),
		}
	}
}

func (n *snailfishNumber) String() string {
	if n.numType == snailfishRegular {
		return strconv.FormatUint(uint64(n.regular), 10)
	}
	return fmt.Sprintf("[%s,%s]", n.left, n.right)
}

func (n *snailfishNumber) doReduce(depth uint, doSplits bool) (addLeft uint, addRight uint, didOp bool) {
	if doSplits {
		if n.numType == snailfishRegular && n.regular >= 10 {
			left := n.regular / 2
			right := (n.regular + 1) / 2
			n.numType = snailfishPair
			n.left = &snailfishNumber{numType: snailfishRegular, regular: left}
			n.right = &snailfishNumber{numType: snailfishRegular, regular: right}
			return 0, 0, true
		}
	} else {
		if depth >= 4 && n.numType == snailfishPair {
			n.numType = snailfishRegular
			n.regular = 0
			return n.left.regular, n.right.regular, true
		}
	}
	if n.numType == snailfishPair {
		if addLeft, addRight, didOp := n.left.doReduce(depth+1, doSplits); didOp {
			if addRight > 0 {
				n.right.addLeftmost(addRight)
			}
			return addLeft, 0, true
		}
		if addLeft, addRight, didOp := n.right.doReduce(depth+1, doSplits); didOp {
			if addLeft > 0 {
				n.left.addRightmost(addLeft)
			}
			return 0, addRight, true
		}
	}
	return 0, 0, false
}

func (n *snailfishNumber) addLeftmost(v uint) {
	if n.numType == snailfishRegular {
		n.regular += v
		return
	}
	n.left.addLeftmost(v)
}

func (n *snailfishNumber) addRightmost(v uint) {
	if n.numType == snailfishRegular {
		n.regular += v
		return
	}
	n.right.addRightmost(v)
}

func (n *snailfishNumber) Magnitude() uint {
	if n.numType == snailfishPair {
		return 3*n.left.Magnitude() + 2*n.right.Magnitude()
	} else {
		return n.regular
	}
}

func solvePuzzle(r io.Reader) (int, error) {
	pr := aocio.NewPuzzleReader(r)
	var numbers []*snailfishNumber
	for pr.NextNonEmptyLine() {
		line := pr.LineRunes()
		n, _, err := readSnailfishNumber(line)
		if err != nil {
			return 0, err
		}
		numbers = append(numbers, n)
	}
	if err := pr.Err(); err != nil {
		return 0, err
	}

	maxMag := uint(0)
	var maxA, maxB, maxSum *snailfishNumber
	sum := new(snailfishNumber)
	for _, a := range numbers {
		for _, b := range numbers {
			if a == b {
				continue
			}
			sum.Add(a, b)
			mag := sum.Magnitude()
			if mag > maxMag {
				maxMag = mag
				maxA = a
				maxB = b
			}
			maxSum = sum.clone()
		}
	}

	fmt.Printf("  %s\n+ %s\n= %s\n\n", maxA, maxB, maxSum)

	return int(maxMag), nil
}

func readSnailfishNumber(b []rune) (*snailfishNumber, []rune, error) {
	if b[0] == '[' {
		left, b, err := readSnailfishNumber(b[1:])
		if err != nil {
			return nil, nil, err
		}
		if b[0] != ',' {
			return nil, nil, ErrInvalidNumber
		}
		right, b, err := readSnailfishNumber(b[1:])
		if err != nil {
			return nil, nil, err
		}
		if b[0] != ']' {
			return nil, nil, ErrInvalidNumber
		}
		n := &snailfishNumber{numType: snailfishPair, left: left, right: right}
		return n, b[1:], nil
	} else {
		i := 0
		for {
			if b[i] >= '0' && b[i] <= '9' {
				i++
				continue
			}
			if b[i] == ',' || b[i] == ']' {
				break
			}
			return nil, nil, ErrInvalidNumber
		}
		v, err := strconv.ParseUint(string(b[:i]), 10, 64)
		if err != nil {
			return nil, nil, ErrInvalidNumber
		}
		n := &snailfishNumber{numType: snailfishRegular, regular: uint(v)}
		return n, b[i:], nil
	}
}
