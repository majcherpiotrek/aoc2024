package day_20

import (
	"fmt"
	"math"
	"slices"
	"time"
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

func simplePrintRaceTrack(racetrack *[][]byte, currentPosition []int, currentDirection []int, currentTime int, possibleCheats *[][]int) {
	toPrint := make([][]byte, len(*racetrack))
	for y := range *racetrack {
		rowToPrint := make([]byte, len((*racetrack)[y]))
		for x := range (*racetrack)[y] {
			if x == currentPosition[0] && y == currentPosition[1] {
				if currentDirection[0] == 1 {
					rowToPrint[x] = '>'
				}
				if currentDirection[0] == -1 {
					rowToPrint[x] = '<'
				}
				if currentDirection[1] == 1 {
					rowToPrint[x] = 'v'
				}
				if currentDirection[1] == -1 {
					rowToPrint[x] = '^'
				}

				continue
			}
			rowToPrint[x] = (*racetrack)[y][x]

			if possibleCheats != nil {
				for _, cheat := range *possibleCheats {
					if x == cheat[0] && y == cheat[1] {
						rowToPrint[x] = 'X'
					}
				}
			}

		}
		toPrint[y] = rowToPrint
	}

	fmt.Printf("Current time %d\n", currentTime)

	for _, row := range toPrint {
		fmt.Println(string(row))
	}

	fmt.Println("")
	time.Sleep(200 * time.Millisecond)
}

func printRaceTrack(racetrack *[][]byte, currentPosition []int, currentDirection []int, currentTime int, cheats *map[string]int, cheatStart []int, cheatDuration int) {
	toPrint := make([][]byte, len(*racetrack))
	for y := range *racetrack {
		rowToPrint := make([]byte, len((*racetrack)[y]))
		for x := range (*racetrack)[y] {
			if x == currentPosition[0] && y == currentPosition[1] {
				if currentDirection[0] == 1 {
					rowToPrint[x] = '>'
				}
				if currentDirection[0] == -1 {
					rowToPrint[x] = '<'
				}
				if currentDirection[1] == 1 {
					rowToPrint[x] = 'v'
				}
				if currentDirection[1] == -1 {
					rowToPrint[x] = '^'
				}
				continue
			}
			_, isInCheats := (*cheats)[fmt.Sprintf("%d,%d", x, y)]
			if isInCheats {
				rowToPrint[x] = 'X'
				continue
			}

			rowToPrint[x] = (*racetrack)[y][x]
		}
		toPrint[y] = rowToPrint
	}

	fmt.Printf("Current time %d, cheat start: %v, cheatDuration: %d\n", currentTime, cheatStart, cheatDuration)

	for _, row := range toPrint {
		fmt.Println(string(row))
	}

	fmt.Println("")
	time.Sleep(500 * time.Millisecond)
}

type Cheat struct {
	Start    []int
	End      []int
	Duration int
	TimeGain int
}

func startNewCheat(start []int) Cheat {
	return Cheat{
		Start:    start,
		Duration: 1,
	}
}

func raceWithCheats(
	racetrack *[][]byte,
	start []int,
	currentDirection []int,
	currentTime int,
	regularRaceTime int,
	cheatField []int,
	allowedCheatDuration int,
	cheatDuration int,
	cheats *map[string]int,
	triedCheatsCache *map[string]struct{},
) {

	if currentTime >= regularRaceTime {
		return
	}

	width := len((*racetrack)[0])
	height := len(*racetrack)

	if start[0] < 0 || start[0] >= width || start[1] < 0 || start[1] >= height {
		return
	}

	chDuration := cheatDuration
	if len(cheatField) > 0 {
		chDuration++
	}

	// printRaceTrack(racetrack, start, currentDirection, currentTime, cheats, cheatField, chDuration)

	field := (*racetrack)[start[1]][start[0]]

	if field == 'E' {
		if len(cheatField) > 0 {
			(*cheats)[fmt.Sprintf("%d,%d", cheatField[0], cheatField[1])] = currentTime
		}
		return
	}

	chField := cheatField
	/*
		hit the wall [0 3]
		no cheat field yet: [], 0

		hit the wall [0 4]
		already in cheat: [0 3], 1

		hit the wall [0 5]
		already in cheat: [0 3], 2

		hit the wall [0 6]
		already in cheat: [0 3], 3

		hit the wall [0 7]
		already in cheat: [0 3], 4
	*/

	if field == '#' {
		cacheKey := fmt.Sprintf("%v%v", start, currentDirection)
		_, alreadyVisited := (*triedCheatsCache)[cacheKey]
		if alreadyVisited {
			return
		}

		if len(cheatField) == 0 {
			chField = start
			chDuration = 1
			(*triedCheatsCache)[cacheKey] = struct{}{}
		} else {
			if chDuration >= allowedCheatDuration {
				return
			}
		}
	}

	next1 := addVectors(start, currentDirection)
	next2 := addVectors(start, rotatePlus90(currentDirection))
	next3 := addVectors(start, rotateMinus90(currentDirection))

	raceWithCheats(racetrack, next1, currentDirection, currentTime+1, regularRaceTime, chField, allowedCheatDuration, chDuration, cheats, triedCheatsCache)
	raceWithCheats(racetrack, next2, rotatePlus90(currentDirection), currentTime+1, regularRaceTime, chField, allowedCheatDuration, chDuration, cheats, triedCheatsCache)
	raceWithCheats(racetrack, next3, rotateMinus90(currentDirection), currentTime+1, regularRaceTime, chField, allowedCheatDuration, chDuration, cheats, triedCheatsCache)
}

type Cheats struct {
	Gain             int
	NumOfCheatFields int
}

type RaceTrack struct {
	Track           [][]byte
	TrackFields     [][]int
	Start           []int
	End             []int
	RegularRaceTime int
	TimesToFields   map[string]int
}

func NewRaceTrack(input *[]string) RaceTrack {
	track := make([][]byte, len(*input))
	var start, end []int
	regularRaceTime := 0
	for y, row := range *input {
		bytes := []byte(row)
		for x, fieldValue := range bytes {
			if fieldValue == 'S' {
				start = []int{x, y}
			}
			if fieldValue == 'E' {
				end = []int{x, y}
			}
			if fieldValue == '.' || fieldValue == 'E' {
				regularRaceTime++
			}

		}
		track[y] = bytes
	}

	timesToFields := make(map[string]int)

	t := 0
	field := start
	var prev *[]int = nil
	currentDirection := []int{1, 0}

	trackFields := make([][]int, 0, regularRaceTime)

	for track[field[1]][field[0]] != 'E' {
		trackFields = append(trackFields, field)
		timesToFields[fmt.Sprintf("%d,%d", field[0], field[1])] = t
		vect := currentDirection
		next := addVectors(field, currentDirection)

		for track[next[1]][next[0]] == '#' {
			vect = rotatePlus90(vect)
			next = addVectors(field, vect)

			if prev != nil && ((*prev)[0] == next[0] && (*prev)[1] == next[1]) {
				vect = rotatePlus90(vect)
				next = addVectors(field, vect)
			}
		}
		prev = &[]int{field[0], field[1]}
		field = next
		currentDirection = vect
		t++
	}

	timesToFields[fmt.Sprintf("%d,%d", end[0], end[1])] = t

	return RaceTrack{
		Track:           track,
		Start:           start,
		End:             end,
		RegularRaceTime: regularRaceTime,
		TimesToFields:   timesToFields,
		TrackFields:     trackFields,
	}
}

func calculateTimeGainForCheat(cheatStart []int, cheatEnd []int, fieldTimes *map[string]int) int {
	timeToStart, hasTime := (*fieldTimes)[fmt.Sprintf("%d,%d", cheatStart[0], cheatStart[1])]
	if !hasTime {
		return 0
	}
	timeToEnd, hasTime := (*fieldTimes)[fmt.Sprintf("%d,%d", cheatEnd[0], cheatEnd[1])]
	if !hasTime {
		return 0
	}

	if timeToStart >= timeToEnd {
		return 0
	}

	xDiff := int(math.Abs(float64(cheatEnd[0] - cheatStart[0])))
	yDiff := int(math.Abs(float64(cheatEnd[1] - cheatStart[1])))

	timeFromCheatStartToCheatEnd := xDiff + yDiff

	timeToEndWithCheat := timeToStart + timeFromCheatStartToCheatEnd

	if timeToEndWithCheat < timeToEnd {
		return timeToEnd - timeToEndWithCheat
	}

	return 0
}

func getFieldsInCheatRange(field []int, cheatRange int, trackWidth int, trackHeight int) [][]int {
	possibleCheats := make([][]int, 0)

	for xDiff := -1 * cheatRange; xDiff <= cheatRange; xDiff++ {
		xDiffAbs := int(math.Abs(float64(xDiff)))
		yRangeStart := -1 * (cheatRange - xDiffAbs)
		yRangeEnd := cheatRange - xDiffAbs

		for yDiff := yRangeStart; yDiff <= yRangeEnd; yDiff++ {
			x := field[0] + xDiff
			y := field[1] + yDiff

			if x >= 0 && x < trackWidth && y >= 0 && y < trackHeight {
				if x != field[0] || y != field[1] {
					possibleCheats = append(possibleCheats, []int{x, y})
				}
			}
		}
	}

	return possibleCheats
}

func Part1(input *[]string) (int, error) {
	racetrack := NewRaceTrack(input)

	cheatsTimeGains := make(map[string]int)

	fmt.Println(racetrack.TrackFields)

	for _, field := range racetrack.TrackFields {
		possibleCheats := getFieldsInCheatRange(field, 2, len((*input)[0]), len(*input))
		// simplePrintRaceTrack(&racetrack.Track, field, []int{0, -1}, 0, &possibleCheats)

		for _, cheat := range possibleCheats {
			cheatTimeGain := calculateTimeGainForCheat(field, cheat, &racetrack.TimesToFields)
			if cheatTimeGain > 0 {
				cheatsTimeGains[fmt.Sprintf("%d,%d-%d,%d", field[0], field[1], cheat[0], cheat[1])] = cheatTimeGain
			}
		}
	}

	inverseCheatTimesMap := make(map[int][]string)

	for cheatField, timeGain := range cheatsTimeGains {
		cheatFields, hasCheatFields := inverseCheatTimesMap[timeGain]
		if hasCheatFields {
			inverseCheatTimesMap[timeGain] = append(cheatFields, cheatField)
		} else {
			inverseCheatTimesMap[timeGain] = []string{cheatField}
		}

	}
	cheatArr := make([]Cheats, 0, len(inverseCheatTimesMap))
	for time, cheats := range inverseCheatTimesMap {
		cheatArr = append(cheatArr, Cheats{
			Gain:             time,
			NumOfCheatFields: len(cheats),
		})
	}

	slices.SortStableFunc(cheatArr, func(a Cheats, b Cheats) int {
		return a.Gain - b.Gain
	})

	bestCheats := 0
	for _, ch := range cheatArr {
		fmt.Printf("There are %d cheats that save %d picoseconds\n", ch.NumOfCheatFields, ch.Gain)
		if ch.Gain >= 100 {
			bestCheats += ch.NumOfCheatFields
		}
	}

	return bestCheats, nil
}

func Part2(input *[]string) (int, error) {
	racetrack := NewRaceTrack(input)

	cheatsTimeGains := make(map[string]int)

	for _, field := range racetrack.TrackFields {
		possibleCheats := getFieldsInCheatRange(field, 20, len((*input)[0]), len(*input))
		// simplePrintRaceTrack(&racetrack.Track, field, []int{0, -1}, 0, &possibleCheats)

		for _, cheat := range possibleCheats {
			cheatTimeGain := calculateTimeGainForCheat(field, cheat, &racetrack.TimesToFields)
			if cheatTimeGain > 0 {
				cheatsTimeGains[fmt.Sprintf("%d,%d-%d,%d", field[0], field[1], cheat[0], cheat[1])] = cheatTimeGain
			}
		}
	}

	inverseCheatTimesMap := make(map[int][]string)

	for cheatField, timeGain := range cheatsTimeGains {
		cheatFields, hasCheatFields := inverseCheatTimesMap[timeGain]
		if hasCheatFields {
			inverseCheatTimesMap[timeGain] = append(cheatFields, cheatField)
		} else {
			inverseCheatTimesMap[timeGain] = []string{cheatField}
		}

	}
	cheatArr := make([]Cheats, 0, len(inverseCheatTimesMap))
	for time, cheats := range inverseCheatTimesMap {
		cheatArr = append(cheatArr, Cheats{
			Gain:             time,
			NumOfCheatFields: len(cheats),
		})
	}

	slices.SortStableFunc(cheatArr, func(a Cheats, b Cheats) int {
		return a.Gain - b.Gain
	})

	bestCheats := 0
	for _, ch := range cheatArr {
		fmt.Printf("There are %d cheats that save %d picoseconds\n", ch.NumOfCheatFields, ch.Gain)
		if ch.Gain >= 100 {
			bestCheats += ch.NumOfCheatFields
		}
	}

	return bestCheats, nil
}
