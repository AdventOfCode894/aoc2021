package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	buf := bufio.NewReader(os.Stdin)
	if err := solvePuzzle(buf, os.Stdout); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func solvePuzzle(r io.RuneReader, w io.Writer) error {
	var zeroes, ones []uint
	i := 0
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return fmt.Errorf("failed to read from input: %v", err)
			}
			break
		}
		switch c {
		case '0', '1':
			if i >= len(zeroes) {
				zeroes = append(zeroes, 0)
				ones = append(ones, 0)
			}
			if c == '0' {
				zeroes[i]++
			} else {
				ones[i]++
			}
			i++
		case '\r', '\n':
			i = 0
		default:
			return fmt.Errorf("invalid character in input: U+%d", c)
		}
	}

	gamma := uint(0)
	epsilon := uint(0)
	for i := range zeroes {
		gamma <<= 1
		epsilon <<= 1
		if zeroes[i] > ones[i] {
			epsilon |= 1
		} else {
			gamma |= 1
		}
	}
	power := gamma * epsilon

	_, _ = fmt.Fprintf(w, "Power consumption: %d\n", power)
	return nil
}
