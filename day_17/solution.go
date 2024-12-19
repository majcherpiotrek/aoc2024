package day_17

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type registers struct {
	A int
	B int
	C int
}

type processor struct {
	registers  registers
	Out        []int
	in         []int
	insPointer int
}

type Instruction int

const (
	adv = iota // division A / (combo operand ^ 2) => A
	bxl        // bitwise XOR : B ^ literal operand => B
	bst        // combo operand % 8 => B
	jnz        // if A == 0 nothing else insPointer = literal operand
	bxc        // B ^ C => B, ignores operand
	out        // combo % A => Out
	bdv        // like adv, but stores in B
	cdv        // like adv, but stores in C
)

func NewProcessor(reg registers, in []int) processor {
	return processor{
		registers:  reg,
		Out:        make([]int, 0),
		in:         in,
		insPointer: 0,
	}
}

func (p *processor) readCurrentInstruction() (Instruction, int, error) {
	if len(p.in) == 0 || p.insPointer >= len(p.in)-1 {
		return -1, -1, fmt.Errorf("Instruction pointer out of bounds")
	}

	if p.insPointer%2 != 0 {
		return -1, -1, fmt.Errorf("Instruction pointer pointing to operand insted of instruction")
	}

	instruction := p.in[p.insPointer]
	operand := p.in[p.insPointer+1]

	return Instruction(instruction), operand, nil
}

func (p *processor) readComboOperand(operand int) (int, error) {
	if operand <= 3 {
		return operand, nil
	}

	operandValue := operand

	switch operand {
	case 4:
		operandValue = p.registers.A
	case 5:
		operandValue = p.registers.B
	case 6:
		operandValue = p.registers.C
	default:
		return -1, fmt.Errorf("Invalid operand: %d", operand)
	}

	return operandValue, nil
}

func (p *processor) division(operandValue int, targetRegister *int) {
	fmt.Printf("Operand value: %d, target register: %p, targetRegisterValue: %d\n", operandValue, targetRegister, *targetRegister)
	*targetRegister = p.registers.A / int(math.Pow(2, float64(operandValue)))
}

func (p *processor) tick() error {
	instruction, operand, err := p.readCurrentInstruction()
	fmt.Printf("Process instruction - pointer: %d: instruciton: %d, operand: %d\n", p.insPointer, instruction, operand)

	if err != nil {
		return err
	}

	if instruction == adv || instruction == bdv || instruction == cdv {
		operandValue, err := p.readComboOperand(operand)
		if err != nil {
			return err
		}

		var targetRegister *int = nil

		switch instruction {
		case adv:
			targetRegister = &p.registers.A
		case bdv:
			targetRegister = &p.registers.B
		case cdv:
			targetRegister = &p.registers.C
		}
		fmt.Println("Running division")
		p.division(operandValue, targetRegister)
		p.insPointer += 2
		return nil
	}

	switch instruction {
	case bxl:
		p.registers.B = p.registers.B ^ operand
		p.insPointer += 2
	case bst:
		operandValue, err := p.readComboOperand(operand)
		if err != nil {
			return err
		}
		p.registers.B = operandValue % 8
		p.insPointer += 2
	case jnz:
		if p.registers.A != 0 {
			p.insPointer = operand
			return nil
		} else {
			p.insPointer += 2
		}
	case bxc:
		p.registers.B = p.registers.B ^ p.registers.C
		p.insPointer += 2
	case out:
		operandValue, err := p.readComboOperand(operand)
		result := operandValue % 8
		fmt.Printf("OUT %d mod 8 = %d\n", operandValue, result)
		if err != nil {
			return err
		}
		p.Out = append(p.Out, result)
		p.insPointer += 2
	}

	return nil
}

func (p *processor) Run() string {
	fmt.Printf("Start program. Pointer: %d, Registers: %v, Program: %v\n", p.insPointer, p.registers, p.in)
	for err := p.tick(); err == nil; err = p.tick() {
		fmt.Printf("STATE: registers: %v, pointer: %d, out: %v\n", p.registers, p.insPointer, p.Out)
	}

	result := ""
	for _, n := range p.Out {
		if len(result) == 0 {
			result = fmt.Sprintf("%d", n)
		} else {
			result = fmt.Sprintf("%s,%d", result, n)
		}
	}

	return result
}

func Part1(input *[]string) (int, error) {
	a, b, c := -1, -1, -1
	var p processor

	for i := 0; i < len(*input); i++ {
		row := (*input)[i]

		if len(row) == 0 {
			continue
		}

		split := strings.Split(row, ": ")
		if len(split) != 2 {
			return -1, fmt.Errorf("Invalid input row: %s", row)
		}

		if i < 3 {
			n, err := strconv.Atoi(split[1])
			if err != nil {
				return -1, err
			}
			if i == 0 {
				a = n
			}
			if i == 1 {
				b = n
			}
			if i == 2 {
				c = n
			}
			continue
		}

		numstr := strings.Split(split[1], ",")
		program := make([]int, 0, len(numstr))

		for _, s := range numstr {
			n, err := strconv.Atoi(s)
			if err != nil {
				return -1, err
			}
			program = append(program, n)
		}

		p = NewProcessor(registers{
			A: a,
			B: b,
			C: c,
		}, program)
	}

	result := p.Run()

	fmt.Println(result)

	return 0, nil

}

func Part2(input *[]string) (int, error) {
	return -1, fmt.Errorf("not implemented")
}
