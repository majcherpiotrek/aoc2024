package main

import (
	"aoc2024/day_01"
	"aoc2024/day_02"
	"aoc2024/day_03"
	"aoc2024/day_04"
	"aoc2024/day_05"
	"aoc2024/day_06"
	"aoc2024/day_07"
	"aoc2024/day_08"
	"aoc2024/day_09"
	"aoc2024/day_10"
	"aoc2024/day_11"
	"aoc2024/day_12"
	"aoc2024/day_14"
	"aoc2024/day_15"
	"aoc2024/day_16"
	"aoc2024/day_17"
	"aoc2024/day_18"
	"aoc2024/day_19"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ChallengeId struct {
	Day    int
	Part   int
	IsTest bool
}

func getDayAndPartFromArgs() (ChallengeId, error) {
	args := os.Args[1:]

	if len(args) < 2 {
		return ChallengeId{}, fmt.Errorf("You need to pass a day and part of the challenge.")
	}

	day, err := strconv.Atoi(args[0])

	if err != nil {
		return ChallengeId{}, fmt.Errorf("Failed to parse challenge day")
	}

	part, err := strconv.Atoi(args[1])

	if err != nil {
		return ChallengeId{}, fmt.Errorf("Failed to parse challenge part")

	}

	if len(args) == 3 && args[2] == "-t" {
		return ChallengeId{Day: day, Part: part, IsTest: true}, nil
	}

	return ChallengeId{Day: day, Part: part, IsTest: false}, nil
}

func readInputFile(challenge ChallengeId) ([]string, error) {
	dayPath := fmt.Sprintf("day_%02d", challenge.Day)

	fileName := "puzzle_input"

	if challenge.IsTest {
		fileName = "test_input"
	}

	filePath := filepath.Join(dayPath, fileName)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file '%s': %w", filePath, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading file '%s': %w", filePath, err)
	}

	return lines, nil
}

func runChallange(challengeId ChallengeId, puzzleInput *[]string) (int, error) {
	switch challengeId.Day {
	case 1:
		if challengeId.Part == 1 {
			return day_01.Part1(puzzleInput)
		} else {
			return day_01.Part2(puzzleInput)
		}

	case 2:
		if challengeId.Part == 1 {
			return day_02.Part1(puzzleInput)
		} else {
			return day_02.Part2(puzzleInput)
		}
	case 3:
		if challengeId.Part == 1 {
			return day_03.Part1(puzzleInput)
		} else {
			return day_03.Part2(puzzleInput)
		}
	case 4:
		if challengeId.Part == 1 {
			return day_04.Part1(puzzleInput)
		} else {
			return day_04.Part2(puzzleInput)
		}
	case 5:
		if challengeId.Part == 1 {
			return day_05.Part1(puzzleInput)
		} else {
			return day_05.Part2(puzzleInput)
		}
	case 6:
		if challengeId.Part == 1 {
			return day_06.Part1(puzzleInput)
		} else {
			return day_06.Part2(puzzleInput)
		}
	case 7:
		if challengeId.Part == 1 {
			return day_07.Part1(puzzleInput)
		} else {
			return day_07.Part2(puzzleInput)
		}
	case 8:
		if challengeId.Part == 1 {
			return day_08.Part1(puzzleInput)
		} else {
			return day_08.Part2(puzzleInput)
		}
	case 9:
		if challengeId.Part == 1 {
			return day_09.Part1(puzzleInput)
		} else {
			return day_09.Part2(puzzleInput)
		}
	case 10:
		if challengeId.Part == 1 {
			return day_10.Part1(puzzleInput)
		} else {
			return day_10.Part2(puzzleInput)
		}
	case 11:
		if challengeId.Part == 1 {
			return day_11.Part1(puzzleInput)
		} else {
			return day_11.Part2(puzzleInput)
		}
	case 12:
		if challengeId.Part == 1 {
			return day_12.Part1(puzzleInput)
		} else {
			return day_12.Part2(puzzleInput)
		}
	case 14:
		if challengeId.Part == 1 {
			return day_14.Part1(puzzleInput)
		} else {
			return day_14.Part2(puzzleInput)
		}
	case 15:
		if challengeId.Part == 1 {
			return day_15.Part1(puzzleInput)
		} else {
			return day_15.Part2(puzzleInput)
		}
	case 16:
		if challengeId.Part == 1 {
			return day_16.Part1(puzzleInput)
		} else {
			return day_16.Part2(puzzleInput)
		}
	case 17:
		if challengeId.Part == 1 {
			return day_17.Part1(puzzleInput)
		} else {
			return day_17.Part2(puzzleInput)
		}
	case 18:
		if challengeId.Part == 1 {
			return day_18.Part1(puzzleInput)
		} else {
			return day_18.Part2(puzzleInput)
		}
	case 19:
		if challengeId.Part == 1 {
			return day_19.Part1(puzzleInput)
		} else {
			return day_19.Part2(puzzleInput)
		}
	default:
		return -1, fmt.Errorf("Not implemented yet")
	}
}

func main() {

	challengeId, err := getDayAndPartFromArgs()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Println(fmt.Sprintf("Advent of Code 2023, day %d, part %d", challengeId.Day, challengeId.Part))

	if challengeId.IsTest {
		fmt.Println("[TEST DATA]")
	}

	puzzleInput, err := readInputFile(challengeId)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	start := time.Now()
	result, err := runChallange(challengeId, &puzzleInput)
	executionTime := time.Since(start)

	if err != nil {
		fmt.Fprintln(os.Stderr, "An error occurred when solving the puzzle:", err)
		os.Exit(1)
	}

	fmt.Println(fmt.Sprintf("Result: %d, execution time: %v", result, executionTime))

	os.Exit(0)
}
