package day_15

import (
	"fmt"
	"strconv"
	"strings"
)

type WarehouseMap struct {
	RobotPosition []int
	Warehouse     [][]byte
}

type Direction byte

const (
	Up    = '^'
	Right = '>'
	Down  = 'v'
	Left  = '<'
)

type WarehouseElement byte

const (
	Wall     = '#'
	Space    = '.'
	Box      = 'O'
	Robot    = '@'
	BoxLeft  = '['
	BoxRight = ']'
)

func NewWarehouseMap(warehouse *[]string) (WarehouseMap, error) {
	robotPosition := make([]int, 0, 2)
	warehouseBytes := make([][]byte, 0, len(*warehouse))

	for y := 0; y < len(*warehouse); y++ {
		rowBytes := make([]byte, 0, len((*warehouse)[y]))

		for x := 0; x < len((*warehouse)[y]); x++ {
			element := (*warehouse)[y][x]
			if element == Robot {
				robotPosition = []int{x, y}
			}

			rowBytes = append(rowBytes, element)
		}

		warehouseBytes = append(warehouseBytes, rowBytes)
	}

	if len(robotPosition) != 2 {
		return WarehouseMap{}, fmt.Errorf("Failed to determine the robot's position")
	}

	return WarehouseMap{
		RobotPosition: robotPosition,
		Warehouse:     warehouseBytes,
	}, nil
}

func (w *WarehouseMap) Print() {
	for _, rowBytes := range w.Warehouse {
		rowStr := string(rowBytes)
		fmt.Println(rowStr)
	}
	fmt.Printf("\n")
}

func (w *WarehouseMap) MoveRobot(dir Direction) {
	updatedPosition, err := moveRecursive(w.RobotPosition, [][]int{w.RobotPosition}, dir, &w.Warehouse)
	if err != nil {
		return
	}
	w.RobotPosition = updatedPosition
}

func moveVectorForDirection(dir Direction) []int {
	switch dir {
	case Up:
		return []int{0, -1}
	case Right:
		return []int{1, 0}
	case Down:
		return []int{0, 1}
	case Left:
		return []int{-1, 0}
	default:
		return []int{0, 0}
	}
}

func moveRecursive(robotLocation []int, block [][]int, direction Direction, warehouse *[][]byte) ([]int, error) {
	blockDestination := make([][]int, 0, len(block))
	moveVector := moveVectorForDirection(direction)

	currentBlockElementsCache := make(map[string]byte)
	for _, b := range block {
		x := b[0]
		y := b[1]

		if x < 0 || x >= len((*warehouse)[0]) || y < 0 || y >= len(*warehouse) {
			return []int{}, fmt.Errorf("coords out of bounds")
		}

		destinationCoords := addVectors(b, moveVector)
		nextX := destinationCoords[0]
		nextY := destinationCoords[1]

		if nextX < 0 || nextX >= len((*warehouse)[0]) || nextY < 0 || nextY >= len(*warehouse) {
			return []int{}, fmt.Errorf("next coords out of bounds")

		}

		key := fmt.Sprintf("%d-%d", b[0], b[1])
		currentBlockElementsCache[key] = (*warehouse)[b[1]][b[0]]

		blockDestination = append(blockDestination, []int{nextX, nextY})
	}

	onlySpacesAhead := true
	blockForNextStep := block

	for _, bDest := range blockDestination {
		key := fmt.Sprintf("%d-%d", bDest[0], bDest[1])
		_, partOfCurrentBlock := currentBlockElementsCache[key]
		if partOfCurrentBlock {
			continue
		}

		destElement := (*warehouse)[bDest[1]][bDest[0]]

		switch destElement {
		case Wall:
			return []int{}, fmt.Errorf("Block hit the wall")
		case Space:
			continue
		case Box:
			onlySpacesAhead = false
			blockForNextStep = append(blockForNextStep, bDest)
			continue
		case BoxLeft:
			onlySpacesAhead = false
			blockForNextStep = append(blockForNextStep, bDest)
			blockForNextStep = append(blockForNextStep, []int{bDest[0] + 1, bDest[1]})
			continue
		case BoxRight:
			onlySpacesAhead = false
			blockForNextStep = append(blockForNextStep, bDest)
			blockForNextStep = append(blockForNextStep, []int{bDest[0] - 1, bDest[1]})
			continue
		default:
			return []int{}, fmt.Errorf("Unknown destination element")
		}
	}

	if onlySpacesAhead {
		for _, b := range block {
			dest := addVectors(b, moveVector)
			key := fmt.Sprintf("%d-%d", b[0], b[1])
			element, isInCache := currentBlockElementsCache[key]

			if !isInCache {
				panic("WTF")
			}

			(*warehouse)[dest[1]][dest[0]] = element
		}

		for _, b := range blockDestination {
			key := fmt.Sprintf("%d-%d", b[0], b[1])
			_, isInCache := currentBlockElementsCache[key]

			if isInCache {
				delete(currentBlockElementsCache, key)
			}
		}

		for key := range currentBlockElementsCache {
			split := strings.Split(key, "-")
			x, err := strconv.Atoi(split[0])
			if err != nil {
				panic("WTF")
			}
			y, err := strconv.Atoi(split[1])
			if err != nil {
				panic("WTF")
			}

			(*warehouse)[y][x] = Space
		}

		return addVectors(robotLocation, moveVector), nil
	} else {
		return moveRecursive(robotLocation, blockForNextStep, direction, warehouse)
	}

}

