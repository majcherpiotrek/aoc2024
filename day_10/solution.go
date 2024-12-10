package day_10

import (
	"fmt"
	"strconv"
)

func parseInput(input *[]string) ([][]int, error) {
	result := make([][]int, 0, len(*input))

	for _, row := range *input {
		bytes := []byte(row)

		parsedRow := make([]int, 0, len(row))

		for _, b := range bytes {
			str := string(b)
			num, err := strconv.Atoi(str)

			if err != nil {
				return [][]int{}, err
			}

			parsedRow = append(parsedRow, num)
		}

		result = append(result, parsedRow)
	}

	return result, nil
}

type Direction int

const (
	Up = iota
	Right
	Down
	Left
)

func followHikingTrail(startingPoint []int, path *[][]int, direction Direction, topoMap *[][]int, uniquePeaksReached *map[string]struct{}) int {
	mapWidth := len((*topoMap)[0])
	mapHeight := len(*topoMap)
	currentFieldHeight := (*topoMap)[startingPoint[1]][startingPoint[0]]
	*path = append(*path, []int{startingPoint[0], startingPoint[1], currentFieldHeight})

	var nextFieldCoordinates []int

	switch direction {
	case Up:
		nextFieldCoordinates = []int{startingPoint[0], startingPoint[1] - 1}
	case Right:
		nextFieldCoordinates = []int{startingPoint[0] + 1, startingPoint[1]}
	case Down:
		nextFieldCoordinates = []int{startingPoint[0], startingPoint[1] + 1}
	case Left:
		nextFieldCoordinates = []int{startingPoint[0] - 1, startingPoint[1]}
	}

	if nextFieldCoordinates[0] >= 0 && nextFieldCoordinates[0] < mapWidth && nextFieldCoordinates[1] >= 0 && nextFieldCoordinates[1] < mapHeight {
		nextFieldHeight := (*topoMap)[nextFieldCoordinates[1]][nextFieldCoordinates[0]]
		// fmt.Printf("Next field height: %d\n", nextFieldHeight)

		heightDiff := nextFieldHeight - currentFieldHeight

		if heightDiff != 1 {
			// fmt.Println("Cutting this trail")
			return 0
		}

		if nextFieldHeight == 9 {
			*path = append(*path, []int{nextFieldCoordinates[0], nextFieldCoordinates[1], nextFieldHeight})
			cs := checksum(path)
			(*uniquePeaksReached)[fmt.Sprintf("%d-%d", nextFieldCoordinates[0], nextFieldCoordinates[1])] = struct{}{}
			fmt.Printf("Path found %v, checksum %d\n", *path, cs)
			return 1
		}

		// fmt.Println("No peak reached, trying other paths")
		clonedUp := clonePath(path)
		clonedRight := clonePath(path)
		clonedDown := clonePath(path)
		clonedLeft := clonePath(path)

		return followHikingTrail(nextFieldCoordinates, &clonedUp, Up, topoMap, uniquePeaksReached) +
			followHikingTrail(nextFieldCoordinates, &clonedRight, Right, topoMap, uniquePeaksReached) +
			followHikingTrail(nextFieldCoordinates, &clonedDown, Down, topoMap, uniquePeaksReached) +
			followHikingTrail(nextFieldCoordinates, &clonedLeft, Left, topoMap, uniquePeaksReached)

	}

	return 0
}

func clonePath(path *[][]int) [][]int {
	newPath := make([][]int, 0, len(*path))
	newPath = append(newPath, *path...)
	return newPath
}

func checksum(path *[][]int) int {
	checksum := 0

	for i, field := range *path {
		checksum = checksum + (field[0]*i+field[1]*(i+1))*i
	}

	return checksum
}

func Part1(rows *[]string) (int, error) {
	topographicMap, err := parseInput(rows)

	if err != nil {
		return -1, err
	}

	fmt.Println(topographicMap)

	var trailHeadPoints []int

	for y := 0; y < len(topographicMap); y++ {
		row := topographicMap[y]
		for x := 0; x < len(row); x++ {
			fieldHeight := row[x]
			//fmt.Printf("[%d, %d] - %d\n", x, y, fieldHeight)

			if fieldHeight == 0 {
				uniquePeakseReached := make(map[string]struct{})
				coords := []int{x, y}
				followHikingTrail(coords, &[][]int{}, Up, &topographicMap, &uniquePeakseReached)
				followHikingTrail(coords, &[][]int{}, Right, &topographicMap, &uniquePeakseReached)
				followHikingTrail(coords, &[][]int{}, Down, &topographicMap, &uniquePeakseReached)
				followHikingTrail(coords, &[][]int{}, Left, &topographicMap, &uniquePeakseReached)

				trailHeadPoints = append(trailHeadPoints, len(uniquePeakseReached))
			}

		}
	}

	fmt.Println(trailHeadPoints)
	sum := 0
	for _, points := range trailHeadPoints {
		sum += points
	}

	return sum, nil
}

func Part2(rows *[]string) (int, error) {
	topographicMap, err := parseInput(rows)

	if err != nil {
		return -1, err
	}

	fmt.Println(topographicMap)

	var trailHeadPoints []int
	var trailHeadRatings []int

	for y := 0; y < len(topographicMap); y++ {
		row := topographicMap[y]
		for x := 0; x < len(row); x++ {
			fieldHeight := row[x]
			//fmt.Printf("[%d, %d] - %d\n", x, y, fieldHeight)

			if fieldHeight == 0 {
				uniquePeakseReached := make(map[string]struct{})

				coords := []int{x, y}
				rating := followHikingTrail(coords, &[][]int{}, Up, &topographicMap, &uniquePeakseReached)
				rating += followHikingTrail(coords, &[][]int{}, Right, &topographicMap, &uniquePeakseReached)
				rating += followHikingTrail(coords, &[][]int{}, Down, &topographicMap, &uniquePeakseReached)
				rating += followHikingTrail(coords, &[][]int{}, Left, &topographicMap, &uniquePeakseReached)

				trailHeadPoints = append(trailHeadPoints, len(uniquePeakseReached))
				trailHeadRatings = append(trailHeadRatings, rating)
			}

		}
	}

	fmt.Println(trailHeadRatings)
	sum := 0
	for _, points := range trailHeadRatings {
		sum += points
	}

	return sum, nil
}
