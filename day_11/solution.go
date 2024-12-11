package day_11

import (
	"fmt"
	"strconv"
	"strings"
)

func parseInput(input *[]string) ([]int, error) {
	if len(*input) != 1 {
		return []int{}, fmt.Errorf("Only one line of input expected")
	}

	fields := strings.Fields((*input)[0])

	result := make([]int, 0, len(fields))

	for _, field := range fields {
		num, err := strconv.Atoi(field)

		if err != nil {
			return []int{}, err
		}

		result = append(result, num)
	}

	return result, nil
}

type Stone struct {
	Num  int
	Prev *Stone
	Next *Stone
}

func (stone *Stone) changeOnBlink() {
	//fmt.Printf("Change stone: {Num: %d, Prev: %v, Next: %v}\n", stone.Num, stone.Prev, stone.Next)
	if stone.Num == 0 {
		stone.Num = 1
		return
	}

	numAsString := strconv.Itoa(stone.Num)
	numOfDigits := len(numAsString)

	if numOfDigits%2 == 0 {
		//fmt.Printf("Splitting digits\n")
		splitIndex := numOfDigits / 2

		firstNumDigits := numAsString[0:splitIndex]
		secondNumDigits := numAsString[splitIndex:]

		firstNum, err := strconv.Atoi(firstNumDigits)
		if err != nil {
			panic("Should not happen")
		}

		secondNum, err := strconv.Atoi(secondNumDigits)
		if err != nil {
			panic("Should not happen")
		}

		stone.Num = firstNum

		secondStone := Stone{
			Num:  secondNum,
			Prev: stone,
			Next: stone.Next,
		}

		stone.Next = &secondStone

		if secondStone.Next != nil {
			secondStone.Next.Prev = &secondStone
		}
		return
	}

	stone.Num = stone.Num * 2024
	return
}

func (stone *Stone) append(stoneToAdd Stone) {
	nextStone := stone

	for nextStone.Next != nil {
		nextStone = nextStone.Next
	}

	nextStone.Next = &stoneToAdd
	stoneToAdd.Prev = nextStone
}

func (stone *Stone) toString() string {
	nextStone := stone
	acc := ""

	for nextStone != nil {
		acc = fmt.Sprintf("%s{Num: %d, Addr: %p, Prev: %v, Next: %v}\n", acc, nextStone.Num, nextStone, nextStone.Prev, nextStone.Next)
		nextStone = nextStone.Next
	}
	return acc
}

func (stone *Stone) toSimpleString() string {
	acc := ""

	stone.forEach(func(s Stone) {
		acc = fmt.Sprintf("%s %d", acc, s.Num)
	})

	return acc
}

func (stone *Stone) forEach(fn func(Stone)) {
	nextStone := stone

	for nextStone != nil {
		fn(*nextStone)
		nextStone = nextStone.Next
	}
}

func Part1(rows *[]string) (int, error) {
	stones, err := parseInput(rows)

	if err != nil {
		return -1, err
	}

	var stonesLinkedList Stone

	for i, stone := range stones {
		if i == 0 {
			stonesLinkedList = Stone{
				Num:  stone,
				Next: nil,
				Prev: nil,
			}
		} else {
			stonesLinkedList.append(Stone{
				Num:  stone,
				Next: nil,
				Prev: nil,
			})
		}
	}

	for i := 0; i < 25; i++ {
		fmt.Printf("Blink %d\n", i+1)
		//fmt.Printf("Before %s\n", stonesLinkedList.toString())

		nextStone := &stonesLinkedList

		for nextStone != nil {
			next := nextStone.Next
			nextStone.changeOnBlink()
			nextStone = next
		}

		//fmt.Printf("After %d blink: %s\n", i+1, stonesLinkedList.toSimpleString())
	}

	length := 0

	stonesLinkedList.forEach(func(s Stone) {
		length++
	})

	return length, nil
}

func Part2(rows *[]string) (int, error) {
	stones, err := parseInput(rows)

	if err != nil {
		return -1, err
	}

	stonesMap := make(map[int]int)

	for _, s := range stones {
		_, alreadyAdded := stonesMap[s]

		if alreadyAdded {
			stonesMap[s]++
		} else {
			stonesMap[s] = 1
		}
	}

	currentMap := &stonesMap
	//fmt.Printf("%v\n\n", *currentMap)

	for i := 0; i < 75; i++ {
		fmt.Printf("Blink %d\n", i+1)
		nextMap := make(map[int]int)

		for stone, count := range *currentMap {
			//fmt.Printf("{ stone: %d, count: %d }\n", stone, count)
			if count == 0 {
				//fmt.Println("No elements")
				continue
			}

			if stone == 0 {
				//fmt.Println("Changing 0 to 1")
				_, alreadyPresent := nextMap[1]
				if alreadyPresent {
					nextMap[1] += count
				} else {
					nextMap[1] = count
				}
				continue
			}

			numAsString := strconv.Itoa(stone)
			numOfDigits := len(numAsString)

			if numOfDigits%2 == 0 {
				//fmt.Println("Splitting")
				splitIndex := numOfDigits / 2

				firstNumDigits := numAsString[0:splitIndex]
				secondNumDigits := numAsString[splitIndex:]

				firstNum, err := strconv.Atoi(firstNumDigits)
				if err != nil {
					panic("Should not happen")
				}

				secondNum, err := strconv.Atoi(secondNumDigits)
				if err != nil {
					panic("Should not happen")
				}

				//fmt.Printf("first: %d, second: %d", firstNum, secondNum)

				_, alreadyPresent := nextMap[firstNum]
				if alreadyPresent {
					nextMap[firstNum] += count
				} else {
					nextMap[firstNum] = count
				}

				_, alreadyPresent = nextMap[secondNum]
				if alreadyPresent {
					nextMap[secondNum] += count
				} else {
					nextMap[secondNum] = count
				}
				continue
			}

			//fmt.Println("Multiplying by 2024")
			multiplied := stone * 2024
			_, alreadyPresent := nextMap[multiplied]

			if alreadyPresent {
				nextMap[multiplied] += count
			} else {
				nextMap[multiplied] = count
			}
		}

		currentMap = &nextMap
		//fmt.Printf("%v\n\n", *currentMap)
	}

	length := 0

	for _, count := range *currentMap {
		length += count
	}

	return length, nil
}