func addVectors(a []int, b []int) []int {
	return []int{a[0] + b[0], a[1] + b[1]}
}

type PuzzleData struct {
	WarehouseMap WarehouseMap
	Moves        []Direction
}

func parseInput(input *[]string, part int) (PuzzleData, error) {
	var wareHouseMapRows []string
	moves := ""

	readingWarehouseMap := true

	for _, row := range *input {
		if len(row) == 0 {
			readingWarehouseMap = false
			continue
		}

		if readingWarehouseMap {
			rowToAdd := row

			if part == 2 {
				extendedRowBytes := make([]byte, 0, len(row)*2)
				rowAsBytes := []byte(row)
				for _, b := range rowAsBytes {
					switch b {
					case Robot:
						extendedRowBytes = append(extendedRowBytes, Robot, Space)
						continue
					case Space:
						extendedRowBytes = append(extendedRowBytes, Space, Space)
						continue
					case Box:
						extendedRowBytes = append(extendedRowBytes, BoxLeft, BoxRight)
						continue
					case Wall:
						extendedRowBytes = append(extendedRowBytes, Wall, Wall)
						continue
					default:
						return PuzzleData{}, fmt.Errorf("Failed to extend symbol %c", b)
					}
				}
				rowToAdd = string(extendedRowBytes)
			}

			wareHouseMapRows = append(wareHouseMapRows, rowToAdd)
		} else {
			moves = fmt.Sprintf("%s%s", moves, row)
		}
	}

	warehouseMap, err := NewWarehouseMap(&wareHouseMapRows)
	if err != nil {
		return PuzzleData{}, err
	}

	moveBytes := []Direction(moves)

	return PuzzleData{
		WarehouseMap: warehouseMap,
		Moves:        moveBytes,
	}, nil
}

func Part1(rows *[]string) (int, error) {
	data, err := parseInput(rows, 1)

	if err != nil {
		return -1, err
	}

	data.WarehouseMap.Print()

	for _, move := range data.Moves {
		fmt.Printf("Move %c\n", move)
		data.WarehouseMap.MoveRobot(move)
		data.WarehouseMap.Print()
	}

	sum := 0

	for y, row := range data.WarehouseMap.Warehouse {
		for x, element := range row {
			if element == Box {
				coords := 100*y + x
				sum += coords
			}
		}
	}

	return sum, nil

}

func Part2(rows *[]string) (int, error) {
	data, err := parseInput(rows, 2)

	if err != nil {
		return -1, err
	}

	data.WarehouseMap.Print()

	for _, move := range data.Moves {
		fmt.Printf("Move %c\n", move)
		data.WarehouseMap.MoveRobot(move)
		data.WarehouseMap.Print()
	}

	sum := 0

	for y, row := range data.WarehouseMap.Warehouse {
		for x, element := range row {
			if element == BoxLeft {
				fromLeft := x
				fromRight := len(row) - (x + 1)
				coords := 100*y + fromLeft

				if fromRight < fromLeft {
					coords = 100*y + fromLeft
				}

				sum += coords
			}
		}
	}

	return sum, nil

}
