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

func calculatePoints(a []int, b []int, currentDirection Direction) int {
	vector := []int{b[0] - a[0], b[1] - a[1]}
	moveDirection := directionFromVector(vector)

	directionChange := int(math.Abs(float64(currentDirection - moveDirection)))

	switch directionChange {
	case 1:
		return 1001
	case 2:
		return 2001
	case 3:
		return 1001
	default:
		return 1
	}
}

const maxInt int = int(^uint(0) >> 1)

type state struct {
	Field            []int
	CurrentDirection []int
	Points           int
	Path             [][]int
}

func (s state) NextState(destination []int) state {
	moveDirection := calculateVector(s.Field, destination)

	points := 1

	vecSum := []int{s.CurrentDirection[0] + moveDirection[0], s.CurrentDirection[1] + moveDirection[1]}

	if vecSum[0] == 0 && vecSum[1] == 0 {
		points += 2000
	}

	if vecSum[0] != 0 && vecSum[1] != 0 {
		points += 1000
	}

	nextState := state{
		Field:            destination,
		CurrentDirection: moveDirection,
		Points:           s.Points + points,
		Path:             append(s.Path, destination),
	}

	return nextState
}

func FindShortestPath(currentField []int, currentPoints int, currentDirection Direction, maze *[]string, pointsMap *map[string]int, visited *map[string]struct{}, depth int) {

	currentLowestPointsForField, alreadyHasPoints := (*pointsMap)[encodeVector(currentField)]
	if alreadyHasPoints && currentLowestPointsForField < currentPoints {
		return
	}

	neighbors := getNeighborsForField(currentField, maze, visited)

	for _, n := range neighbors {
		neighborValue := (*maze)[n[1]][n[0]]

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
				FindShortestPath(n, updatedPoints, moveDirection, maze, pointsMap, depth+1)
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

	FindShortestPath(start, 0, Right, maze, &pointsMap, 0)
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

func printMazeWithVisitedFields(maze *[]string, visitedFields [][]int) {
	mazeToPrint := make([][]byte, 0, len(*maze))

	for _, row := range *maze {
		mazeToPrint = append(mazeToPrint, []byte(row))
	}

	for _, field := range visitedFields {
		mazeToPrint[field[1]][field[0]] = 'O'
	}

	for _, row := range mazeToPrint {
		fmt.Println(string(row))
	}
}

func encodeVisitedKey(field []int, direction Direction) string {
	visitedKey := encodeVector(field)
	visitedKey = fmt.Sprintf("%s-%d", visitedKey, direction)

	return visitedKey
}

func calculateVector(a []int, b []int) []int {
	return []int{b[0] - a[0], b[1] - a[1]}
}

func validateNeighbor(field []int, neighbor []int, maze *[]string, visited *map[string]struct{}) bool {
	mazeHeight := len(*maze)
	mazeWidth := len((*maze)[0])

	if neighbor[0] < 0 || neighbor[0] >= mazeWidth || neighbor[1] < 0 || neighbor[1] >= mazeHeight {
		return false
	}

	if (*maze)[neighbor[1]][neighbor[0]] == '#' {
		return false
	}

	if visited != nil {
		vector := calculateVector(field, neighbor)
		direction := directionFromVector(vector)
		visitedKey := encodeVisitedKey(neighbor, direction)
		_, alreadyVisited := (*visited)[visitedKey]

		return !alreadyVisited
	}
	return true
}

// [1 13] [1 12] [1 11] [1 10] [1 9] [2 9] [3 9] [3 8] [3 7] [4 7] [5 7]
func getNeighborsForField(field []int, maze *[]string, visited *map[string]struct{}) [][]int {
	neighbors := make([][]int, 0, 4)
	x := field[0]
	y := field[1]

	n1 := []int{x + 1, y}
	if validateNeighbor(field, n1, maze, visited) {
		neighbors = append(neighbors, n1)
	}

	n2 := []int{x - 1, y}
	if validateNeighbor(field, n2, maze, visited) {
		neighbors = append(neighbors, n2)
	}

	n3 := []int{x, y + 1}
	if validateNeighbor(field, n3, maze, visited) {
		neighbors = append(neighbors, n3)
	}

	n4 := []int{x, y - 1}
	if validateNeighbor(field, n4, maze, visited) {
		neighbors = append(neighbors, n4)
	}

	return neighbors
}

func getPadding(depth int) string {
	padding := ""
	for i := 0; i < depth; i++ {
		padding = fmt.Sprintf("%s ", padding)
	}
	return padding
}
