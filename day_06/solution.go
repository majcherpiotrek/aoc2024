package day_06

import (
	"fmt"
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

func Part1(rows *[]string) (int, error) {
	guardPosition := make([]int, 0, 2)
	width := len((*rows)[0])
	height := len(*rows)

	for y, row := range *rows {

		guardX := strings.IndexByte(row, '^')

		if guardX != -1 {
			guardPosition = []int{guardX, y}
			break
		}
	}

	if len(guardPosition) == 0 {
		return -1, fmt.Errorf("Could not find the guard")
	}

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

	return len(visitedPositions), nil
}

func Part2(rows *[]string) (int, error) {
	return -1, fmt.Errorf("not implemented")
}
