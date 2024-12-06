package day_06

import (
	"fmt"
	"strconv"
	"strings"
)

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

func getMoveForDirection(direction Direction) []int {
	switch direction {
	case Up:
		return []int{0, -1}
	case Right:
		return []int{1, 0}
	case Down:
		return []int{0, 1}
	case Left:
		return []int{-1, 0}
	}
	return []int{0, 0}
}

func applyMove(pos []int, move []int) []int {
	return []int{pos[0] + move[0], pos[1] + move[1]}
}

func getGuardStartPosition(rows *[]string) ([]int, error) {
	guardPosition := make([]int, 0, 2)

	for y, row := range *rows {
		guardX := strings.IndexByte(row, '^')

		if guardX != -1 {
			guardPosition = []int{guardX, y}
			break
		}
	}

	if len(guardPosition) == 0 {
		return []int{}, fmt.Errorf("Could not find the guard")
	}

	return guardPosition, nil
}

func getVisitedPositions(rows *[]string, guardStartPosition []int) map[string]struct{} {
	guardPosition := []int{guardStartPosition[0], guardStartPosition[1]}
	width := len((*rows)[0])
	height := len(*rows)

	visitedPositions := make(map[string]struct{})
	direction := Up

	for guardPosition[0] >= 0 && guardPosition[0] < width && guardPosition[1] >= 0 && guardPosition[1] < height {
		positionKey := fmt.Sprintf("%d-%d", guardPosition[0], guardPosition[1])
		visitedPositions[positionKey] = struct{}{}

		var nextPosition []int

		move := getMoveForDirection(direction)
		nextPosition = applyMove(guardPosition, move)

		if nextPosition[0] < 0 || nextPosition[0] >= width || nextPosition[1] < 0 || nextPosition[1] >= height {
			break
		}

		if (*rows)[nextPosition[1]][nextPosition[0]] == '#' {
			// Turn right
			direction = (direction + 1) % 4
		} else {
			guardPosition = nextPosition
		}
	}

	return visitedPositions
}

func Part1(rows *[]string) (int, error) {
	guardPosition, err := getGuardStartPosition(rows)

	if err != nil {
		return -1, err
	}

	visitedPositions := getVisitedPositions(rows, guardPosition)

	return len(visitedPositions), nil
}

func Part2(rows *[]string) (int, error) {
	guardStartPosition, err := getGuardStartPosition(rows)
	width := len((*rows)[0])
	height := len(*rows)

	if err != nil {
		return -1, err
	}

	visitedPositions := getVisitedPositions(rows, guardStartPosition)

	var validBlockades [][]int

	for pos := range visitedPositions {
		posStrings := strings.Split(pos, "-")

		x, err := strconv.Atoi(posStrings[0])
		if err != nil {
			return -1, err
		}

		y, err := strconv.Atoi(posStrings[1])
		if err != nil {
			return -1, err
		}

		if x == guardStartPosition[0] && y == guardStartPosition[1] {
			continue
		}

		blockade := []int{x, y}
		loopCheck := make(map[string]struct{})
		guardPosition := []int{guardStartPosition[0], guardStartPosition[1]}
		direction := Up

		for guardPosition[0] >= 0 && guardPosition[0] < width && guardPosition[1] >= 0 && guardPosition[1] < height {
			positionKey := fmt.Sprintf("%d-%d-%d", guardPosition[0], guardPosition[1], direction)
			_, hasLooped := loopCheck[positionKey]

			if hasLooped {
				validBlockades = append(validBlockades, blockade)
				break
			}

			loopCheck[positionKey] = struct{}{}

			var nextPosition []int

			move := getMoveForDirection(direction)
			nextPosition = applyMove(guardPosition, move)

			if nextPosition[0] < 0 || nextPosition[0] >= width || nextPosition[1] < 0 || nextPosition[1] >= height {
				break
			}

			if (*rows)[nextPosition[1]][nextPosition[0]] == '#' || (nextPosition[0] == blockade[0] && nextPosition[1] == blockade[1]) {
				// Turn right
				direction = (direction + 1) % 4
			} else {
				guardPosition = nextPosition
			}
		}
	}

	return len(validBlockades), nil
}
