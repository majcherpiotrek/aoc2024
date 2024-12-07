package day_07

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Operation int

const (
	Add = iota
	Mul
	Concat
)

type CalibrationEquation struct {
	TestValue int
	Numbers   []int
}

func parseInput(rows *[]string) ([]CalibrationEquation, error) {
	var equations []CalibrationEquation

	for _, row := range *rows {
		split := strings.Split(row, ": ")

		if len(split) != 2 {
			return []CalibrationEquation{}, fmt.Errorf("Unexpected input row %s", row)
		}

		testValue, err := strconv.Atoi(split[0])

		if err != nil {
			return []CalibrationEquation{}, err
		}

		split = strings.Fields(split[1])

		numbers := make([]int, 0, len(split))

		for _, numStr := range split {
			num, err := strconv.Atoi(numStr)

			if err != nil {
				return []CalibrationEquation{}, err
			}

			numbers = append(numbers, num)
		}

		equations = append(equations, CalibrationEquation{TestValue: testValue, Numbers: numbers})
	}

	return equations, nil
}

func generatePossibilities(n int) [][]Operation {
	if n == 0 {
		return make([][]Operation, 0)
	}

	max := int(math.Pow(2, float64(n)))

	var allOperations [][]Operation

	for i := 0; i < max; i++ {
		var operations []Operation
		binary := fmt.Sprintf("%0*b", n, i)

		for _, str := range binary {
			b, err := strconv.Atoi(string(str))
			if err != nil {
				panic("test")
			}
			op := Operation(b)
			operations = append(operations, op)
		}

		allOperations = append(allOperations, operations)
	}

	return allOperations
}

func Part1(rows *[]string) (int, error) {
	equations, err := parseInput(rows)

	if err != nil {
		return -1, err
	}

	sumOfValidTestValues := 0

	possibilitiesCache := make(map[int][][]Operation)

	for _, equation := range equations {
		n := len(equation.Numbers) - 1
		possibilities, alreadyComputed := possibilitiesCache[n]

		if !alreadyComputed {
			possibilities = generatePossibilities(n)
			possibilitiesCache[n] = possibilities
		}

		isValid := false
		for _, possibility := range possibilities {
			if len(equation.Numbers) == 0 {
				break
			}

			if len(equation.Numbers) == 1 {
				if equation.Numbers[0] == equation.TestValue {
					isValid = true
					break
				} else {

					isValid = false
					break
				}
			}

			sum := equation.Numbers[0]

			if sum > equation.TestValue {
				isValid = false
				break
			}

			for i := 1; i < len(equation.Numbers); i++ {
				op := possibility[i-1]

				switch op {
				case Add:
					sum = sum + equation.Numbers[i]
				case Mul:
					sum = sum * equation.Numbers[i]
				}

				if sum > equation.TestValue {
					isValid = false
					break
				}
			}

			if sum == equation.TestValue {
				isValid = true
			}

			if isValid {
				break
			}
		}

		if isValid {
			sumOfValidTestValues += equation.TestValue
		}
	}

	return sumOfValidTestValues, nil
}

func checkIsValid(equation CalibrationEquation) bool {
	if len(equation.Numbers) == 0 {
		return false
	}

	if len(equation.Numbers) == 1 {
		if equation.TestValue == equation.Numbers[0] {
			return true
		}
		return false
	}

	if equation.Numbers[0] > equation.TestValue {
		return false
	}

	a := equation.Numbers[0]
	b := equation.Numbers[1]
	rest := equation.Numbers[2:]

	// Try to add
	sum := a + b
	newNumbers := []int{sum}
	newNumbers = append(newNumbers, rest...)

	if sum == equation.TestValue && len(rest) == 0 {
		return true
	}

	if len(equation.Numbers) > 2 {
		isValid := checkIsValid(CalibrationEquation{TestValue: equation.TestValue, Numbers: newNumbers})
		if isValid {
			return true
		}
	}

	// Try to multiply
	sum = a * b
	newNumbers = []int{sum}
	newNumbers = append(newNumbers, rest...)
	if sum == equation.TestValue && len(rest) == 0 {
		return true
	}
	if len(equation.Numbers) > 2 {
		isValid := checkIsValid(CalibrationEquation{TestValue: equation.TestValue, Numbers: newNumbers})
		if isValid {
			return true
		}
	}

	// Try to concat
	sumStr := fmt.Sprintf("%d%d", a, b)
	sum, err := strconv.Atoi(sumStr)

	if err != nil {
		panic("Unexpected error")
	}

	newNumbers = []int{sum}
	newNumbers = append(newNumbers, rest...)
	if sum == equation.TestValue && len(rest) == 0 {
		return true
	}

	if len(equation.Numbers) > 2 {
		isValid := checkIsValid(CalibrationEquation{TestValue: equation.TestValue, Numbers: newNumbers})
		if isValid {
			return true
		}
	}

	return false

}

func Part2(rows *[]string) (int, error) {
	equations, err := parseInput(rows)

	if err != nil {
		return -1, err
	}

	sumOfValidTestValues := 0

	for _, equation := range equations {
		isValid := checkIsValid(equation)

		if isValid {
			sumOfValidTestValues += equation.TestValue
		}

	}

	return sumOfValidTestValues, nil

}
