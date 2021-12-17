/*

Reader implementation.

*/

package bitio

import (
	"io"
)

type ReaderAndByteReader interface {
	io.Reader
	io.ByteReader
}

// Reader is the bit reader implementation.
//
// If you need the number of processed bits, use CountReader.
//
// For convenience, it also implements io.Reader and io.ByteReader.
type Reader struct {
	in    ReaderAndByteReader
	cache byte // unread bits are stored here
	bits  byte // number of unread bits in cache

	// tryError holds the first error occurred in TryXXX() methods.
	tryError error
}

func (r *Reader) Init(in ReaderAndByteReader) {
	r.in = in
	r.cache = 0
	r.bits = 0
	r.tryError = nil
}

// readBits reads n bits and returns them as the lowest n bits of u.
func (r *Reader) readBits(n uint8) (u uint64, err error) {
	// Some optimization, frequent cases
	if n < r.bits {
		// cache has all needed bits, and there are some extra which will be left in cache
		shift := r.bits - n
		u = uint64(r.cache >> shift)
		r.cache &= 1<<shift - 1
		r.bits = shift
		return
	}

	if n > r.bits {
		// all cache bits needed, and it's not even enough so more will be read
		if r.bits > 0 {
			u = uint64(r.cache)
			n -= r.bits
		}
		// Read whole bytes
		for n >= 8 {
			b, err2 := r.in.ReadByte()
			if err2 != nil {
				return 0, err2
			}
			u = u<<8 + uint64(b)
			n -= 8
		}
		// Read last fraction, if any
		if n > 0 {
			if r.cache, err = r.in.ReadByte(); err != nil {
				return 0, err
			}
			shift := 8 - n
			u = u<<n + uint64(r.cache>>shift)
			r.cache &= 1<<shift - 1
			r.bits = shift
		} else {
			r.bits = 0
		}
		return u, nil
	}

	// cache has exactly as many as needed
	r.bits = 0 // no need to clear cache, will be overwritten on next read
	return uint64(r.cache), nil
}

// ReadBits tries to read n bits.
//
// If there was a previous Err, it does nothing. Else it calls readBits(),
// returns the data it provides and stores the error for return by Err().
func (r *Reader) ReadBits(n uint8) (u uint64) {
	if r.tryError == nil {
		u, r.tryError = r.readBits(n)
	}
	return
}

func (r *Reader) Err() error {
	return r.tryError
}
