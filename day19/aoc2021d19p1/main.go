package main

import (
	"fmt"
	"io"

	"github.com/AdventOfCode894/aoc2021/internal/aocio"
	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
)

const scannerRange = 1000
const minScannerOverlap = 12

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

type scannerRegion struct {
	relativeBeacons map[point3D]struct{}
}

type point3D struct {
	X int
	Y int
	Z int
}

type rotation3D struct {
	m [9]int
}

type scannerOverlap struct {
	relativePos       point3D
	secondOrientation rotation3D
}

type scannerTransform struct {
	globalPos         point3D
	globalOrientation rotation3D
}

var scannerRotations = []rotation3D{
	{m: [...]int{1, 0, 0, 0, 1, 0, 0, 0, 1}},
	{m: [...]int{1, 0, 0, 0, 0, -1, 0, 1, 0}},   // x90
	{m: [...]int{1, 0, 0, 0, -1, 0, 0, 0, -1}},  // x180
	{m: [...]int{1, 0, 0, 0, 0, 1, 0, -1, 0}},   // x270
	{m: [...]int{-1, 0, 0, 0, 1, 0, 0, 0, -1}},  // y180
	{m: [...]int{-1, 0, 0, 0, 0, -1, 0, -1, 0}}, // y180 x90
	{m: [...]int{-1, 0, 0, 0, -1, 0, 0, 0, 1}},  // y180 x180
	{m: [...]int{-1, 0, 0, 0, 0, 1, 0, 1, 0}},   // y180 x270
	{m: [...]int{0, 0, 1, 0, 1, 0, -1, 0, 0}},   // y90
	{m: [...]int{0, 0, 1, 1, 0, 0, 0, 1, 0}},    // y90 z90
	{m: [...]int{0, 0, 1, 0, -1, 0, 1, 0, 0}},   // y90 z180
	{m: [...]int{0, 0, 1, -1, 0, 0, 0, -1, 0}},  // y90 z270
	{m: [...]int{0, 0, -1, 0, 1, 0, 1, 0, 0}},   // y270
	{m: [...]int{0, 0, -1, 1, 0, 0, 0, -1, 0}},  // y270 z90
	{m: [...]int{0, 0, -1, 0, -1, 0, -1, 0, 0}}, // y270 z180
	{m: [...]int{0, 0, -1, -1, 0, 0, 0, 1, 0}},  // y270 z270
	{m: [...]int{0, -1, 0, 1, 0, 0, 0, 0, 1}},   // z90
	{m: [...]int{0, -1, 0, 0, 0, 1, -1, 0, 0}},  // z90 y90
	{m: [...]int{0, -1, 0, -1, 0, 0, 0, 0, -1}}, // z90 y180
	{m: [...]int{0, -1, 0, 0, 0, -1, 1, 0, 0}},  // z90 y270
	{m: [...]int{0, 1, 0, -1, 0, 0, 0, 0, 1}},   // z270
	{m: [...]int{0, 1, 0, 0, 0, -1, -1, 0, 0}},  // z270 y90
	{m: [...]int{0, 1, 0, 1, 0, 0, 0, 0, -1}},   // z270 y180
	{m: [...]int{0, 1, 0, 0, 0, 1, 1, 0, 0}},    // z270 y270
}

func (p *point3D) Rotate(r rotation3D, src point3D) {
	p.X = r.m[0]*src.X + r.m[1]*src.Y + r.m[2]*src.Z
	p.Y = r.m[3]*src.X + r.m[4]*src.Y + r.m[5]*src.Z
	p.Z = r.m[6]*src.X + r.m[7]*src.Y + r.m[8]*src.Z
}

func (p *point3D) Add(p1 point3D, p2 point3D) {
	p.X = p1.X + p2.X
	p.Y = p1.Y + p2.Y
	p.Z = p1.Z + p2.Z
}

func (p *point3D) Sub(subFrom point3D, subAmount point3D) {
	p.X = subFrom.X - subAmount.X
	p.Y = subFrom.Y - subAmount.Y
	p.Z = subFrom.Z - subAmount.Z
}

func (p *point3D) Negate(other point3D) {
	p.X = -other.X
	p.Y = -other.Y
	p.Z = -other.Z
}

func (r *rotation3D) Multiply(first rotation3D, second rotation3D) {
	r.m[0] = first.m[0]*second.m[0] + first.m[1]*second.m[3] + first.m[2]*second.m[6]
	r.m[1] = first.m[0]*second.m[1] + first.m[1]*second.m[4] + first.m[2]*second.m[7]
	r.m[2] = first.m[0]*second.m[2] + first.m[1]*second.m[5] + first.m[2]*second.m[8]
	r.m[3] = first.m[3]*second.m[0] + first.m[4]*second.m[3] + first.m[5]*second.m[6]
	r.m[4] = first.m[3]*second.m[1] + first.m[4]*second.m[4] + first.m[5]*second.m[7]
	r.m[5] = first.m[3]*second.m[2] + first.m[4]*second.m[5] + first.m[5]*second.m[8]
	r.m[6] = first.m[6]*second.m[0] + first.m[7]*second.m[3] + first.m[8]*second.m[6]
	r.m[7] = first.m[6]*second.m[1] + first.m[7]*second.m[4] + first.m[8]*second.m[7]
	r.m[8] = first.m[6]*second.m[2] + first.m[7]*second.m[5] + first.m[8]*second.m[8]
}

