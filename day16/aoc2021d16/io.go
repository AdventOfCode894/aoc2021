package aoc2021d16

import (
	"errors"

	"github.com/icza/bitio"
)

type bitReader interface {
	ReadBits(n uint8) uint64
	Err() error
}

type bitioReader bitio.Reader

func (r *bitioReader) ReadBits(n uint8) uint64 {
	return (*bitio.Reader)(r).TryReadBits(n)
}

func (r *bitioReader) Err() error {
	return (*bitio.Reader)(r).TryError
}

type limitedBitReader struct {
	R             bitReader
	RemainingBits uint64

	overflowed bool
}

var errBitOverflow = errors.New("read more bits than were available")

func (r *limitedBitReader) ReadBits(n uint8) uint64 {
	if uint64(n) > r.RemainingBits {
		n = uint8(r.RemainingBits)
		r.overflowed = true
	}
	r.RemainingBits -= uint64(n)
	return r.R.ReadBits(n)
}

func (r *limitedBitReader) Err() error {
	if r.overflowed {
		return errBitOverflow
	}
	return r.R.Err()
}
