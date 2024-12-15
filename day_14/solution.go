package day_14

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const WIDTH int = 101
const HEIGHT int = 103

type Robot struct {
	Position []int
	Velocity []int
}

func (r *Robot) move(mapWidth int, mapHeight int) {
	nextPosition := addVectors(r.Position, r.Velocity)

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

func calculateMatrix(robots *[]*Robot) [][]int {
	matrix := make([][]int, 0, HEIGHT)

	for i := 0; i < HEIGHT; i++ {
		row := make([]int, 0, WIDTH)
		for j := 0; j < WIDTH; j++ {
			row = append(row, 0)
		}
		matrix = append(matrix, row)
	}

	for _, robot := range *robots {
		x := robot.Position[0]
		y := robot.Position[1]
		matrix[y][x] += 1
	}

	return matrix
}

func printMatrix(matrix *[][]int, subMatrix *[][]int) {
	if subMatrix != nil {
		minX := (*subMatrix)[0][0]
		minY := (*subMatrix)[0][1]
		maxX := (*subMatrix)[1][0]
		maxY := (*subMatrix)[1][1]

		for y, row := range *matrix {
			str := ""
			for x, num := range row {
				if x >= minX && x <= maxX && y >= minY && y <= maxY {
					if num == 0 {
						str = fmt.Sprintf("%s.", str)
					} else {
						str = fmt.Sprintf("%s%d", str, num)
					}
				} else {

					str = fmt.Sprintf("%sx", str)
				}
			}
			fmt.Println(str)
		}
	} else {
		for _, row := range *matrix {
			str := ""
			for _, num := range row {
				if num == 0 {
					str = fmt.Sprintf("%s.", str)
				} else {
					str = fmt.Sprintf("%s%d", str, num)
				}
			}
			fmt.Println(str)
		}
	}
	fmt.Println("-")
}

func Part1(rows *[]string) (int, error) {
	robots, err := parseInput(rows)

	if err != nil {
		return -1, err
	}

	fmt.Println("START")
	matrix := calculateMatrix(&robots)
	fmt.Println("-----------")

	for i := 0; i < 100; i++ {
		for _, robot := range robots {
			robot.move(WIDTH, HEIGHT)
		}
		fmt.Printf("After %d second\n", i+1)
		matrix = calculateMatrix(&robots)
		fmt.Println("-----------")
	}

	matrix = calculateMatrix(&robots)
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

func sliceEq(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	isEqual := false

	for i := 0; i < len(a); i++ {
		isEqual = a[i] == b[i]

		if !isEqual {
			break
		}
	}

	return isEqual
}

func searchForChristmasTreePattern(start []int, matrix *[][]int, visited *map[string]struct{}) bool {
	pattern := [][]int{
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 1, 1, 1, 0, 0},
		{0, 1, 1, 1, 1, 1, 0},
		{1, 1, 1, 1, 1, 1, 1},
	}

	x := start[0]
	y := start[1]

	key := fmt.Sprintf("%d-%d", x, y)

	_, alreadyVisited := (*visited)[key]

	if alreadyVisited {
		return false
	}

	if y < 0 || y+3 >= len(*matrix) {
		return false
	}

	if x < 0 || x+6 >= len((*matrix)[y]) {
		return false
	}

	hasPattern := false

	//	clearConsole()
	//	printMatrix(matrix, &[][]int{{x, y}, {x + 6, y + 3}})

	for i := 0; i <= 3; i++ {
		row := (*matrix)[y+i]
		rowToCompare := row[x : x+7]

		hasPattern = sliceEq(rowToCompare, pattern[i])

		if !hasPattern {
			break
		}
	}

	//fmt.Println("WHOLE PATTERN", hasPattern)

	if hasPattern {
		return true
	}

	(*visited)[key] = struct{}{}

	right := searchForChristmasTreePattern([]int{x + 1, y}, matrix, visited)

	if right {
		return true
	}

	down := searchForChristmasTreePattern([]int{x, y + 1}, matrix, visited)

	if down {
		return true
	}

	left := searchForChristmasTreePattern([]int{x - 1, y}, matrix, visited)

	if left {
		return true
	}

	top := searchForChristmasTreePattern([]int{x, y - 1}, matrix, visited)

	if top {
		return true
	}

	return false
}

func Part2(rows *[]string) (int, error) {
	robots, err := parseInput(rows)

	if err != nil {
		return -1, err
	}

	matrix := calculateMatrix(&robots)

	frame := 1
	christmasTreeFrame := -1

	//testMatrix := [][]int{
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
	//	{0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//}

	//testVisited := make(map[string]struct{})
	//testRes := searchForChristmasTreePattern([]int{0, 0}, &testMatrix, &testVisited)

	//fmt.Println("TEST RES", testRes)

	fmt.Printf("START\n\n")
	printMatrix(&matrix, nil)

	for frame < 10000 {
		for _, robot := range robots {
			robot.move(WIDTH, HEIGHT)
		}
		matrix = calculateMatrix(&robots)
		clearConsole()
		printMatrix(&matrix, nil)
		fmt.Printf("Frame %d\n\n", frame)

		visited := make(map[string]struct{})
		patternFound := searchForChristmasTreePattern([]int{0, 0}, &matrix, &visited)

		if patternFound {
			christmasTreeFrame = frame
			break
		}

		frame++
	}

	return christmasTreeFrame, nil
}

func clearConsole() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}
