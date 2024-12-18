package day_16

import (
	"fmt"
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
	return fmt.Sprintf("%d-%d", vector[0], vector[1])
}

func decodeVector(vectorStr string) []int {
	split := strings.Split(vectorStr, "-")

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

const maxInt int = int(^uint(0) >> 1)

type state struct {
	Field            []int
	CurrentDirection []int
	Points           int
	Path             [][]int
}

func (s state) NextState(destination []int) state {
	moveDirection := calculateVector(s.Field, destination)

	pointsForDirectionChange := PointsForDirectionChange(s.CurrentDirection, moveDirection)
	points := 1 + pointsForDirectionChange

	nextState := state{
		Field:            destination,
		CurrentDirection: moveDirection,
		Points:           s.Points + points,
		Path:             append(s.Path, destination),
	}

	return nextState
}

func PointsForDirectionChange(currentDirection []int, direction []int) int {
	vecSum := []int{currentDirection[0] + direction[0], currentDirection[1] + direction[1]}

	points := 0

	if vecSum[0] == 0 && vecSum[1] == 0 {
		points += 2000
	}

	if vecSum[0] != 0 && vecSum[1] != 0 {
		points += 1000
	}

	return points
}

func (s state) KeyWithDirection() string {
	key := encodeVector(s.Field)
	key = fmt.Sprintf("%s-%d-%d", key, s.CurrentDirection[0], s.CurrentDirection[1])

	return key
}

func (s state) Key() string {
	return encodeVector(s.Field)
}

func (s state) IsStrictlyBetter(otherState state) bool {
	if vectorEq(s.CurrentDirection, otherState.CurrentDirection) {
		return s.Points < otherState.Points
	} else {
		pointsForDirectionChange := PointsForDirectionChange(otherState.CurrentDirection, s.CurrentDirection)

		return s.Points < otherState.Points+pointsForDirectionChange
	}
}

func (s state) IsBetterOrEqual(otherState state) bool {
	if vectorEq(s.CurrentDirection, otherState.CurrentDirection) {
		return s.Points <= otherState.Points
	} else {
		pointsForDirectionChange := PointsForDirectionChange(otherState.CurrentDirection, s.CurrentDirection)

		return s.Points <= otherState.Points+pointsForDirectionChange
	}
}

func vectorEq(a []int, b []int) bool {
	return a[0] == b[0] && a[1] == b[1]
}

func FindShortestPath(current state, maze *[]string, pointsMap *map[string]state) {
	currentLowestState, alreadyHasPoints := (*pointsMap)[current.Key()]
	if alreadyHasPoints && !current.IsBetterOrEqual(currentLowestState) {
		return
	}

	neighbors := getNeighborsForField(current, maze, nil)

	for _, n := range neighbors {
		nextState := current.NextState(n)
		nextStateKey := nextState.Key()

		currentBestPointsForNeighbor, hasBestState := (*pointsMap)[nextStateKey]
		if hasBestState {
			if nextState.IsStrictlyBetter(currentBestPointsForNeighbor) {
				(*pointsMap)[nextStateKey] = nextState
				if (*maze)[nextState.Field[1]][nextState.Field[0]] != 'E' {
					FindShortestPath(nextState, maze, pointsMap)
				}
			}
		} else {
			(*pointsMap)[nextStateKey] = nextState
			if (*maze)[nextState.Field[1]][nextState.Field[0]] != 'E' {
				FindShortestPath(nextState, maze, pointsMap)
			}
		}
	}
}

func FindAllShortestPaths(current state, maze *[]string, pointsMap *map[string]state, shortestPathPoints int, endStates *[]state) {
	if current.Points > shortestPathPoints {
		return
	}

	neighbors := getNeighborsForField(current, maze, nil)

	for _, n := range neighbors {
		nextState := current.NextState(n)
		nextStateKey := nextState.Key()

		currentBestPointsForNeighbor, hasBestState := (*pointsMap)[nextStateKey]
		if hasBestState {
			if nextState.IsBetterOrEqual(currentBestPointsForNeighbor) {
				(*pointsMap)[nextStateKey] = nextState
				if (*maze)[nextState.Field[1]][nextState.Field[0]] != 'E' {
					FindAllShortestPaths(nextState, maze, pointsMap, shortestPathPoints, endStates)
				} else {
					*endStates = append(*endStates, nextState)
				}
			}
		} else {
			(*pointsMap)[nextStateKey] = nextState
			if (*maze)[nextState.Field[1]][nextState.Field[0]] != 'E' {
				FindAllShortestPaths(nextState, maze, pointsMap, shortestPathPoints, endStates)
			} else {
				*endStates = append(*endStates, nextState)
			}
		}
	}
}

func Part1(maze *[]string) (int, error) {
	pointsMap := make(map[string]state)
	start, finish, err := findStartAndFinish(maze)

	if err != nil {
		return -1, err
	}

	FindShortestPath(state{
		Field:            start,
		CurrentDirection: []int{1, 0},
		Points:           0,
		Path:             [][]int{start},
	}, maze, &pointsMap)

	points, hasPoints := pointsMap[encodeVector(finish)]
	if !hasPoints {
		return -1, fmt.Errorf("Shortest path not found")
	}

	return points.Points, nil
}

func Part2(maze *[]string) (int, error) {
	pointsMap := make(map[string]state)
	start, finish, err := findStartAndFinish(maze)

	if err != nil {
		return -1, err
	}

	FindShortestPath(state{
		Field:            start,
		CurrentDirection: []int{1, 0},
		Points:           0,
		Path:             [][]int{start},
	}, maze, &pointsMap)

	shortestPathLen, hasPoints := pointsMap[encodeVector(finish)]

	if !hasPoints {
		return -1, fmt.Errorf("Shortest path not found")
	}

	fmt.Println("Shortest path", shortestPathLen)

	var endStates []state
	pointsMap = make(map[string]state)

	FindAllShortestPaths(
		state{
			Field:            start,
			CurrentDirection: []int{1, 0},
			Points:           0,
			Path:             [][]int{start},
		}, maze, &pointsMap, shortestPathLen.Points, &endStates,
	)
	fmt.Printf("end states: %v, points map: %v\n", endStates, pointsMap)

	visitedFieldsMap := make(map[string]struct{})

	for _, endState := range endStates {
		fmt.Printf("Path: %v\n", endState.Path)
		for _, field := range endState.Path {
			visitedFieldsMap[encodeVector(field)] = struct{}{}
		}
	}
	visitedFields := make([][]int, 0, len(visitedFieldsMap))
	for key := range visitedFieldsMap {
		visitedFields = append(visitedFields, decodeVector(key))
	}

	printMazeWithVisitedFields(maze, visitedFields)

	return len(visitedFieldsMap), nil
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

func calculateVector(a []int, b []int) []int {
	return []int{b[0] - a[0], b[1] - a[1]}
}

func validateNeighbor(current state, neighbor []int, maze *[]string, visited *map[string]struct{}) bool {
	mazeHeight := len(*maze)
	mazeWidth := len((*maze)[0])
	moveDirection := calculateVector(current.Field, neighbor)

	vecDiff := []int{moveDirection[0] + current.CurrentDirection[0], moveDirection[1] + current.CurrentDirection[1]}
	if vecDiff[0] == 0 && vecDiff[1] == 0 {
		return false
	}

	if neighbor[0] < 0 || neighbor[0] >= mazeWidth || neighbor[1] < 0 || neighbor[1] >= mazeHeight {
		return false
	}

	if (*maze)[neighbor[1]][neighbor[0]] == '#' {
		return false
	}

	if visited != nil {
		nextState := current.NextState(neighbor)
		_, alreadyVisited := (*visited)[nextState.KeyWithDirection()]

		return !alreadyVisited
	}
	return true
}

func getNeighborsForField(current state, maze *[]string, visited *map[string]struct{}) [][]int {
	neighbors := make([][]int, 0, 4)
	x := current.Field[0]
	y := current.Field[1]

	n1 := []int{x + 1, y}
	if validateNeighbor(current, n1, maze, visited) {
		neighbors = append(neighbors, n1)
	}

	n2 := []int{x - 1, y}
	if validateNeighbor(current, n2, maze, visited) {
		neighbors = append(neighbors, n2)
	}

	n3 := []int{x, y + 1}
	if validateNeighbor(current, n3, maze, visited) {
		neighbors = append(neighbors, n3)
	}

	n4 := []int{x, y - 1}
	if validateNeighbor(current, n4, maze, visited) {
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
