package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

	for _, line := range puzzleInput {
		fmt.Println(line)
	}

}
