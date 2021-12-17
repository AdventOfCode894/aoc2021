package aoc2021d16

type bufPool struct {
	opValues          [][]ExpressionValue
	limitedBitReaders []*limitedBitReader
}

func (bp *bufPool) MakeOpValueBuffer() []ExpressionValue {
	if len(bp.opValues) < 1 {
		bp.opValues = append(bp.opValues, make([]ExpressionValue, 0, 100))
	}
	b := bp.opValues[len(bp.opValues)-1]
	bp.opValues = bp.opValues[:len(bp.opValues)-1]
	return b[:0]
}

func (bp *bufPool) ReturnOpValueBuffer(b []ExpressionValue) {
	bp.opValues = append(bp.opValues, b)
}

func (bp *bufPool) MakeLimitedBitReader() *limitedBitReader {
	if len(bp.limitedBitReaders) < 1 {
		bp.limitedBitReaders = append(bp.limitedBitReaders, new(limitedBitReader))
	}
	lr := bp.limitedBitReaders[len(bp.limitedBitReaders)-1]
	bp.limitedBitReaders = bp.limitedBitReaders[:len(bp.limitedBitReaders)-1]
	return lr
}

func (bp *bufPool) ReturnLimitedBitReader(lr *limitedBitReader) {
	bp.limitedBitReaders = append(bp.limitedBitReaders, lr)
}
