package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"

	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
	"github.com/icza/bitio"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

var totalVersionNumber uint64

func solvePuzzle(r io.Reader) (int, error) {
	br := bitio.NewReader(hex.NewDecoder(r))

	totalVersionNumber = 0
	if err := parsePacketStream(br); err != nil {
		return 0, err
	}

	return int(totalVersionNumber), nil
}

type packetTypeID uint64

const (
	literalValueType packetTypeID = 4
)

func parsePacketStream(br *bitio.Reader) error {
	for {
		parsed, err := parseOnePacket(br)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		if !parsed {
			break
		}
	}
	return nil
}

func parseOnePacket(br *bitio.Reader) (bool, error) {
	version := br.TryReadBits(3)
	packetType := packetTypeID(br.TryReadBits(3))
	if br.TryError != nil {
		return false, br.TryError
	}

	// For part one
	totalVersionNumber += version

	switch packetType {
	case literalValueType:
		_, err := parseLiteralValue(br)
		if err != nil {
			return false, err
		}
	default:
		if err := parseOperator(br); err != nil {
			return false, err
		}
	}
	return true, nil
}

func parseLiteralValue(br *bitio.Reader) (uint64, error) {
	answer := uint64(0)
	readBits := 0
	for {
		more := br.TryReadBits(1)
		subValue := br.TryReadBits(4)
		readBits += 4
		if readBits > 64 {
			return 0, errors.New("literal value overflow")
		}
		answer <<= 4
		answer |= subValue
		if more < 1 {
			break
		}
	}
	return answer, nil
}

func parseOperator(br *bitio.Reader) error {
	lengthTypeID := br.TryReadBits(1)
	if lengthTypeID == 0 {
		subBitLength := br.TryReadBits(15)
		var subBuffer bytes.Buffer
		bw := bitio.NewWriter(&subBuffer)
		const copyChunkSize = 64
		for i := uint64(0); i < subBitLength; i += copyChunkSize {
			toRead := subBitLength - i
			if toRead > copyChunkSize {
				toRead = copyChunkSize
			}
			chunk := br.TryReadBits(uint8(toRead))
			_ = bw.WriteBits(chunk, uint8(toRead))
		}
		if err := bw.Close(); err != nil {
			return err
		}
		if err := parsePacketStream(bitio.NewReader(&subBuffer)); err != nil {
			return err
		}
	} else {
		subCount := br.TryReadBits(11)
		for i := uint64(0); i < subCount; i++ {
			parsed, err := parseOnePacket(br)
			if err != nil {
				return err
			}
			if !parsed {
				return errors.New("truncated sub-stream")
			}
		}
	}
	return nil
}
