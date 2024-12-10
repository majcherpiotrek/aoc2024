package day_09

import (
	"fmt"
	"strconv"
)

type DiscMap struct {
	DiscMap    []int
	FreeSpaces [][]int
	Files      [][]int
}

func parseDiscMap(str string) (DiscMap, error) {
	bytes := []byte(str)
	nextFileId := 0
	var discMap []int
	var freeSpaces [][]int
	var files [][]int

	for i, b := range bytes {
		isFile := i%2 == 0
		length, err := strconv.Atoi(string(b))

		if err != nil {
			return DiscMap{}, err
		}

		fileOrSpace := make([]int, 3)

		fileOrSpace[1] = length

		if isFile {
			fileOrSpace[2] = nextFileId
		}

		for j := 0; j < length; j++ {
			if isFile {
				discMap = append(discMap, nextFileId)
			} else {
				discMap = append(discMap, -1)
			}

			if j == 0 {
				fileOrSpace[0] = len(discMap) - 1
			}
		}

		if isFile {
			nextFileId++
			files = append(files, fileOrSpace)
		} else {
			freeSpaces = append(freeSpaces, fileOrSpace)
		}

	}

	return DiscMap{
		DiscMap:    discMap,
		FreeSpaces: freeSpaces,
		Files:      files,
	}, nil
}

func organiseDisc(discMap DiscMap) []int {
	result := make([]int, 0, len(discMap.DiscMap))
	result = append(result, discMap.DiscMap...)

	i := len(discMap.Files) - 1
	j := 0
	file := discMap.Files[i]
	freeSpace := discMap.FreeSpaces[j]

	for file[0] > freeSpace[0] {
		fmt.Printf("File: %v\n", file)
		fmt.Printf("Free space: %v\n", freeSpace)

		freeSpaceIndex := freeSpace[0]
		freeSpaceSize := freeSpace[1]
		fileIndex := file[0]
		fileSize := file[1]

		if fileSize <= freeSpaceSize {
			toMove := discMap.DiscMap[fileIndex : fileIndex+fileSize]
			for x, el := range toMove {
				result[freeSpaceIndex+x] = el
				result[fileIndex+x] = -1
			}

			freeSpace = []int{freeSpaceIndex + fileSize, freeSpaceSize - fileSize, -1}
			i--
			file = discMap.Files[i]
		} else {
			fileEnd := fileIndex + fileSize
			copyStart := fileEnd - freeSpaceSize
			toMove := discMap.DiscMap[copyStart:fileEnd]
			for x, el := range toMove {
				result[freeSpaceIndex+x] = el
				result[copyStart+x] = -1
			}

			file = []int{fileIndex, fileSize - freeSpaceSize, file[2]}
			j++
			freeSpace = discMap.FreeSpaces[j]
		}
	}

	return result
}

func calculateChecksum(organised []int) int {
	sum := 0

	for i, fileId := range organised {

		if fileId == -1 {
			break
		}

		sum = sum + (i * fileId)
	}

	return sum
}

func calculateChecksum2(organised []int) int {
	sum := 0

	for i, fileId := range organised {

		if fileId == -1 {
			continue
		}

		sum = sum + (i * fileId)
	}

	return sum
}

func Part1(rows *[]string) (int, error) {
	discMap, err := parseDiscMap((*rows)[0])

	if err != nil {
		return -1, err
	}

	organised := organiseDisc(discMap)

	fmt.Printf("Organised: %v\n", organised)

	return calculateChecksum(organised), nil
}

func Part2(rows *[]string) (int, error) {
	discMap, err := parseDiscMap((*rows)[0])

	if err != nil {
		return -1, err
	}

	numOfFiles := len(discMap.Files)

	for f := numOfFiles - 1; f >= 0; f-- {
		file := discMap.Files[f]
		fmt.Printf("File: %v\n", file)
		fileIndex := file[0]
		fileSize := file[1]
		fileID := file[2]

		for _, space := range discMap.FreeSpaces {
			spaceIndex := space[0]
			spaceSize := space[1]

			if spaceIndex != -1 && spaceSize >= fileSize && spaceIndex < fileIndex {
				for i := 0; i < fileSize; i++ {
					discMap.DiscMap[spaceIndex+i] = fileID
					discMap.DiscMap[fileIndex+i] = -1
				}

				diff := spaceSize - fileSize
				if diff > 0 {
					space[0] = spaceIndex + fileSize
					space[1] = spaceSize - fileSize
				} else {
					space[0] = -1
				}
				break
			}

		}
	}

	return calculateChecksum2(discMap.DiscMap), nil

}
