package day_21

import (
	"fmt"
	"strings"
)

type Pad struct {
	ButtonCoords map[byte][]int
	Gap          []int
}

func (pad *Pad) shortestPath(a byte, b byte) []byte {
	aCoords, hasA := pad.ButtonCoords[a]
	if !hasA {
		return []byte{}
	}
	bCoords, hasB := pad.ButtonCoords[b]
	if !hasB {
		return []byte{}
	}

	diffX := bCoords[0] - aCoords[0]
	diffY := bCoords[1] - aCoords[1]

	var horizontalMove string
	var veritcalMove string

	if diffX > 0 {
		horizontalMove = strings.Repeat(">", diffX)
	} else {
		horizontalMove = strings.Repeat("<", diffX)
	}

	if diffY > 0 {
		veritcalMove = strings.Repeat("v", diffY)
	} else {
		veritcalMove = strings.Repeat("^", diffY)
	}

	if aCoords[0]+diffX != pad.Gap[0] || aCoords[1] != pad.Gap[1] {
		// Safe to move horizontally first
		moveStr := fmt.Sprintf("%s%s", horizontalMove, veritcalMove)
		return []byte(moveStr)

	} else {
		// Safe to move vertically first
		moveStr := fmt.Sprintf("%s%s", veritcalMove, horizontalMove)
		return []byte(moveStr)
	}
}

var NumericPad = Pad{
	ButtonCoords: map[byte][]int{
		'7': {0, 0},
		'8': {1, 0},
		'9': {2, 0},
		'4': {0, 1},
		'5': {1, 1},
		'6': {2, 1},
		'1': {0, 2},
		'2': {1, 2},
		'3': {2, 2},
		'0': {1, 3},
		'A': {2, 3},
	},
	Gap: []int{0, 3},
}

var DirectionalPad = Pad{
	ButtonCoords: map[byte][]int{
		'^': {1, 0},
		'<': {0, 1},
		'v': {1, 1},
		'>': {2, 1},
		'A': {2, 0},
	},
	Gap: []int{0, 0},
}

func Part1(input *[]string) (int, error) {

	return -1, fmt.Errorf("not implemented")
}

func Part2(input *[]string) (int, error) {

	return -1, fmt.Errorf("not implemented")
}
