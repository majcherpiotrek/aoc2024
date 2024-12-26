package day_21

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

type Pad struct {
	ButtonCoords map[byte][]int
	Gap          []int
	Pad          [][]byte
}

type Button struct {
	Text   string
	Coords []int
}

func resetCursor() {
	fmt.Print("\033[2K\r")
}

func clearConsole() {
	fmt.Print("\033[H\033[2J")
}

func (pad *Pad) print(activeKey byte) {
	height := len(pad.Pad)
	width := len(pad.Pad[0])

	for y := 0; y < height; y++ {
		fmt.Printf("%s+\n", strings.Repeat("+---", width))
		resetCursor()

		row := ""
		for x := 0; x < width; x++ {
			button := pad.Pad[y][x]

			if button == activeKey {
				row = fmt.Sprintf("%s| %s%c%s ", row, red, button, reset)

			} else {
				row = fmt.Sprintf("%s| %c ", row, button)
			}
		}
		fmt.Printf("%s|\n", row)
		resetCursor()
	}

	fmt.Printf("%s+\n", strings.Repeat("+---", width))
	resetCursor()
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

	if a == b {
		return []byte{}
	}

	diffX := bCoords[0] - aCoords[0]
	diffY := bCoords[1] - aCoords[1]
	diffXAbs := int(math.Abs(float64(diffX)))
	diffYAbs := int(math.Abs(float64(diffY)))

	horizontalMove := ">"
	if diffX < 0 {
		horizontalMove = "<"
	}
	veritcalMove := "^"
	if diffY > 0 {
		veritcalMove = "v"
	}

	horizontalSequence := []byte(strings.Repeat(horizontalMove, diffXAbs))
	verticalSequence := []byte(strings.Repeat(veritcalMove, diffYAbs))

	willHitGapHorizontally := pad.Gap[1] == aCoords[1] &&
		((pad.Gap[0] <= aCoords[0] && pad.Gap[0] >= aCoords[0]+diffX) ||
			(pad.Gap[0] >= aCoords[0] && pad.Gap[0] <= aCoords[0]+diffX))

	willHitGapVertically := pad.Gap[0] == aCoords[0] &&
		((pad.Gap[1] <= aCoords[1] && pad.Gap[1] >= aCoords[1]+diffY) ||
			(pad.Gap[1] >= aCoords[1] && pad.Gap[1] <= aCoords[1]+diffY))

	// Prefer vertical keys first, because it can be 'v', which is more expensive than '>' that is left for horizontal moves. '^' and '>' are equally expensive
	visitFirst := verticalSequence
	visitSecond := horizontalSequence

	// Prefer < first because it's the most expensive key to reach
	if (horizontalMove == "<" && !willHitGapHorizontally) || willHitGapVertically {
		visitFirst = horizontalSequence
		visitSecond = verticalSequence
	}

	return append(visitFirst, visitSecond...)
}

func moveToVector(move byte) []int {
	switch move {
	case '<':
		return []int{-1, 0}
	case '>':
		return []int{1, 0}
	case '^':
		return []int{0, -1}
	case 'v':
		return []int{0, 1}
	default:
		return []int{0, 0}
	}
}

func addVectors(a []int, b []int) []int {
	return []int{a[0] + b[0], a[1] + b[1]}
}

type RobotState struct {
	RobotPad  *Pad
	Pointer   byte
	NextRobot *RobotState
}

func (rs *RobotState) movePointer(move byte) {
	nextPointer, err := rs.RobotPad.nextKey(rs.Pointer, move)
	resetCursor()
	if err == nil {
		rs.Pointer = nextPointer
	}
}

func (rs *RobotState) pressButton() {
	if rs.NextRobot != nil {
		if rs.Pointer == 'A' {
			rs.NextRobot.pressButton()
		} else {
			rs.NextRobot.movePointer(rs.Pointer)
		}
	} else {
		fmt.Printf("%c\n", rs.Pointer)
	}
}

const red = "\033[31m"
const green = "\033[32m"
const reset = "\033[0m"

