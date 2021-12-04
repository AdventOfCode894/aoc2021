package aoc2021d4

import "fmt"

type BingoCard struct {
	cells         [][]uint64 // Row-major
	marked        [][]bool
	cellLocations map[uint64]cardLocation
	won           bool

	width  int
	height int
}

type cardLocation struct {
	row int
	col int
}

func newBingoCard(cells []uint64, width int) (*BingoCard, error) {
	if len(cells)%width != 0 {
		return nil, fmt.Errorf("bingo card has impossible geometry: width = %d but cell count = %d", width, len(cells))
	}
	height := len(cells) / width

	bc := &BingoCard{
		cells:         make([][]uint64, height),
		marked:        make([][]bool, height),
		cellLocations: make(map[uint64]cardLocation),
		width:         width,
		height:        height,
	}
	for row := 0; row < height; row++ {
		bc.cells[row] = make([]uint64, width)
		bc.marked[row] = make([]bool, width)
		for col := 0; col < width; col++ {
			cell := cells[0]
			cells = cells[1:]
			bc.cells[row][col] = cell
			bc.cellLocations[cell] = cardLocation{row: row, col: col}
		}
	}

	return bc, nil
}

func (bc *BingoCard) Mark(draw uint64) (winner bool) {
	loc, found := bc.cellLocations[draw]
	if found {
		bc.marked[loc.row][loc.col] = true
		if bc.checkWin(loc.row, 0, 0, 1) ||
			bc.checkWin(0, loc.col, 1, 0) {
			bc.won = true
		}
	}
	return bc.won
}

func (bc *BingoCard) Won() bool { return bc.won }

func (bc *BingoCard) UnmarkedSum() uint64 {
	sum := uint64(0)
	for row := 0; row < bc.height; row++ {
		for col := 0; col < bc.width; col++ {
			if !bc.marked[row][col] {
				sum += bc.cells[row][col]
			}
		}
	}
	return sum
}

func (bc *BingoCard) checkWin(startRow int, startCol int, stepRow int, stepCol int) bool {
	for startRow < bc.height && startCol < bc.width {
		if !bc.marked[startRow][startCol] {
			return false
		}
		startRow += stepRow
		startCol += stepCol
	}
	return true
}
