package day_08

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func findAntennas(input *[]string) map[byte][][]int {
	antennasMap := make(map[byte][][]int)

	for y, row := range *input {
		bytes := []byte(row)

		for x, b := range bytes {
			if b != '.' {
				locations, hasEntry := antennasMap[b]

				if hasEntry {
					locations = append(locations, []int{x, y})
					antennasMap[b] = locations
				} else {
					antennasMap[b] = [][]int{{x, y}}
				}

			}
		}
	}

	return antennasMap
}

func getVector(a []int, b []int) []int {
	vx := b[0] - a[0]
	vy := b[1] - a[1]
	return []int{vx, vy}
}

func addVector(a []int, vect []int) []int {
	return []int{a[0] + vect[0], a[1] + vect[1]}
}

func flipVector(vect []int) []int {
	return []int{-1 * vect[0], -1 * vect[1]}
}

func isInBounds(a []int, xMax int, yMax int) bool {
	return a[0] >= 0 && a[0] <= xMax && a[1] >= 0 && a[1] <= yMax
}

func findAntinodesForPoints(a []int, b []int, xMax int, yMax int) [][]int {
	antinode1 := addVector(b, getVector(a, b))
	antinode2 := addVector(a, getVector(b, a))

	nodes := make([][]int, 0, 2)

	if isInBounds(antinode1, xMax, yMax) {
		nodes = append(nodes, antinode1)
	}

	if isInBounds(antinode2, xMax, yMax) {
		nodes = append(nodes, antinode2)
	}

	return nodes
}

func findAntinodesForPointsPart2(a []int, b []int, xMax int, yMax int) [][]int {
	var nodes [][]int

	vect := minimizeVector(getVector(a, b))
	fmt.Printf("vect AB %v", vect)

	for node := addVector(a, vect); isInBounds(node, xMax, yMax); node = addVector(node, vect) {
		nodes = append(nodes, node)
	}

	vect = minimizeVector(getVector(b, a))
	fmt.Printf("vect BA %v", vect)

	for node := addVector(b, vect); isInBounds(node, xMax, yMax); node = addVector(node, vect) {
		nodes = append(nodes, node)
	}

	return nodes

}

func gcd(a int, b int) int {
	if a == 0 {
		return int(math.Abs(float64(b)))
	}

	return gcd(b%a, a)
}

func minimizeVector(vect []int) []int {
	divisor := gcd(vect[0], vect[1])
	fmt.Printf("GCD %d, %d = %d", vect[0], vect[1], divisor)

	return []int{vect[0] / divisor, vect[1] / divisor}
}

func encodePoint(a []int) string {
	return fmt.Sprintf("%d-%d", a[0], a[1])
}

func decodePoint(str string) []int {
	split := strings.Split(str, "-")

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

func allAntennaPairs(antennas [][]int) [][][]int {
	var allPairs [][][]int

	for i, j := 0, 1; j < len(antennas); i, j = i+1, j+1 {
		for k := j; k < len(antennas); k++ {
			pair := [][]int{antennas[i], antennas[k]}
			allPairs = append(allPairs, pair)
		}
	}

	return allPairs
}

func Part1(rows *[]string) (int, error) {
	width := len((*rows)[0])
	height := len(*rows)
	antennasMap := findAntennas(rows)
	uniqueAntinodes := make(map[string]struct{})

	for key, antennas := range antennasMap {
		fmt.Printf("Analyzing antennas '%c'\n", key)
		fmt.Printf("Locations: %v\n", antennas)

		allPairs := allAntennaPairs(antennas)

		for _, pair := range allPairs {
			antinodes := findAntinodesForPoints(pair[0], pair[1], width-1, height-1)

			for _, node := range antinodes {
				encoded := encodePoint(node)
				uniqueAntinodes[encoded] = struct{}{}
			}
		}
	}

	return len(uniqueAntinodes), nil
}

func Part2(rows *[]string) (int, error) {
	width := len((*rows)[0])
	height := len(*rows)
	antennasMap := findAntennas(rows)
	uniqueAntinodes := make(map[string]struct{})

	for key, antennas := range antennasMap {
		fmt.Printf("Analyzing antennas '%c'\n", key)
		fmt.Printf("Locations: %v\n", antennas)

		allPairs := allAntennaPairs(antennas)
		fmt.Printf("all pairs: %v\n", allPairs)

		for _, pair := range allPairs {
			fmt.Printf("Pair: %v\n", pair)
			antinodes := findAntinodesForPointsPart2(pair[0], pair[1], width-1, height-1)
			fmt.Printf("Antinodes: %v\n", antinodes)

			for _, node := range antinodes {
				encoded := encodePoint(node)
				uniqueAntinodes[encoded] = struct{}{}
			}
		}
	}

	matrix := make([][]byte, 0, len(*rows))

	for _, row := range *rows {
		bytes := []byte(row)
		matrix = append(matrix, bytes)
	}

	for key := range uniqueAntinodes {
		decoded := decodePoint(key)
		matrix[decoded[1]][decoded[0]] = '#'
	}

	for _, row := range matrix {
		fmt.Printf("%s\n", string(row))
	}

	return len(uniqueAntinodes), nil
}
