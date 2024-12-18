package day_18

import (
	"container/heap"
	"fmt"
	"strconv"
	"strings"
)

func decodeVector(str string) []int {
	split := strings.Split(str, ",")
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

func encodeVector(vect []int) string {
	return fmt.Sprintf("%d,%d", vect[0], vect[1])
}

func parseInput(input *[]string) [][]int {
	result := make([][]int, 0, len(*input))

	for _, row := range *input {
		vector := decodeVector(row)
		result = append(result, vector)
	}

	return result
}

func printMemory(memory *[][]byte) {
	for _, row := range *memory {
		fmt.Println(string(row))
	}
}

const gridSize = 71

func Part1(input *[]string) (int, error) {
	corruptedFields := parseInput(input)
	memoryMap := make([][]byte, gridSize)

	for y := 0; y < gridSize; y++ {
		row := make([]byte, gridSize)
		memoryMap[y] = row
		for x := 0; x < gridSize; x++ {
			memoryMap[y][x] = '.'
		}
	}

	for i := 0; i < 1024 && i < len(corruptedFields); i++ {
		field := corruptedFields[i]
		memoryMap[field[1]][field[0]] = '#'
	}

	printMemory(&memoryMap)

	h := &MinHeap{}
	heap.Init(h)
	distances := make(map[string]int)
	maxInt := int(^uint(0) >> 1)

	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			if memoryMap[y][x] == '.' {
				point := Point{
					X:        x,
					Y:        y,
					Distance: maxInt,
				}
				if x == 0 && y == 0 {
					point.Distance = 0
				}

				distances[encodeVector([]int{x, y})] = point.Distance
				heap.Push(h, point)
			}
		}
	}

	distance := 0

	for current := heap.Pop(h).(Point); h.Len() > 0; current = heap.Pop(h).(Point) {
		if current.X == gridSize-1 && current.Y == gridSize-1 {
			distance = current.Distance
			break
		}

		neighbors := [][]int{
			{current.X + 1, current.Y},
			{current.X, current.Y + 1},
			{current.X - 1, current.Y},
			{current.X, current.Y - 1},
		}

		for _, n := range neighbors {
			if n[0] < 0 || n[1] < 0 || n[0] >= gridSize || n[1] >= gridSize {
				continue
			}

			if memoryMap[n[1]][n[0]] == '#' {
				continue
			}

			currentDistance, hasDistance := distances[encodeVector([]int{current.X, current.Y})]
			if !hasDistance {
				panic("Should not happen")
			}
			neighborShortestDistance, hasDistance := distances[encodeVector(n)]
			if !hasDistance {
				panic("Should not happen")
			}

			newNeighborDistance := currentDistance + 1
			if newNeighborDistance < neighborShortestDistance {
				distances[encodeVector(n)] = newNeighborDistance
				h.UpdateDistance(n[0], n[1], newNeighborDistance)
			}

		}

	}

	return distance, nil
}

func Part2(input *[]string) (int, error) {
	corruptedFields := parseInput(input)
	memoryMap := make([][]byte, gridSize)
	maxInt := int(^uint(0) >> 1)

	for y := 0; y < gridSize; y++ {
		row := make([]byte, gridSize)
		memoryMap[y] = row
		for x := 0; x < gridSize; x++ {
			memoryMap[y][x] = '.'
		}
	}

	pathCutOff := []int{-1, -1}
	for numberOfCorruptedBytes := 1025; numberOfCorruptedBytes < len(corruptedFields); numberOfCorruptedBytes++ {
		for i := 0; i < numberOfCorruptedBytes && i < len(corruptedFields); i++ {
			field := corruptedFields[i]
			memoryMap[field[1]][field[0]] = '#'
		}

		fmt.Printf("Byte %d, value: %v\n", numberOfCorruptedBytes-1, corruptedFields[numberOfCorruptedBytes-1])
		printMemory(&memoryMap)

		h := &MinHeap{}
		heap.Init(h)
		distances := make(map[string]int)

		for y := 0; y < gridSize; y++ {
			for x := 0; x < gridSize; x++ {
				if memoryMap[y][x] == '.' {
					point := Point{
						X:        x,
						Y:        y,
						Distance: maxInt,
					}
					if x == 0 && y == 0 {
						point.Distance = 0
					}

					distances[encodeVector([]int{x, y})] = point.Distance
					heap.Push(h, point)
				}
			}
		}

		distance := -1
		fmt.Printf("Heap len %d\n", h.Len())

		for current := heap.Pop(h).(Point); h.Len() > 0; current = heap.Pop(h).(Point) {
			if current.Distance == maxInt {
				break
			}

			if current.X == gridSize-1 && current.Y == gridSize-1 {
				distance = current.Distance
				break
			}

			neighbors := [][]int{
				{current.X + 1, current.Y},
				{current.X, current.Y + 1},
				{current.X - 1, current.Y},
				{current.X, current.Y - 1},
			}

			for _, n := range neighbors {
				if n[0] < 0 || n[1] < 0 || n[0] >= gridSize || n[1] >= gridSize {
					continue
				}

				if memoryMap[n[1]][n[0]] == '#' {
					continue
				}

				currentDistance, hasDistance := distances[encodeVector([]int{current.X, current.Y})]
				if !hasDistance {
					panic("Should not happen")
				}
				neighborShortestDistance, hasDistance := distances[encodeVector(n)]
				if !hasDistance {
					panic("Should not happen")
				}

				newNeighborDistance := currentDistance + 1
				if newNeighborDistance < neighborShortestDistance {
					distances[encodeVector(n)] = newNeighborDistance
					h.UpdateDistance(n[0], n[1], newNeighborDistance)
				}
			}
		}

		if distance == -1 {
			pathCutOff = corruptedFields[numberOfCorruptedBytes-1]
			break
		}
		fmt.Printf("Dist %d\n", distance)
		fmt.Println("----")
	}

	fmt.Printf("Can't find a path at byte with coords: %v\n", pathCutOff)

	return 0, nil
}

type Point struct {
	X        int
	Y        int
	Distance int
}

type MinHeap []Point

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Less(i, j int) bool {
	return h[i].Distance < h[j].Distance
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(Point))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

func (h *MinHeap) UpdateDistance(x, y, newDistance int) {
	for i := range *h {
		if (*h)[i].X == x && (*h)[i].Y == y {
			(*h)[i].Distance = newDistance
			heap.Fix(h, i)
			return
		}
	}
}
