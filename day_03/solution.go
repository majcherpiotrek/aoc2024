package day_03

import (
	"fmt"
	"regexp"
	"strconv"
)

const mulInstructionPattern string = `mul\(\d{1,3},\d{1,3}\)`
const doInstructionPattern = `do\(\)`
const dontInstructionPattern = `don\'t\(\)`
const numbersPattern = `\d{1,3}`

var mulInstructionRegexp = regexp.MustCompile(mulInstructionPattern)
var doInstructionRegexp = regexp.MustCompile(doInstructionPattern)
var dontInstructionRegexp = regexp.MustCompile(dontInstructionPattern)
var allInstructionsRegexp = regexp.MustCompile(fmt.Sprintf("(%s)|(%s)|(%s)", mulInstructionPattern, doInstructionPattern, dontInstructionPattern))
var numbersRegexp = regexp.MustCompile(numbersPattern)

type Instruction interface {
	isInstruction()
}

type Do struct{}

func (Do) isInstruction() {}
func (Do) String() string {
	return "do()"
}

type Dont struct{}

func (Dont) isInstruction() {}

func (Dont) String() string {
	return "don't()"
}

type Mul struct {
	A int
	B int
}

func (m Mul) String() string {
	return fmt.Sprintf("mul(%d, %d)", m.A, m.B)
}

func (Mul) isInstruction() {}

func parseInstruction(str string) (Instruction, error) {
	instruction := mulInstructionRegexp.FindString(str)

	if len(instruction) > 0 {
		numbers := numbersRegexp.FindAllString(instruction, -1)
		if len(numbers) != 2 {
			return nil, fmt.Errorf("Expected two numbers in mul instruction")
		}

		A, err := strconv.Atoi(numbers[0])
		if err != nil {
			return nil, err
		}

		B, err := strconv.Atoi(numbers[1])
		if err != nil {
			return nil, err
		}

		return Mul{A, B}, nil
	}

	instruction = doInstructionRegexp.FindString(str)

	if len(instruction) > 0 {
		return Do{}, nil
	}

	instruction = dontInstructionRegexp.FindString(str)

	if len(instruction) > 0 {
		return Dont{}, nil
	}

	return nil, fmt.Errorf("Unknown instruction: %s", str)
}

func parseInstructions(str string) ([]Instruction, error) {
	instructionStrings := allInstructionsRegexp.FindAllString(str, -1)

	numOfInstructions := len(instructionStrings)

	instructions := make([]Instruction, numOfInstructions)

	for i := 0; i < numOfInstructions; i++ {
		instruction, err := parseInstruction(instructionStrings[i])

		if err != nil {
			return []Instruction{}, err
		}

		instructions[i] = instruction
	}

	return instructions, nil
}

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

func processInstructions(instructions *[]Instruction) int {
	isMulProcessingEnabled := true
	sum := 0

	for _, instruction := range *instructions {
		fmt.Println(fmt.Sprintf("Instruction: %s, isMulEnabled: %t, sum: %d", instruction, isMulProcessingEnabled, sum))
		switch ins := instruction.(type) {
		case Do:
			isMulProcessingEnabled = true
		case Dont:
			isMulProcessingEnabled = false
		case Mul:
			if isMulProcessingEnabled {
				sum += ins.A * ins.B
			}
		}
	}

	return sum
}

func Part1(input *[]string) (int, error) {
	if len(*input) < 1 {
		return -1, fmt.Errorf("At least one line of input date expected")
	}

	sum := 0

	for _, memory := range *input {
		mulInstructions := mulInstructionRegexp.FindAllString(memory, -1)

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
	if len(*input) < 1 {
		return -1, fmt.Errorf("At least one line of input date expected")
	}

	var allInstructions []Instruction

	for _, memory := range *input {
		instructions, err := parseInstructions(memory)
		if err != nil {
			return -1, err
		}

		allInstructions = append(allInstructions, instructions...)
	}

	return processInstructions(&allInstructions), nil
}
