package day_04

import (
	"fmt"
	"strings"
)

func reverseString(str string) string {
	bytes := []byte(str)

	i := 0
	j := len(bytes) - 1
	for ; i < j; i, j = i+1, j-1 {
		buf := bytes[i]
		bytes[i] = bytes[j]
		bytes[j] = buf
	}

	return string(bytes)
}

func reverseAll(data []string) []string {
	result := make([]string, len(data))

	for i, str := range data {
		result[i] = reverseString(str)
	}

	return result
}

func extendDataRight(data []string) []string {
	height := len(data)
	widthIncrease := height - 1

	result := make([]string, height)

	for i, row := range data {
		extendedRow := fmt.Sprintf("%s%s%s", createBlankString(i), row, createBlankString(widthIncrease-i))
		result[i] = extendedRow
	}

	return result
}

func extendDataLeft(data []string) []string {
	height := len(data)
	widthIncrease := height - 1

	result := make([]string, height)

	for i, row := range data {
		extendedRow := fmt.Sprintf("%s%s%s", createBlankString(widthIncrease-i), row, createBlankString(i))
		result[i] = extendedRow
	}

	return result
}

func createBlankString(lenght int) string {
	bytes := make([]byte, lenght)

	for i := range bytes {
		bytes[i] = '.'
	}

	return string(bytes)
}

func getColumns(data []string) []string {
	width := len(data[0])
	height := len(data)
	columns := make([]string, width)

	for x := 0; x < width; x++ {
		bytes := make([]byte, height)

		for y := 0; y < height; y++ {
			bytes[y] = data[y][x]
		}

		columns[x] = string(bytes)
	}

	return columns
}

func Part1(rows *[]string) (int, error) {
	rowsReversed := reverseAll(*rows)
	columns := getColumns(*rows)
	columsReversed := reverseAll(columns)
	principalDiagonals := getColumns(extendDataLeft(*rows))
	principalDiagonalsReversed := reverseAll(principalDiagonals)
	counterDiagonals := getColumns(extendDataRight(*rows))
	counterDiagonalsReversed := reverseAll(counterDiagonals)

	var all []string
	all = append(all, *rows...)
	all = append(all, rowsReversed...)
	all = append(all, columns...)
	all = append(all, columsReversed...)
	all = append(all, principalDiagonals...)
	all = append(all, principalDiagonalsReversed...)
	all = append(all, counterDiagonals...)
	all = append(all, counterDiagonalsReversed...)

	count := 0

	for _, str := range all {
		count += strings.Count(str, "XMAS")
	}

	return count, nil
}

func Part2(input *[]string) (int, error) {

	return -1, fmt.Errorf("not implemented")
}
