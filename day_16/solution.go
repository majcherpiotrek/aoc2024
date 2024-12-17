package day_16

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func findStartAndFinish(maze *[]string) ([]int, []int, error) {
	var start []int
	var finish []int

	for y, row := range *maze {
		for x := 0; x < len(row); x++ {
			element := row[x]
			if element == 'S' {
				start = []int{x, y}
			} else if element == 'E' {
				finish = []int{x, y}
			}

			if len(start) > 0 && len(finish) > 0 {
				break
			}
		}
	}

	if len(start) == 0 || len(finish) == 0 {
		return []int{}, []int{}, fmt.Errorf("Invalid input, can't find start and finish in the maze")
	}

	return start, finish, nil
}

type Direction int

const (
	Up = iota
	Right
	Down
	Left
)

func encodeVector(vector []int) string {
	return fmt.Sprintf("%d,%d", vector[0], vector[1])
}

func decodeVector(vectorStr string) []int {
	split := strings.Split(vectorStr, ",")

	if len(split) != 2 {
		panic(fmt.Sprintf("Invalid vector string %s", vectorStr))
	}

	x, err := strconv.Atoi(split[0])
	if err != nil {
		panic(err)
	}

	y, err := strconv.Atoi(split[1])

	if err != nil {
		panic(err)
	}

	return []int{x, y}
}

func directionFromVector(vector []int) Direction {
	switch encodeVector(vector) {
	case "1,0":
		return Right
	case "-1,0":
		return Left
	case "0,1":
		return Down
	case "0,-1":
		return Up
	default:
		panic(fmt.Sprintf("Invalid vector %v, can't convert to Direction", vector))
	}
}

func calculatePoints(a []int, b []int, currentDirection Direction) int {
	vector := []int{b[0] - a[0], b[1] - a[1]}
	moveDirection := directionFromVector(vector)

	// Why modulo?
	// 0 -> 0
	// 1 -> 1000
	// 2 -> 2000
	// 3 -> 1000
	directionChange := int(math.Abs(float64(currentDirection-moveDirection))) % 2
	//fmt.Printf("%sDirection change: %d\n", padding, directionChange)

	directionChangePoints := directionChange * 1000

	return 1 + directionChangePoints
}

const maxInt int = int(^uint(0) >> 1)

func GoThroughMaze(currentField []int, currentPoints int, currentDirection Direction, maze *[]string, pointsMap *map[string]int, depth int) {
	neighbors := make([][]int, 0, 4)

	//padding := ""
	//for i := 0; i < depth; i++ {
	//	padding = fmt.Sprintf("%s ", padding)
	//}

	mazeHeight := len(*maze)
	mazeWidth := len((*maze)[0])
	x := currentField[0]
	y := currentField[1]

	currentLowestPointsForField, alreadyHasPoints := (*pointsMap)[encodeVector(currentField)]
	if alreadyHasPoints && currentLowestPointsForField < currentPoints {
		return
	}

	if x+1 < mazeWidth {
		neighbors = append(neighbors, []int{x + 1, y})
	}
	if x-1 >= 0 {
		neighbors = append(neighbors, []int{x - 1, y})
	}
	if y+1 < mazeHeight {
		neighbors = append(neighbors, []int{x, y + 1})
	}
	if y-1 >= 0 {
		neighbors = append(neighbors, []int{x, y - 1})
	}

	for _, n := range neighbors {
		neighborValue := (*maze)[n[1]][n[0]]
		//fmt.Printf("%sCurrentField: %v, currentPoints: %d, currentDirection: %d, neighbor: %v, value: %c\n", padding, currentField, currentPoints, currentDirection, n, neighborValue)

		if neighborValue == '#' {
			//fmt.Printf("%sSkipping, wall ahead\n\n", padding)
			continue
		}

		key := encodeVector(n)
		points, hasPoints := (*pointsMap)[key]
		if !hasPoints {
			points = maxInt
		}

		vector := []int{n[0] - currentField[0], n[1] - currentField[1]}
		moveDirection := directionFromVector(vector)
		//fmt.Printf("%sMove direction: %d\n", padding, moveDirection)
		pointsForMove := calculatePoints(currentField, n, currentDirection)
		//fmt.Printf("%sPoints for move: %d\n", padding, pointsForMove)
		updatedPoints := pointsForMove + currentPoints

		//fmt.Printf("%sCurrent neighbor points: %d, updatedPoints: %d\n\n", padding, points, updatedPoints)

		if updatedPoints < points {
			(*pointsMap)[key] = updatedPoints
			if neighborValue != 'E' {
				GoThroughMaze(n, updatedPoints, moveDirection, maze, pointsMap, depth+1)
			}
		}
	}
}

func Part1(maze *[]string) (int, error) {
	pointsMap := make(map[string]int)
	start, finish, err := findStartAndFinish(maze)

	if err != nil {
		return -1, err
	}

	GoThroughMaze(start, 0, Right, maze, &pointsMap, 0)
	//GoThroughMaze(start, 1000, Up, maze, &pointsMap)
	//GoThroughMaze(start, 1000, Down, maze, &pointsMap)
	//GoThroughMaze(start, 2000, Left, maze, &pointsMap)

	//fmt.Printf("%v", pointsMap)

	result, hasResult := pointsMap[encodeVector(finish)]
	if !hasResult {
		return -1, fmt.Errorf("Could not find path")
	}

	return result, nil
}

func Part2(rows *[]string) (int, error) {

	return -1, fmt.Errorf("not implemented")
}
