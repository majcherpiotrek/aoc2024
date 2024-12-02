package day_01

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func parseInput(input *[]string) ([][]int, error) {
	var column_1 []int
	var column_2 []int

	for i, line := range *input {
		words := strings.Fields(line)

		if len(words) != 2 {
			return [][]int{}, fmt.Errorf("Invalid input at line %d - expected exactly two columns", i)
		}

		value_1, err := strconv.Atoi(words[0])

		if err != nil {
			return [][]int{}, fmt.Errorf("Invalid value at line %d, column 0: %s - Integer expected", i, words[0])
		}

		value_2, err := strconv.Atoi(words[1])

		if err != nil {
			return [][]int{}, fmt.Errorf("Invalid value at line %d, column 1: %s - Integer expected", i, words[1])
		}

		column_1 = append(column_1, value_1)
		column_2 = append(column_2, value_2)
	}

	return append([][]int{}, column_1, column_2), nil

}

func Part1(input *[]string) (int, error) {
	data, err := parseInput(input)

	if err != nil {
		return -1, err
	}

	column_1 := data[0]
	column_2 := data[1]

	slices.Sort(column_1)
	slices.Sort(column_2)

	totalDistance := 0

	for i, value_1 := range column_1 {

		distance := value_1 - column_2[i]

		if distance < 0 {
			totalDistance += -distance
		} else {
			totalDistance += distance
		}
	}

	return totalDistance, nil
}

func calculateOccurrences(secondList *[]int) map[int]int {
	occurrencesMap := make(map[int]int)

	for _, el := range *secondList {
		sum, exists := occurrencesMap[el]

		if exists {
			occurrencesMap[el] = sum + 1
		} else {
			occurrencesMap[el] = 1
		}
	}

	return occurrencesMap
}

func Part2(input *[]string) (int, error) {
	data, err := parseInput(input)

	if err != nil {
		return -1, err
	}

	column_1 := data[0]
	column_2 := data[1]

	occurrencesMap := calculateOccurrences(&column_2)

	sum := 0

	for _, value := range column_1 {
		occurrences, hasAnyOccurrences := occurrencesMap[value]

		multiplied := 0

		if hasAnyOccurrences {
			multiplied = value * occurrences
		}

		sum += multiplied
	}

	return sum, nil
}
