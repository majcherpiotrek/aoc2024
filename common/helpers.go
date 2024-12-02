package common

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseInput(input *[]string) ([][]int, error) {
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
