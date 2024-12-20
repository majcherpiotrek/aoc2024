package day_20

import (
	"fmt"
)

func parseRacetrack(input *[]string) [][]byte {
	res := make([][]byte, 0, len(*input))
	for _, row := range *input {
		bytes := []byte(row)
		res = append(res, bytes)
	}

	return res
}

func addVectors(a []int, b []int) []int {
	return []int{a[0] + b[0], a[1] + b[1]}
}

func rotatePlus90(direction []int) []int {
	if direction[0] == 0 && direction[1] == 1 {
		return []int{1, 0}
	}
	if direction[0] == 1 && direction[1] == 0 {
		return []int{0, -1}
	}
	if direction[0] == 0 && direction[1] == -1 {
		return []int{-1, 0}
	}
	if direction[0] == -1 && direction[1] == 0 {
		return []int{0, 1}
	}

	panic("incorrect direction vector")
}

func rotateMinus90(direction []int) []int {
	rotated := rotatePlus90(direction)
	return []int{rotated[0] * -1, rotated[1] * -1}
}

func raceWithCheats(racetrack *[][]byte, start []int, currentDirection []int, currentTime int, cheatField []int, cheats *map[string]int, cache *map[string]struct{}) {
	cacheKey := fmt.Sprintf("%v%v%v", start, currentDirection, cheatField)
	_, alreadyVisited := (*cache)[cacheKey]
	if alreadyVisited {
		return
	}

	(*cache)[cacheKey] = struct{}{}

	width := len((*racetrack)[0])
	height := len(*racetrack)

	if start[0] < 0 || start[0] >= width || start[1] < 0 || start[1] >= height {
		return
	}

	field := (*racetrack)[start[1]][start[0]]

	if field == 'E' {
		if len(cheatField) > 0 {
			(*cheats)[fmt.Sprintf("%d,%d", cheatField[0], cheatField[1])] = currentTime
		}
		return
	}

	chField := cheatField

	if field == '#' {
		if len(cheatField) > 0 {
			return
		} else {
			chField = start
		}
	}

	next1 := addVectors(start, currentDirection)
	next2 := addVectors(start, rotatePlus90(currentDirection))
	next3 := addVectors(start, rotateMinus90(currentDirection))

	raceWithCheats(racetrack, next1, currentDirection, currentTime+1, chField, cheats, cache)
	raceWithCheats(racetrack, next2, rotatePlus90(currentDirection), currentTime+1, chField, cheats, cache)
	raceWithCheats(racetrack, next3, rotateMinus90(currentDirection), currentTime+1, chField, cheats, cache)
}

func Part1(input *[]string) (int, error) {
	racetrack := parseRacetrack(input)

	start := []int{-1, -1}

	for y := range racetrack {
		for x := range racetrack[y] {
			if racetrack[y][x] == 'S' {
				start = []int{x, y}
				break
			}
		}
	}

	cheatTimes := make(map[string]int)
	cache := make(map[string]struct{})

	raceWithCheats(&racetrack, start, []int{0, -1}, 0, []int{}, &cheatTimes, &cache)

	inverseCheatTimesMap := make(map[int][]string)

	for cheatField, raceTime := range cheatTimes {
		if raceTime < 84 {
			timeSaved := 84 - raceTime
			cheatFields, hasCheatFields := inverseCheatTimesMap[timeSaved]
			if hasCheatFields {
				inverseCheatTimesMap[timeSaved] = append(cheatFields, cheatField)
			} else {
				inverseCheatTimesMap[timeSaved] = []string{cheatField}
			}

		}
	}

	for time, cheats := range inverseCheatTimesMap {
		fmt.Printf("%d: %d - %v\n", time, len(cheats), cheats)
	}

	return -1, fmt.Errorf("not implemented")
}

func Part2(input *[]string) (int, error) {

	return -1, fmt.Errorf("not implemented")
}