func (pad *Pad) nextKey(currentPointer byte, move byte) (byte, error) {
	currentCoords, hasButton := pad.ButtonCoords[currentPointer]
	if !hasButton {
		return 0, fmt.Errorf("Out of bounds")
	}
	newCoords := addVectors(currentCoords, moveToVector(move))
	// fmt.Printf("current: %c, move: %c, coords: %v, vector: %v, new coords: %v\n", currentPointer, move, currentCoords, moveToVector(move), newCoords)
	resetCursor()

	for button, coords := range pad.ButtonCoords {
		if coords[0] == newCoords[0] && coords[1] == newCoords[1] {
			return button, nil
		}
	}

	return 0, fmt.Errorf("Invalid coords")
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
	Pad: [][]byte{
		{'7', '8', '9'},
		{'4', '5', '6'},
		{'1', '2', '3'},
		{' ', '0', 'A'},
	},
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
	Pad: [][]byte{
		{' ', '^', 'A'},
		{'<', 'v', '>'},
	},
}

func translateNumericToDirectional(sequence []byte) []byte {
	path := make([]byte, 0)
	current := 'A'

	for _, b := range sequence {
		pathToButton := NumericPad.shortestPath(byte(current), b)
		path = append(path, pathToButton...)
		path = append(path, 'A')
		current = rune(b)
	}

	return path
}

func translateDirectionalToDirectional(sequence []byte) []byte {
	path := make([]byte, 0)
	current := 'A'

	for _, b := range sequence {
		pathToButton := DirectionalPad.shortestPath(byte(current), b)
		path = append(path, pathToButton...)
		path = append(path, 'A')
		current = rune(b)
	}

	return path
}

func runSimulation() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	numPadRobot := RobotState{
		NextRobot: nil,
		Pointer:   'A',
		RobotPad:  &NumericPad,
	}

	firstDirPadRobot := RobotState{
		NextRobot: &numPadRobot,
		Pointer:   'A',
		RobotPad:  &DirectionalPad,
	}

	secondDirPadRobot := RobotState{
		NextRobot: &firstDirPadRobot,
		Pointer:   'A',
		RobotPad:  &DirectionalPad,
	}

	// Read a single byte
	buf := make([]byte, 3) // Arrow keys send a sequence of 3 bytes

	for {

		clearConsole()
		resetCursor()
		secondDirPadRobot.RobotPad.print(secondDirPadRobot.Pointer)
		resetCursor()
		fmt.Println("")
		resetCursor()
		firstDirPadRobot.RobotPad.print(firstDirPadRobot.Pointer)
		resetCursor()
		fmt.Println("")
		resetCursor()
		numPadRobot.RobotPad.print(numPadRobot.Pointer)
		resetCursor()
		fmt.Println("")
		resetCursor()

		n, err := os.Stdin.Read(buf)
		if err != nil {
			panic(err)
		}

		// Interpret the byte sequence
		if n == 1 && buf[0] == 27 { // ESC key
			resetCursor()
			fmt.Println("ESC pressed. Exiting...")
			resetCursor()
			break
		} else if n == 3 && buf[0] == 27 && buf[1] == 91 { // Arrow keys
			switch buf[2] {
			case 65:
				// Up arrow pressed
				secondDirPadRobot.movePointer('^')
			case 66:
				// Down arrow pressed
				secondDirPadRobot.movePointer('v')
			case 67:
				// Right arrow pressed
				secondDirPadRobot.movePointer('>')
			case 68:
				// Left arrow pressed
				secondDirPadRobot.movePointer('<')
			}
		} else if n == 1 && (buf[0] == 'A' || buf[0] == 'a') { // A key
			secondDirPadRobot.pressButton()
		}
		resetCursor()
		fmt.Println("-----------------------")
		resetCursor()
	}

}
func Part1(input *[]string) (int, error) {
	// runSimulation()
	// return 0, nil
	sumOfComplexities := 0

	NumericPad.print('A')
	DirectionalPad.print('A')

	for _, sequence := range *input {
		numPart := sequence[0 : len(sequence)-1]
		num, err := strconv.Atoi(numPart)
		if err != nil {
			return -1, err
		}

		fmt.Printf("Sequence: %s\n", sequence)

		level1Seq := translateNumericToDirectional([]byte(sequence))
		level2Seq := translateDirectionalToDirectional(level1Seq)
		level3Seq := translateDirectionalToDirectional(level2Seq)

		fmt.Printf("%s: %s\n", sequence, string(level3Seq))
		fmt.Printf("%s: %s\n", sequence, string(level2Seq))
		fmt.Printf("%s: %s\n", sequence, string(level1Seq))

		complexity := num * len(level3Seq)
		fmt.Printf("%s: %d * %d = %d\n", sequence, num, len(level3Seq), complexity)
		sumOfComplexities += complexity
	}

	return sumOfComplexities, nil
}

func Part2(input *[]string) (int, error) {

	return -1, fmt.Errorf("not implemented")
}
