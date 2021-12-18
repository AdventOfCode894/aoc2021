package aocio

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type PuzzleReader struct {
	s       *bufio.Scanner
	scanned bool
	err     error

	tokenReader *TokenReader
}

type TokenReader struct {
	b   []byte
	err *error
}

const MaxPuzzleLineLen = 1 * 1024 * 1024
const initialLineBufLen = 4 * 1024

func NewPuzzleReader(r io.Reader) *PuzzleReader {
	pr := &PuzzleReader{
		s: bufio.NewScanner(r),
	}
	pr.s.Buffer(make([]byte, initialLineBufLen), MaxPuzzleLineLen)
	return pr
}

func (pr *PuzzleReader) ensureScanned() {
	if pr.scanned {
		return
	}
	pr.NextLine()
}

func (pr *PuzzleReader) Err() error {
	if pr.err != nil {
		return pr.err
	}
	return pr.s.Err()
}

func (pr *PuzzleReader) NextLine() bool {
	if pr.Err() != nil {
		return false
	}
	pr.scanned = true
	return pr.s.Scan()
}

func (pr *PuzzleReader) SkipEmptyLines() {
	if pr.Err() != nil {
		return
	}
	for pr.IsLineEmpty() {
		pr.NextLine()
	}
}

func (pr *PuzzleReader) NextNonEmptyLine() bool {
	for {
		if !pr.NextLine() {
			return false
		}
		if pr.IsLineEmpty() {
			continue
		}
		break
	}
	return true
}

func (pr *PuzzleReader) LineReader() *bytes.Reader {
	pr.ensureScanned()
	if pr.Err() != nil {
		return new(bytes.Reader)
	}
	return bytes.NewReader(pr.s.Bytes())
}

func (pr *PuzzleReader) LineTokenReader() *TokenReader {
	pr.ensureScanned()
	if err := pr.Err(); err != nil {
		return &TokenReader{b: nil, err: &err}
	}
	return &TokenReader{b: pr.s.Bytes(), err: &pr.err}
}

func (pr *PuzzleReader) LineRunes() []rune {
	pr.ensureScanned()
	if pr.Err() != nil {
		return nil
	}
	return bytes.Runes(pr.s.Bytes())
}

func (pr *PuzzleReader) LineString() string {
	pr.ensureScanned()
	if pr.Err() != nil {
		return ""
	}
	return pr.s.Text()
}

func (pr *PuzzleReader) LineLen() int {
	pr.ensureScanned()
	if pr.Err() != nil {
		return 0
	}
	return len(pr.s.Bytes())
}

func (pr *PuzzleReader) NewTokenReader(b []byte) *TokenReader {
	return &TokenReader{b: b}
}

func (pr *PuzzleReader) IsLineEmpty() bool {
	br := pr.LineReader()
	for {
		c, _, err := br.ReadRune()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return false
			}
			break
		}
		if !unicode.IsSpace(c) {
			return false
		}
	}
	return true
}

func (tr *TokenReader) NextRune() (rune, bool) {
	if *tr.err != nil {
		return utf8.RuneError, false
	}
	if len(tr.b) < 1 {
		return utf8.RuneError, false
	}
	c, sz := utf8.DecodeRune(tr.b)
	if c == utf8.RuneError {
		*tr.err = errors.New("input file contained invalid UTF-8")
		return utf8.RuneError, false
	}
	tr.b = tr.b[sz:]
	return c, true
}

func (tr *TokenReader) PeekRune() (rune, bool) {
	if *tr.err != nil {
		return utf8.RuneError, false
	}
	if len(tr.b) < 1 {
		return utf8.RuneError, false
	}
	c, _ := utf8.DecodeRune(tr.b)
	if c == utf8.RuneError {
		*tr.err = errors.New("input file contained invalid UTF-8")
		return utf8.RuneError, false
	}
	return c, true
}

const (
	EOLDelim rune = utf8.MaxRune + 1
	NoDelim  rune = utf8.MaxRune + 2
)

func (tr *TokenReader) NextToken(delim rune) ([]byte, bool) {
	if *tr.err != nil {
		return nil, false
	}
	if len(tr.b) < 1 {
		return nil, false
	}
	var idx int
	skip := 0
	switch delim {
	case EOLDelim:
		idx = len(tr.b)
	case NoDelim:
		_, idx = utf8.DecodeRune(tr.b)
	default:
		idx = bytes.IndexRune(tr.b, delim)
		if idx < 0 {
			idx = len(tr.b)
		} else {
			skip = utf8.RuneLen(delim)
		}
	}
	token := tr.b[:idx]
	tr.b = tr.b[idx+skip:]
	return token, true
}

func (tr *TokenReader) NextString(delim rune) (string, bool) {
	token, ok := tr.NextToken(delim)
	return string(token), ok
}

func (tr *TokenReader) ConsumeRepeating(cutset string) {
	if *tr.err != nil {
		return
	}
	tr.b = bytes.TrimLeft(tr.b, cutset)
}

func (tr *TokenReader) ConsumeSpaces() {
	if *tr.err != nil {
		return
	}
	tr.b = bytes.TrimLeftFunc(tr.b, func(c rune) bool {
		return unicode.IsSpace(c)
	})
}

