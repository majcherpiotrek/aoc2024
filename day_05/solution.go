package day_05

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func parseLine(input *string, separator string) ([]int, error) {
	split := strings.Split(*input, separator)

	if len(split) == 0 {
		return []int{}, fmt.Errorf("Invalid input: %s", *input)
	}

	parsed := make([]int, len(split))

	for i, el := range split {
		num, err := strconv.Atoi(el)

		if err != nil {
			return []int{}, err
		}

		parsed[i] = num
	}

	return parsed, nil
}

func parseOrderingRule(input *string) ([]int, error) {
	parsedLine, err := parseLine(input, "|")

	if err != nil {
		return []int{}, err
	}

	if len(parsedLine) != 2 {
		return []int{}, fmt.Errorf("Invalid ordering rule: %s", *input)
	}

	return parsedLine, nil
}

type PrintConfig struct {
	OrderingRules [][]int
	UpdateConfigs [][]int
}

func (pc PrintConfig) String() string {
	return fmt.Sprintf("OrderingRules:\n%v,\nUpdateConfigs:\n%v\n", pc.OrderingRules, pc.UpdateConfigs)
}

func parsePrintConfig(input *[]string) (PrintConfig, error) {
	var orderingRules [][]int
	var updateConfigs [][]int

	orderingRulesParsed := false
	for _, line := range *input {
		if len(line) == 0 {
			orderingRulesParsed = true
			continue
		}

		if orderingRulesParsed {
			pagesToUpdate, err := parseLine(&line, ",")

			if err != nil {
				return PrintConfig{}, err
			}

			updateConfigs = append(updateConfigs, pagesToUpdate)
		} else {
			rule, err := parseOrderingRule(&line)

			if err != nil {
				return PrintConfig{}, err
			}

			orderingRules = append(orderingRules, rule)
		}
	}

	return PrintConfig{
		OrderingRules: orderingRules,
		UpdateConfigs: updateConfigs,
	}, nil
}

type OrderingRule struct {
	Before map[int]struct{}
	After  map[int]struct{}
}

func (or OrderingRule) String() string {
	before := make([]int, 0, len(or.Before))
	for key := range or.Before {
		before = append(before, key)
	}

	after := make([]int, 0, len(or.After))
	for key := range or.After {
		after = append(after, key)
	}

	return fmt.Sprintf("{\n\t\tBefore: %v,\n\t\tAfter: %v\n\t}", before, after)
}

type OrderingRulesMap map[int]OrderingRule

func (orm *OrderingRulesMap) ToString() string {
	str := "{\n"

	for key, value := range *orm {
		str = fmt.Sprintf("%s\t%d: %v,\n", str, key, value)

	}

	return fmt.Sprintf("%s}\n", str)
}

func (orm *OrderingRulesMap) ValidateUpdate(update []int) bool {
	for i, element := range update {
		var before []int
		var after []int

		if i > 0 {
			before = update[0:i]
		}

		if i < len(update)-1 {
			after = update[i:]
		}

		rulesForElement, hasRulesForElement := (*orm)[element]

		if !hasRulesForElement {
			continue
		}

		for _, preceedingElement := range before {
			_, hasElement := rulesForElement.Before[preceedingElement]

			if hasElement {
				return false
			}
		}

		for _, trailingElement := range after {
			_, hasElement := rulesForElement.After[trailingElement]

			if hasElement {
				return false
			}
		}
	}

	return true
}

func NewOrderingRulesMap(rules *[][]int) OrderingRulesMap {
	rulesMap := map[int]OrderingRule{}

	for _, rawRule := range *rules {
		before := rawRule[0]
		after := rawRule[1]

		existingRule, hasExistingRule := rulesMap[before]

		if hasExistingRule {
			existingRule.Before[after] = struct{}{}
		} else {
			Before := make(map[int]struct{})
			Before[after] = struct{}{}
			After := make(map[int]struct{})
			rulesMap[before] = OrderingRule{
				Before,
				After,
			}
		}

		existingRule, hasExistingRule = rulesMap[after]

		if hasExistingRule {
			existingRule.After[before] = struct{}{}
		} else {
			Before := make(map[int]struct{})
			After := make(map[int]struct{})
			After[before] = struct{}{}
			rulesMap[after] = OrderingRule{
				Before,
				After,
			}
		}
	}

	return rulesMap
}

func Part1(rows *[]string) (int, error) {
	printConfig, err := parsePrintConfig(rows)

	if err != nil {
		return -1, err
	}

	fmt.Println(fmt.Sprintf("%v", printConfig))

	rulesMap := NewOrderingRulesMap(&printConfig.OrderingRules)

	fmt.Println("Rules map:")
	fmt.Println(fmt.Sprintf("%s", rulesMap.ToString()))

	sum := 0

	for _, updateConfig := range printConfig.UpdateConfigs {
		isValid := rulesMap.ValidateUpdate(updateConfig)

		if isValid {
			idx := len(updateConfig) / 2
			sum += updateConfig[idx]
		}

		fmt.Println(fmt.Sprintf("%v - %t", updateConfig, isValid))

	}

	return sum, nil
}

func Part2(rows *[]string) (int, error) {
	printConfig, err := parsePrintConfig(rows)

	if err != nil {
		return -1, err
	}

	rulesMap := NewOrderingRulesMap(&printConfig.OrderingRules)

	sum := 0

	for _, updateConfig := range printConfig.UpdateConfigs {
		isValid := rulesMap.ValidateUpdate(updateConfig)

		if !isValid {
			fmt.Println(fmt.Sprintf("Before sort: %v", updateConfig))
			slices.SortStableFunc(updateConfig, func(a int, b int) int {
				rulesForA, hasRulesForA := rulesMap[a]

				res := 0

				if hasRulesForA {
					_, aAfterB := rulesForA.After[b]

					if aAfterB {
						res = 1
					}
					_, aBeforeB := rulesForA.Before[b]

					if aBeforeB {
						res = -1
					}
				}

				return res
			})

			fmt.Println(fmt.Sprintf("After sort: %v", updateConfig))

			idx := len(updateConfig) / 2
			sum += updateConfig[idx]
		}
	}

	return sum, nil
}
