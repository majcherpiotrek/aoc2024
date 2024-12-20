package day_19

import (
	"fmt"
	"slices"
	"strings"
)

func canBuildDesign(design string, availablePatterns *[]string, cache *map[string]bool) bool {
	if len(design) == 0 {
		return true
	}

	isPossibleDesign, isInCache := (*cache)[design]
	if isInCache {
		return isPossibleDesign
	}

	for _, pattern := range *availablePatterns {
		if design == pattern {
			(*cache)[design] = true
			return true
		}

		indexOfPattern := strings.Index(design, pattern)

		if indexOfPattern == -1 {
			continue
		}

		leadingSlice := design[0:indexOfPattern]
		leadingOk := canBuildDesign(leadingSlice, availablePatterns, cache)

		trailingSlice := design[indexOfPattern+len(pattern):]
		trailingOk := canBuildDesign(trailingSlice, availablePatterns, cache)

		if leadingOk && trailingOk {
			(*cache)[design] = true
			return true
		}
	}

	(*cache)[design] = false
	return false
}

func getFirstPrintPadding(depth int) string {
	padding := ""
	for i := 0; i < depth*4; i++ {
		if i == 0 {
			padding = "|"
		} else {
			padding = fmt.Sprintf("%s-", padding)
		}
	}
	return padding
}

func getPadding(depth int) string {
	padding := ""
	for i := 0; i < depth*4; i++ {
		padding = fmt.Sprintf("%s ", padding)
	}
	return padding
}

func findArrangements(design string, availablePatterns *[]string, possibleDesigns *map[string]bool, cache *map[string]int) int {
	isPossible, isInPossibleCache := (*possibleDesigns)[design]
	if isInPossibleCache && !isPossible {
		return 0
	}

	cacheValue, isInCache := (*cache)[design]
	if isInCache {
		return cacheValue
	}

	arrangements := 0

	for _, pattern := range *availablePatterns {
		patternLen := len(pattern)
		if len(design) < patternLen {
			continue
		}

		if design == pattern {
			arrangements++
			continue
		}

		if design[0:patternLen] == pattern {
			designSlice := design[patternLen:]
			tailArrangements := findArrangements(designSlice, availablePatterns, possibleDesigns, cache)
			if tailArrangements > 0 {
				arrangements += tailArrangements
			}
		}
	}

	(*cache)[design] = arrangements
	return arrangements
}

func Part1(input *[]string) (int, error) {
	availablePatterns := make([]string, 0)

	for _, pattern := range strings.Split((*input)[0], ", ") {
		availablePatterns = append(availablePatterns, pattern)
	}

	slices.SortStableFunc(availablePatterns, func(a string, b string) int {
		return len(b) - len(a)
	})

	possibleDesigns := 0
	possibleDesignsCache := make(map[string]bool)
	for i := 2; i < len(*input); i++ {
		fmt.Printf("Design %d\n", i-1)
		design := (*input)[i]

		if canBuildDesign(design, &availablePatterns, &possibleDesignsCache) {
			possibleDesigns++
			fmt.Println(design)
		}
	}

	return possibleDesigns, nil
}

func Part2(input *[]string) (int, error) {
	availablePatterns := make([]string, 0)

	for _, pattern := range strings.Split((*input)[0], ", ") {
		availablePatterns = append(availablePatterns, pattern)
	}

	slices.SortStableFunc(availablePatterns, func(a string, b string) int {
		return len(b) - len(a)
	})

	possibleDesignsCache := make(map[string]bool)
	for i := 2; i < len(*input); i++ {
		design := (*input)[i]
		canBuildDesign(design, &availablePatterns, &possibleDesignsCache)
	}

	arrangementsCache := make(map[string]int)

	sum := 0
	for i := 2; i < len(*input); i++ {
		design := (*input)[i]
		arrangements := findArrangements(design, &availablePatterns, &possibleDesignsCache, &arrangementsCache)
		fmt.Printf("%s: %d\n", design, arrangements)
		sum += arrangements
	}

	return sum, nil
}