func (tr *TokenReader) ConsumeString(s string) bool {
	if *tr.err != nil {
		return false
	}
	b := tr.b
	tr.b = bytes.TrimPrefix(tr.b, []byte(s))
	if len(tr.b) >= len(b) {
		*tr.err = fmt.Errorf("failed to consume fixed string \"%s\"", s)
		return false
	}
	return true
}

func (tr *TokenReader) ConsumeEOL() bool {
	if *tr.err != nil {
		return false
	}
	tr.ConsumeSpaces()
	if len(tr.b) > 0 {
		*tr.err = fmt.Errorf("found garbage at end of line: \"%s\"", tr.b)
		return false
	}
	return true
}

func (tr *TokenReader) NextUint(delim rune, base int) (uint, bool) {
	token, ok := tr.NextToken(delim)
	if !ok {
		return 0, false
	}
	x, err := strconv.ParseUint(string(token), base, 32)
	if err != nil {
		*tr.err = err
		return 0, false
	}
	return uint(x), true
}

func (tr *TokenReader) NextInt(delim rune, base int) (int, bool) {
	token, ok := tr.NextToken(delim)
	if !ok {
		return 0, false
	}
	x, err := strconv.ParseInt(string(token), base, 32)
	if err != nil {
		*tr.err = err
		return 0, false
	}
	return int(x), true
}

func (tr *TokenReader) WalkArray(split rune, f func(b []byte) error) error {
	for {
		tr.ConsumeSpaces()
		b, ok := tr.NextToken(split)
		if !ok {
			break
		}
		if err := f(b); err != nil {
			*tr.err = fmt.Errorf("error during array walk: %w", err)
			break
		}
	}
	return *tr.err
}

func (tr *TokenReader) NextIntArray(split rune, term rune, base int) []int {
	token, ok := tr.NextToken(term)
	if !ok {
		return nil
	}
	pieces := &TokenReader{b: token, err: tr.err}
	var out []int
	if err := pieces.WalkArray(split, func(b []byte) error {
		x, err := strconv.ParseInt(string(token), base, 32)
		if err != nil {
			return err
		}
		out = append(out, int(x))
		return nil
	}); err != nil {
		return nil
	}
	return out
}

func (tr *TokenReader) NextUintArray(split rune, term rune, base int) []uint {
	token, ok := tr.NextToken(term)
	if !ok {
		return nil
	}
	pieces := &TokenReader{b: token, err: tr.err}
	var out []uint
	if err := pieces.WalkArray(split, func(b []byte) error {
		x, err := strconv.ParseUint(string(b), base, 32)
		if err != nil {
			return err
		}
		out = append(out, uint(x))
		return nil
	}); err != nil {
		return nil
	}
	return out
}

func (tr *TokenReader) NextStringArray(split rune, term rune) []string {
	token, ok := tr.NextToken(term)
	if !ok {
		return nil
	}
	pieces := &TokenReader{b: token, err: tr.err}
	var out []string
	if err := pieces.WalkArray(split, func(b []byte) error {
		out = append(out, string(b))
		return nil
	}); err != nil {
		return nil
	}
	return out
}

func (pr *PuzzleReader) readFullLine(f func(tr *TokenReader)) {
	tr := pr.LineTokenReader()
	tr.ConsumeSpaces()
	f(tr)
	tr.ConsumeEOL()
}

func (pr *PuzzleReader) ReadIntLine(base int) int {
	var x int
	pr.readFullLine(func(tr *TokenReader) {
		x, _ = tr.NextInt(EOLDelim, base)
	})
	return x
}

func (pr *PuzzleReader) ReadUintLine(base int) uint {
	var x uint
	pr.readFullLine(func(tr *TokenReader) {
		x, _ = tr.NextUint(EOLDelim, base)
	})
	return x
}

func (pr *PuzzleReader) ReadIntArrayLine(split rune, base int) []int {
	var x []int
	pr.readFullLine(func(tr *TokenReader) {
		x = tr.NextIntArray(split, EOLDelim, base)
	})
	return x
}

func (pr *PuzzleReader) ReadUintArrayLine(split rune, base int) []uint {
	var x []uint
	pr.readFullLine(func(tr *TokenReader) {
		x = tr.NextUintArray(split, EOLDelim, base)
	})
	return x
}

func (pr *PuzzleReader) ReadStringArrayLine(split rune) []string {
	var x []string
	pr.readFullLine(func(tr *TokenReader) {
		x = tr.NextStringArray(split, EOLDelim)
	})
	return x
}

func (pr *PuzzleReader) Read2DIntArray(split rune, base int) (arr [][]int, width int, height int) {
	for !pr.IsLineEmpty() {
		row := pr.ReadIntArrayLine(split, base)
		arr = append(arr, row)
		if !pr.NextLine() {
			break
		}
	}
	if len(arr) < 1 || pr.Err() != nil {
		return nil, 0, 0
	}
	return arr, len(arr[0]), len(arr)
}

func (pr *PuzzleReader) Read2DUintArray(split rune, base int) (arr [][]uint, width int, height int) {
	for !pr.IsLineEmpty() {
		row := pr.ReadUintArrayLine(split, base)
		arr = append(arr, row)
		if !pr.NextLine() {
			break
		}
	}
	if len(arr) < 1 || pr.Err() != nil {
		return nil, 0, 0
	}
	return arr, len(arr[0]), len(arr)
}
