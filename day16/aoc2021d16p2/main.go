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

func solvePuzzle(r io.Reader) (int, error) {
	br := bitio.NewCountReader(hex.NewDecoder(r))

	answer, err := parseOnePacket(br)
	if err != nil {
		return 0, err
	}

	return int(answer), nil
}

type packetTypeID uint64

type expressionValue uint64

const expressionMaxBits = 64

const (
	sumType          packetTypeID = 0
	productType      packetTypeID = 1
	minimumType      packetTypeID = 2
	maximumType      packetTypeID = 3
	literalValueType packetTypeID = 4
	greaterThan      packetTypeID = 5
	lessThan         packetTypeID = 6
	equalTo          packetTypeID = 7
)

func parseOnePacket(br *bitio.CountReader) (expressionValue, error) {
	version := br.TryReadBits(3)
	version = version
	packetType := packetTypeID(br.TryReadBits(3))
	if br.TryError != nil {
		return 0, br.TryError
	}

	if packetType == literalValueType {
		v, err := parseLiteralValue(br)
		if err != nil {
			return 0, err
		}
		return v, nil
	} else {
		v, err := parseOperator(br, packetType)
		if err != nil {
			return 0, err
		}
		return v, nil
	}
}

func parseLiteralValue(br *bitio.CountReader) (expressionValue, error) {
	v := expressionValue(0)
	readBits := 0
	for {
		more := br.TryReadBits(1)
		subValue := expressionValue(br.TryReadBits(4))
		readBits += 4
		if readBits > expressionMaxBits {
			return 0, errors.New("literal value overflow")
		}
		v <<= 4
		v |= subValue
		if more < 1 {
			break
		}
	}
	return v, nil
}

func parseOperator(br *bitio.CountReader, packetType packetTypeID) (expressionValue, error) {
	lengthTypeID := br.TryReadBits(1)
	var subValues []expressionValue
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
			return 0, err
		}
		cr := bitio.NewCountReader(&subBuffer)
		for uint64(cr.BitsCount) < subBitLength {
			v, err := parseOnePacket(cr)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return 0, err
			}
			subValues = append(subValues, v)
		}
	} else {
		subCount := br.TryReadBits(11)
		for i := uint64(0); i < subCount; i++ {
			v, err := parseOnePacket(br)
			if err != nil {
				return 0, err
			}
			subValues = append(subValues, v)
		}
	}
	if len(subValues) < 1 {
		return 0, errors.New("no sub-values")
	}

	v, err := evaluateExpression(packetType, subValues)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func evaluateExpression(packetType packetTypeID, subValues []expressionValue) (expressionValue, error) {
	var v expressionValue
	switch packetType {
	case sumType:
		v = subValues[0]
		for i := 1; i < len(subValues); i++ {
			newSum := v + subValues[i]
			if newSum < v {
				return 0, errors.New("sum overflow")
			}
			v = newSum
		}
	case productType:
		v = subValues[0]
		for i := 1; i < len(subValues); i++ {
			newProduct := v * subValues[i]
			if newProduct < v && v != 0 && subValues[i] != 0 {
				return 0, errors.New("product overflow")
			}
			v = newProduct
		}
	case minimumType:
		v = subValues[0]
		for i := 1; i < len(subValues); i++ {
			if subValues[i] < v {
				v = subValues[i]
			}
		}
	case maximumType:
		v = subValues[0]
		for i := 1; i < len(subValues); i++ {
			if subValues[i] > v {
				v = subValues[i]
			}
		}
	case greaterThan:
		if len(subValues) != 2 {
			return 0, errors.New("invalid packet count for greater-than operator")
		}
		if subValues[0] > subValues[1] {
			v = 1
		} else {
			v = 0
		}
	case lessThan:
		if len(subValues) != 2 {
			return 0, errors.New("invalid packet count for less-than operator")
		}
		if subValues[0] < subValues[1] {
			v = 1
		} else {
			v = 0
		}
	case equalTo:
		if len(subValues) != 2 {
			return 0, errors.New("invalid packet count for equal-to operator")
		}
		if subValues[0] == subValues[1] {
			v = 1
		} else {
			v = 0
		}
	}
	return v, nil
}
