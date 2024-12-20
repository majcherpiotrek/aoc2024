package day_19

import (
	"fmt"
	"slices"
	"strings"
)

type PatternColor byte

const (
	white = 'w'
	blue  = 'u'
	black = 'b'
	red   = 'r'
	green = 'g'
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

func findArrangements(design string, availablePatterns *[]string, canBuildDesignCache *map[string]bool, arrCache *map[string]map[string]struct{}, depth int) map[string]struct{} {
	firstPrintPadding := getFirstPrintPadding(depth)
	padding := getPadding(depth)
	// fmt.Printf("%sCalled with design %s\n", firstPrintPadding, design)
	// fmt.Printf("CACHE %v\n\n", *cache)
	fmt.Printf("%s%dStarting design %s\n", firstPrintPadding, depth, design)
	canBuild, alreadyChecked := (*canBuildDesignCache)[design]
	if !alreadyChecked {
		canBuild = canBuildDesign(design, availablePatterns, canBuildDesignCache)
	}

	if !canBuild {
		fmt.Printf("%s%d Can't build this design\n", padding, depth)
		return map[string]struct{}{}
	}

	cachedArrangements, isInCache := (*arrCache)[design]
	if isInCache {
		fmt.Printf("%s%dCACHE HIT design %s\n", padding, depth, design)
		return cachedArrangements
	}
	fmt.Printf("%s%d No cache hit design %s\n", padding, depth, design)

	arrangements := make(map[string]struct{})

	for _, pattern := range *availablePatterns {
		fmt.Printf("%s%dChecking pattern %s for design %s\n", padding, depth, pattern, design)
		if design == pattern {
			// fmt.Printf("%sPattern equals design\n", padding)
			arrangements[pattern] = struct{}{}
			(*arrCache)[pattern] = map[string]struct{}{
				pattern: {},
			}

			if len(design) == 1 {
				return arrangements
			}
			continue
		}

		indexOfPattern := strings.Index(design, pattern)

		if indexOfPattern == -1 {
			// fmt.Printf("%sPattern not found\n", padding)
			continue
		}
		// fmt.Printf("%sPattern %s present at %d. Going deeper ->\n", padding, pattern, indexOfPattern)

		leadingSlice := design[0:indexOfPattern]
		//fmt.Printf("%sLeading slice: %v\n", padding, leadingSlice)
		arrangementsLeading := make(map[string]struct{})
		if len(leadingSlice) > 0 {
			arrangementsLeading = findArrangements(leadingSlice, availablePatterns, canBuildDesignCache, arrCache, depth+1)
		}

		trailingSlice := design[indexOfPattern+len(pattern):]
		//fmt.Printf("%sTrailing slice: %v\n", padding, leadingSlice)
		arrangementsTrailing := make(map[string]struct{})
		if len(trailingSlice) > 0 {
			arrangementsTrailing = findArrangements(trailingSlice, availablePatterns, canBuildDesignCache, arrCache, depth+1)
		}

		if len(leadingSlice) == 0 {
			fmt.Printf("%s%d Leading slice empty", padding, depth)
			for arr := range arrangementsTrailing {
				arrangements[fmt.Sprintf("%s,%s", pattern, arr)] = struct{}{}
			}
			continue
		}

		if len(trailingSlice) == 0 {
			fmt.Printf("%s%d Trailing slice empty", padding, depth)
			for arr := range arrangementsLeading {
				arrangements[fmt.Sprintf("%s,%s", arr, pattern)] = struct{}{}
			}
			continue
		}

		if len(arrangementsLeading) > 0 && len(arrangementsTrailing) > 0 {
			fmt.Printf("%s%d Merging possibilities from leading (%d) and trailing (%d)\n", padding, depth, len(arrangementsLeading), len(arrangementsTrailing))
			for leading := range arrangementsLeading {
				for trailing := range arrangementsTrailing {
					arr := fmt.Sprintf("%s,%s,%s", leading, pattern, trailing)
					arrangements[arr] = struct{}{}
				}
			}
		}
	}

	(*arrCache)[design] = arrangements
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

	arrangementsCache := make(map[string]map[string]struct{})
	canBuildDesignCache := make(map[string]bool)
	sum := 0
	for i := 2; i < len(*input); i++ {
		fmt.Printf("Design %d\n", i-1)
		design := (*input)[i]
		arrangements := findArrangements(design, &availablePatterns, &canBuildDesignCache, &arrangementsCache, 0)
		fmt.Printf("%s: %d\n", design, len(arrangements))
		sum += len(arrangements)
	}

	return sum, nil
}
