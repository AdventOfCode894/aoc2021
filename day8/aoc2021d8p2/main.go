package main

import (
	"errors"
	"io"
	"strconv"

	"github.com/AdventOfCode894/aoc2021/internal/aocio"

	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

func solvePuzzle(r io.Reader) (int, error) {
	// Sorry for this mess! Yikes.
	total := 0
	pr := aocio.NewPuzzleReader(r)
	for pr.NextNonEmptyLine() {
		tokens := pr.ReadStringArrayLine(' ')
		if err := pr.Err(); err != nil {
			return 0, err
		}
		if len(tokens) != 10+1+4 || tokens[10] != "|" {
			return 0, errors.New("invalid input data")
		}
		patterns := tokens[:10]
		outputs := tokens[11:]

		var right1, right2, top, midL1, midL2, botL1, botL2 rune
		// 1
		for _, p := range patterns {
			if len(p) != 2 {
				continue
			}
			right1 = rune(p[0])
			right2 = rune(p[1])
		}
		// 7
		for _, p := range patterns {
			if len(p) != 3 {
				continue
			}
			for _, c := range p {
				if c != right1 && c != right2 {
					top = c
					break
				}
			}
		}
		// 4
		for _, p := range patterns {
			if len(p) != 4 {
				continue
			}
			for _, c := range p {
				if c != right1 && c != right2 {
					midL1, midL2 = midL2, c
				}
			}
		}
		// 8
		for _, p := range patterns {
			if len(p) != 7 {
				continue
			}
			for _, c := range p {
				if c != right1 && c != right2 && c != top && c != midL1 && c != midL2 {
					botL1, botL2 = botL2, c
				}
			}
		}

		result := ""
		for _, o := range outputs {
			switch len(o) {
			case 2: //1
				result += "1"
			case 4: //4
				result += "4"
			case 3: //7
				result += "7"
			case 7: //8
				result += "8"
			case 5: // 2 or 3 or 5
				botLs := 0
				rights := 0
				for _, c := range o {
					if c == botL1 || c == botL2 {
						botLs++
					}
					if c == right1 || c == right2 {
						rights++
					}
				}
				if botLs > 1 {
					result += "2"
					break
				}
				if rights > 1 {
					result += "3"
				} else {
					result += "5"
				}
			case 6: // 0 or 6 or 9
				botLs := 0
				rights := 0
				for _, c := range o {
					if c == botL1 || c == botL2 {
						botLs++
					}
					if c == right1 || c == right2 {
						rights++
					}
				}
				if botLs < 2 {
					result += "9"
					break
				}
				if rights > 1 {
					result += "0"
				} else {
					result += "6"
				}
			default:
				panic("impossible digit")
			}
		}
		numericResult, _ := strconv.Atoi(result)
		total += numericResult
	}

	return total, nil
}
