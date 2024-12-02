package day_01

import (
	"aoc2024/common"
	"slices"
)

func Part1(input *[]string) (int, error) {
	data, err := common.ParseInput(input)

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
	data, err := common.ParseInput(input)

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
