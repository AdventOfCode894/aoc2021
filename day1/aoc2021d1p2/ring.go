package main

type measurementRing struct {
	measurements []uint
	next         int
}

func newMeasurementRing(capacity uint) *measurementRing {
	return &measurementRing{
		measurements: make([]uint, 0, capacity),
	}
}

func (mr *measurementRing) IsFull() bool {
	return len(mr.measurements) >= cap(mr.measurements)
}

func (mr *measurementRing) Record(depth uint) {
	if !mr.IsFull() {
		mr.measurements = append(mr.measurements, depth)
		return
	}
	mr.measurements[mr.next] = depth
	mr.next++
	if mr.next >= len(mr.measurements) {
		mr.next = 0
	}
}

func (mr *measurementRing) Sum() uint {
	sum := uint(0)
	for _, depth := range mr.measurements {
		sum += depth
	}
	return sum
}
