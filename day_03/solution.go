package day_03

import (
	"fmt"
	"regexp"
	"strconv"
)

const mulInstructionPattern string = `mul\(\d{1,3},\d{1,3}\)`
const numbersPattern = `\d{1,3}`

var mulIntructionRegexp = regexp.MustCompile(mulInstructionPattern)
var numbersRegexp = regexp.MustCompile(numbersPattern)

func processMulInstruction(instruction *string) (int, error) {
	numbers := numbersRegexp.FindAllString(*instruction, -1)

	if len(numbers) != 2 {
		return -1, fmt.Errorf("Expected two numbers in mul instruction")
	}

	a, err := strconv.Atoi(numbers[0])

	if err != nil {
		return -1, err
	}

	b, err := strconv.Atoi(numbers[1])

	if err != nil {
		return -1, err
	}

	return a * b, nil

}

func Part1(input *[]string) (int, error) {
	if len(*input) < 1 {
		return -1, fmt.Errorf("At least one line of input date expected")
	}

	sum := 0

	for _, memory := range *input {
		mulInstructions := mulIntructionRegexp.FindAllString(memory, -1)

		for _, instruction := range mulInstructions {
			result, err := processMulInstruction(&instruction)

			if err != nil {
				return -1, err
			}

			sum += result
		}
	}

	return sum, nil
}

func Part2(input *[]string) (int, error) {
	return -1, fmt.Errorf("Not implemented")
}
