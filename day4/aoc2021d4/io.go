package aoc2021d4

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func nextLine(s *bufio.Scanner) error {
	for s.Scan() && len(s.Bytes()) < 1 {
		// Skip empty line
	}
	return s.Err()
}

func ReadDraws(s *bufio.Scanner) ([]uint64, error) {
	if err := nextLine(s); err != nil {
		return nil, err
	}

	sequenceNums := strings.Split(s.Text(), ",")
	draws := make([]uint64, len(sequenceNums))
	for i, num := range sequenceNums {
		draw, err := strconv.ParseUint(num, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid draw number \"%s\": %v", num, err)
		}
		draws[i] = draw
	}

	return draws, nil
}

func ReadBingoCard(s *bufio.Scanner) (*BingoCard, error) {
	if err := nextLine(s); err != nil {
		return nil, err
	}
	if len(s.Text()) < 1 {
		return nil, nil
	}

	cells, err := appendCells(nil, s.Text())
	if err != nil {
		return nil, err
	}
	width := len(cells)

	for s.Scan() {
		if len(s.Text()) < 1 {
			break
		}
		if cells, err = appendCells(cells, s.Text()); err != nil {
			return nil, err
		}
	}

	return newBingoCard(cells, width)
}

func appendCells(cells []uint64, str string) ([]uint64, error) {
	r := strings.NewReader(str)
	for {
		var cell uint64
		if _, err := fmt.Fscanf(r, "%d", &cell); err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, fmt.Errorf("bingo card line contained invalid number: \"%s\": %v", str, err)
			}
			break
		}
		cells = append(cells, cell)
	}
	return cells, nil
}
