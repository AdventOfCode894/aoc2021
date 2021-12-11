package aocmain

import (
	"fmt"
	"io"
	"os"
)

func HandlePuzzle(solver func(r io.Reader) (int, error)) {
	var in io.Reader = os.Stdin
	if len(os.Args) == 2 {
		var err error
		if in, err = os.Open(os.Args[1]); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error attempting to open input file \"%s\": %v", os.Args[1], err)
			os.Exit(1)
		}
	}
	solution, err := solver(in)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(solution)
}
