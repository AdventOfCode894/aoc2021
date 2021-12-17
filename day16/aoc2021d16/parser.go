package aoc2021d16

import (
	"errors"
	"io"

	"github.com/icza/bitio"
)

type ExpressionValue uint64

const ExpressionValueBits = 64

var ErrExpressionOverflow = errors.New("expression value overflow")

func EvaluateExpression(r io.Reader) (value ExpressionValue, versionSum uint64, err error) {
	br := (*bitioReader)(bitio.NewReader(r))

	value, versionSum, err = parseOnePacket(br)
	if err != nil {
		return 0, 0, err
	}

	return value, versionSum, err
}

type packetTypeID uint64

const (
	sumType packetTypeID = iota
	productType
	minimumType
	maximumType
	literalValueType
	greaterThan
	lessThan
	equalTo
)

const minPacketLength = 3 + 3 + 5

func parseOnePacket(br bitReader) (ExpressionValue, uint64, error) {
	versionSum := br.ReadBits(3)
	packetType := packetTypeID(br.ReadBits(3))
	if err := br.Err(); err != nil {
		return 0, 0, err
	}

	if packetType == literalValueType {
		v, err := parseLiteralValue(br)
		if err != nil {
			return 0, 0, err
		}
		return v, versionSum, nil
	} else {
		v, vs, err := parseOperator(br, packetType)
		if err != nil {
			return 0, 0, err
		}
		versionSum += vs
		return v, versionSum, nil
	}
}

func parseLiteralValue(br bitReader) (ExpressionValue, error) {
	v := ExpressionValue(0)
	readBits := 0
	for {
		more := br.ReadBits(1)
		subValue := ExpressionValue(br.ReadBits(4))
		readBits += 4
		if readBits > ExpressionValueBits {
			return 0, ErrExpressionOverflow
		}
		v <<= 4
		v |= subValue
		if more < 1 {
			break
		}
	}
	return v, nil
}

func parseOperator(br bitReader, packetType packetTypeID) (ExpressionValue, uint64, error) {
	lengthTypeID := br.ReadBits(1)

	const estimatedPacketCount = 10
	subValues := make([]ExpressionValue, 0, estimatedPacketCount)
	versionSum := uint64(0)
	if lengthTypeID == 0 {
		subBitLength := br.ReadBits(15)
		lr := &limitedBitReader{
			R:             br,
			RemainingBits: subBitLength,
		}
		for lr.RemainingBits > 0 {
			v, vs, err := parseOnePacket(lr)
			if err != nil {
				return 0, 0, err
			}
			subValues = append(subValues, v)
			versionSum += vs
		}
	} else {
		subCount := br.ReadBits(11)
		for i := uint64(0); i < subCount; i++ {
			v, vs, err := parseOnePacket(br)
			if err != nil {
				return 0, 0, err
			}
			subValues = append(subValues, v)
			versionSum += vs
		}
	}
	if len(subValues) < 1 {
		return 0, 0, errors.New("operation contained no sub-values")
	}

	v, err := evaluateExpression(packetType, subValues)
	if err != nil {
		return 0, 0, err
	}
	return v, versionSum, nil
}

func evaluateExpression(packetType packetTypeID, subValues []ExpressionValue) (ExpressionValue, error) {
	var v ExpressionValue
	switch packetType {
	case sumType:
		v = subValues[0]
		for i := 1; i < len(subValues); i++ {
			newSum := v + subValues[i]
			if newSum < v {
				return 0, ErrExpressionOverflow
			}
			v = newSum
		}
	case productType:
		v = subValues[0]
		for i := 1; i < len(subValues); i++ {
			newProduct := v * subValues[i]
			if newProduct < v && v != 0 && subValues[i] != 0 {
				return 0, ErrExpressionOverflow
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