func solvePuzzle(r io.Reader) (int, error) {
	pr := aocio.NewPuzzleReader(r)
	var scanners []*scannerRegion
	for pr.NextNonEmptyLine() {
		tr := pr.LineTokenReader()
		tr.ConsumeString("--- scanner ")
		scannerNum, _ := tr.NextUint(' ', 10)
		tr.ConsumeString("---")
		tr.ConsumeEOL()
		if err := pr.Err(); err != nil {
			return 0, err
		}
		if scannerNum != uint(len(scanners)) {
			return 0, fmt.Errorf("unexpected scanner number: %d", scannerNum)
		}

		relativeBeacons := make(map[point3D]struct{})
		for pr.NextLine() {
			if pr.IsLineEmpty() {
				break
			}
			tr = pr.LineTokenReader()
			x, _ := tr.NextInt(',', 10)
			y, _ := tr.NextInt(',', 10)
			z, _ := tr.NextInt(aocio.EOLDelim, 10)
			relativeBeacons[point3D{X: x, Y: y, Z: z}] = struct{}{}

			if _, weCanSee := relativeBeacons[point3D{X: x, Y: y, Z: z}]; !weCanSee {
				continue
			}
		}
		scanners = append(scanners, &scannerRegion{relativeBeacons: relativeBeacons})
		if err := pr.Err(); err != nil {
			return 0, err
		}
	}

	overlapGraph := make(map[int]map[int]*scannerOverlap)
	for i, s1 := range scanners {
		for j, s2 := range scanners {
			if j <= i {
				continue
			}
			if overlap, overlapping := s1.CheckOverlap(s2); overlapping {
				if _, ok := overlapGraph[i]; !ok {
					overlapGraph[i] = make(map[int]*scannerOverlap)
				}
				overlapGraph[i][j] = overlap
				if _, ok := overlapGraph[j]; !ok {
					overlapGraph[j] = make(map[int]*scannerOverlap)
				}
				overlapGraph[j][i], _ = s2.CheckOverlap(s1)
			}
		}
	}

	scannerTransforms := make([]scannerTransform, len(scanners))
	seenScanners := make([]bool, len(scanners))
	mapScanners(overlapGraph, scannerTransforms, seenScanners, 0, point3D{}, scannerRotations[0])

	globalBeacons := make(map[point3D]struct{})
	for i, region := range scanners {
		t := scannerTransforms[i]
		for p := range region.relativeBeacons {
			var oriented, globalPos point3D
			oriented.Rotate(t.globalOrientation, p)
			globalPos.Add(t.globalPos, oriented)
			globalBeacons[globalPos] = struct{}{}
		}
	}

	return len(globalBeacons), nil
}

func mapScanners(overlapGraph map[int]map[int]*scannerOverlap, scannerTransforms []scannerTransform, seenScanners []bool, atScanner int, atPosition point3D, orientation rotation3D) {
	if seenScanners[atScanner] {
		return
	}
	seenScanners[atScanner] = true
	scannerTransforms[atScanner].globalPos = atPosition
	scannerTransforms[atScanner].globalOrientation = orientation

	for to, overlap := range overlapGraph[atScanner] {
		var offset, newPos point3D
		offset.Rotate(orientation, overlap.relativePos)
		newPos.Add(atPosition, offset)

		var newOrientation rotation3D
		newOrientation.Multiply(orientation, overlap.secondOrientation)

		mapScanners(overlapGraph, scannerTransforms, seenScanners, to, newPos, newOrientation)
	}
}

func (sr *scannerRegion) CheckOverlap(other *scannerRegion) (*scannerOverlap, bool) {
	for _, otherRotation := range scannerRotations {
		for otherBeacon := range other.relativeBeacons {
			var otherBeaconPos point3D
			otherBeaconPos.Rotate(otherRotation, otherBeacon)
			for ourMatchingBeacon := range sr.relativeBeacons {
				var offset point3D
				offset.Sub(otherBeaconPos, ourMatchingBeacon)
				overlapping := 1
				for otherOtherBeacon := range other.relativeBeacons {
					if otherOtherBeacon == otherBeacon {
						continue
					}
					var otherOtherBeaconPos point3D
					otherOtherBeaconPos.Rotate(otherRotation, otherOtherBeacon)
					var otherOtherRelativePos point3D
					otherOtherRelativePos.Sub(otherOtherBeaconPos, offset)
					if !sr.beaconWithinRange(otherOtherRelativePos) {
						continue
					}
					if _, weCanSee := sr.relativeBeacons[otherOtherRelativePos]; !weCanSee {
						continue
					}
					overlapping++
				}
				if overlapping >= minScannerOverlap {
					offset.Negate(offset)
					return &scannerOverlap{
						relativePos:       offset,
						secondOrientation: otherRotation,
					}, true
				}
			}
		}
	}
	return nil, false
}

func (sr *scannerRegion) beaconWithinRange(p point3D) bool {
	return p.X >= -scannerRange && p.X <= scannerRange && p.Y >= -scannerRange && p.Y <= scannerRange && p.Z >= -scannerRange && p.Z <= scannerRange
}
