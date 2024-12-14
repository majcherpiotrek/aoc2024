package day_14

import (
	"fmt"
	"strconv"
	"strings"
)

type Robot struct {
	Position []int
	Velocity []int
}

func (r *Robot) move(mapWidth int, mapHeight int) {
	nextPosition := addVectors(r.Position, r.Velocity)

	//fmt.Printf("Moving robot: (%d, %d) -> (%d, %d) = (%d, %d)", r.Position[0], r.Position[1], r.Velocity[0], r.Velocity[1], nextPosition[0], nextPosition[1])
	x := nextPosition[0]

	if x >= mapWidth {
		x = x % mapWidth
	}

	if x < 0 {
		x = mapWidth + x
	}

	y := nextPosition[1]

	if y >= mapHeight {
		y = y % mapHeight
	}

	if y < 0 {
		y = mapHeight + y
	}
	//fmt.Printf(", actually (%d, %d)\n", x, y)

	r.Position = []int{x, y}
}

func addVectors(a []int, b []int) []int {
	return []int{a[0] + b[0], a[1] + b[1]}
}

func parseInput(input *[]string) ([]*Robot, error) {
	robotList := make([]*Robot, 0, len(*input))
	for _, row := range *input {
		robot, err := parseLine(row)

		if err != nil {
			return []*Robot{}, err
		}

		robotList = append(robotList, &robot)

	}

	return robotList, nil
}

func parseLine(line string) (Robot, error) {
	fields := strings.Fields(line)

	positionSplit := strings.Split(fields[0], "=")
	velocitySplit := strings.Split(fields[1], "=")

	positon, err := parseVector(positionSplit[1])
	if err != nil {
		return Robot{}, err
	}
	velocity, err := parseVector(velocitySplit[1])
	if err != nil {
		return Robot{}, err
	}

	return Robot{
		Position: positon,
		Velocity: velocity,
	}, nil
}

func parseVector(str string) ([]int, error) {
	split := strings.Split(str, ",")

	if len(split) != 2 {
		return []int{}, fmt.Errorf("Expected a vector of length 2")
	}

	x, err := strconv.Atoi(split[0])
	if err != nil {
		return []int{}, err
	}

	y, err := strconv.Atoi(split[1])
	if err != nil {
		return []int{}, err
	}

	return []int{x, y}, nil
}

func sumRobots(matrix [][]int) []int {
	halfHeight := len(matrix) / 2
	halfWidth := len(matrix[0]) / 2

	sumQ1 := 0
	sumQ2 := 0
	sumQ3 := 0
	sumQ4 := 0

	for y := 0; y < len(matrix); y++ {
		for x := 0; x < len(matrix[y]); x++ {

			if x < halfWidth && y < halfHeight {
				sumQ1 += matrix[y][x]
			}

			if x > halfWidth && y < halfHeight {
				sumQ2 += matrix[y][x]
			}

			if x > halfWidth && y > halfHeight {
				sumQ3 += matrix[y][x]
			}

			if x < halfWidth && y > halfHeight {
				sumQ4 += matrix[y][x]
			}
		}
	}

	return []int{sumQ1, sumQ2, sumQ3, sumQ4}
}

func printRobots(robots *[]*Robot) [][]int {
	matrix := make([][]int, 0, 103)

	for i := 0; i < 103; i++ {
		row := make([]int, 0, 101)
		for j := 0; j < 101; j++ {
			row = append(row, 0)
		}
		matrix = append(matrix, row)
	}

	for _, robot := range *robots {
		x := robot.Position[0]
		y := robot.Position[1]
		matrix[y][x] += 1
	}

	for _, row := range matrix {
		fmt.Println(row)
	}

	return matrix
}

func Part1(rows *[]string) (int, error) {

	robots, err := parseInput(rows)

	if err != nil {
		return -1, err
	}

	//fmt.Println("START")
	//matrix := printRobots(&robots)
	//fmt.Println("-----------")

	for i := 0; i < 100; i++ {
		for _, robot := range robots {
			robot.move(101, 103)
		}
		//fmt.Printf("After %d second\n", i+1)
		//matrix = printRobots(&robots)
		//fmt.Println("-----------")
	}

	matrix := printRobots(&robots)
	sumInQuadrants := sumRobots(matrix)

	fmt.Println(sumInQuadrants)

	sum := -1
	for _, s := range sumInQuadrants {
		if sum == -1 {
			sum = s
		} else {

			sum *= s
		}
	}

	return sum, nil
}

func Part2(rows *[]string) (int, error) {

	return -1, fmt.Errorf("not implemented")
}
