package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	var in io.Reader = os.Stdin
	if len(os.Args) == 2 {
		var err error
		if in, err = os.Open(os.Args[1]); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error attempting to open input file \"%s\": %v", os.Args[0], err)
		}
	}
	if err := solvePuzzle(in, os.Stdout); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func solvePuzzle(r io.Reader, w io.Writer) error {
	// Sorry for this mess! Yikes.
	total := 0
	for {
		var patterns []string
		for i := 0; i < 10; i++ {
			var pattern string
			if _, err := fmt.Fscanf(r, "%s", &pattern); err != nil {
				if !errors.Is(err, io.EOF) {
					return fmt.Errorf("failed to read pattern: %v", err)
				}
				break
			}
			patterns = append(patterns, pattern)
		}
		var delim rune
		if _, err := fmt.Fscanf(r, "%c", &delim); err != nil {
			if !errors.Is(err, io.EOF) {
				return fmt.Errorf("failed to read delimeter: %v", err)
			}
			break
		}
		if delim != '|' {
			return fmt.Errorf("wrong delimeter: %c", delim)
		}
		var outputs []string
		for i := 0; i < 4; i++ {
			var output string
			if _, err := fmt.Fscanf(r, "%s", &output); err != nil {
				if !errors.Is(err, io.EOF) {
					return fmt.Errorf("failed to read output: %v", err)
				}
				break
			}
			outputs = append(outputs, output)
		}

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
		fmt.Println(result, numericResult)
		total += numericResult
	}

	_, _ = fmt.Fprintf(w, "Decoded output total %d\n", total)
	return nil
}
