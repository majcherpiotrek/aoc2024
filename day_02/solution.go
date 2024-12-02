package day_02

import (
	"fmt"
	"strconv"
	"strings"
)

func parseInput(input *[]string) ([][]int, error) {
	var result [][]int

	for rowNumber, line := range *input {
		words := strings.Fields(line)
		row := make([]int, len(words))

		for columnNumber, word := range words {
			num, err := strconv.Atoi(word)

			if err != nil {
				return [][]int{}, fmt.Errorf("Invalid value at line %d, column %d: %s - Integer expected", rowNumber, columnNumber, word)
			}

			row[columnNumber] = num
		}

		result = append(result, row)
	}

	return result, nil

}

var safeLevelIncreaseStep = map[int]bool{
	1: true,
	2: true,
	3: true,
}

var safeLevelDecreaseStep = map[int]bool{
	-1: true,
	-2: true,
	-3: true,
}

func checkIsStepSafe(step int, isIncreasing bool) bool {
	if isIncreasing {
		_, isAllowed := safeLevelIncreaseStep[step]
		return isAllowed
	} else {
		_, isAllowed := safeLevelDecreaseStep[step]
		return isAllowed
	}
}

func isReportSafe(report *[]int) bool {
	isSafe := len(*report) > 1
	var isIncreasing bool

	reportLength := len(*report)

	for i := 1; i < reportLength && isSafe; i++ {
		step := (*report)[i] - (*report)[i-1]

		if i == 1 {
			isIncreasing = step > 0
		}

		isSafe = checkIsStepSafe(step, isIncreasing)
	}

	return isSafe
}

func Part1(input *[]string) (int, error) {
	data, err := parseInput(input)

	if err != nil {
		return -1, err
	}

	numOfSafeReports := 0

	for _, report := range data {
		if isReportSafe(&report) {
			numOfSafeReports += 1
		}

	}

	return numOfSafeReports, nil
}

func splitPossibilities(report *[]int) [][]int {
	reportLength := len(*report)
	possibilities := make([][]int, 0, reportLength)

	for i := 0; i < reportLength; i++ {
		start := append([]int{}, (*report)[0:i]...)
		end := append([]int{}, (*report)[i+1:]...)
		combined := append(start, end...)
		possibilities = append(possibilities, combined)
	}

	return possibilities
}

func Part2(input *[]string) (int, error) {
	data, err := parseInput(input)

	if err != nil {
		return -1, err
	}

	numOfSafeReports := 0

	for _, report := range data {
		if isReportSafe(&report) {
			numOfSafeReports += 1
		} else {
			otherPossibilities := splitPossibilities(&report)
			for _, sub := range otherPossibilities {
				if isReportSafe(&sub) {
					numOfSafeReports += 1
					break
				}

			}
		}

	}

	return numOfSafeReports, nil
}
