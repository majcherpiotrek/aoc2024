package day_12

import (
	"fmt"
	"math"
	"slices"
)

type Field struct {
	X         int
	Y         int
	Perimeter int
}

func (field *Field) touchesField(x int, y int) bool {
	diffX := int(math.Abs(float64(x - field.X)))
	diffY := int(math.Abs(float64(y - field.Y)))

	return (diffX == 0 && diffY == 1) || (diffX == 1 && diffY == 0)
}

func (field Field) String() string {
	return fmt.Sprintf("{ x: %d, y: %d, perimeter: %d }", field.X, field.Y, field.Perimeter)
}

type Region struct {
	Plant     byte
	Fields    map[string]Field
	Perimeter map[string]Field
}

func (region *Region) calculatePrice() int {
	perimeter := 0

	for _, field := range region.Fields {
		perimeter += field.Perimeter
	}

	return len(region.Fields) * perimeter
}

func encodeFieldKey(x int, y int) string {
	return fmt.Sprintf("%d-%d", x, y)
}

type Neighbours struct {
	Top   *Field
	Right *Field
	Down  *Field
	Left  *Field
	All   []Field
}

func (region *Region) getFieldNeighbours(x int, y int) Neighbours {
	top := encodeFieldKey(x, y-1)
	right := encodeFieldKey(x+1, y)
	down := encodeFieldKey(x, y+1)
	left := encodeFieldKey(x-1, y)

	result := Neighbours{
		Top:   nil,
		Right: nil,
		Down:  nil,
		Left:  nil,
		All:   []Field{},
	}

	neighbours := make([]Field, 0, 4)

	topField, exists := region.Fields[top]
	if exists {
		neighbours = append(neighbours, topField)
		result.Top = &topField
	}
	rightField, exists := region.Fields[right]
	if exists {
		neighbours = append(neighbours, rightField)
		result.Right = &rightField
	}
	downField, exists := region.Fields[down]
	if exists {
		neighbours = append(neighbours, downField)
		result.Down = &downField
	}
	leftField, exists := region.Fields[left]
	if exists {
		neighbours = append(neighbours, leftField)
		result.Left = &leftField
	}

	result.All = neighbours

	return result
}

func (region *Region) touchesField(x int, y int) bool {
	neighbours := region.getFieldNeighbours(x, y)

	return len(neighbours.All) > 0
}

func (region *Region) addField(x int, y int) {
	neighbours := region.getFieldNeighbours(x, y)
	newField := Field{
		X:         x,
		Y:         y,
		Perimeter: 4 - len(neighbours.All),
	}

	newFieldKey := encodeFieldKey(x, y)
	region.Fields[newFieldKey] = newField

	if newField.Perimeter > 0 {
		region.Perimeter[newFieldKey] = newField
	}

	for _, n := range neighbours.All {
		neighbourFieldKey := encodeFieldKey(n.X, n.Y)
		updatedNeighbour := Field{
			X:         n.X,
			Y:         n.Y,
			Perimeter: n.Perimeter - 1,
		}
		region.Fields[neighbourFieldKey] = updatedNeighbour

		if n.Perimeter > 0 && updatedNeighbour.Perimeter == 0 {
			delete(region.Perimeter, neighbourFieldKey)
		}

		if updatedNeighbour.Perimeter > 0 {
			region.Perimeter[neighbourFieldKey] = updatedNeighbour
		}
	}
}

func (region *Region) merge(regionToMerge Region) {
	for _, field := range regionToMerge.Fields {
		region.addField(field.X, field.Y)
	}
}

func (region Region) String() string {
	return fmt.Sprintf("{ plant: %s, fields: %v }", string(region.Plant), region.Fields)
}

func flood(start []int, plane *[]string, region *Region, visited *map[string]struct{}) {
	planeHeight := len(*plane)
	planeWidth := len((*plane)[0])

	x := start[0]
	y := start[1]

	if x < 0 || x >= planeWidth || y < 0 || y >= planeHeight {
		return
	}

	fieldKey := encodeFieldKey(x, y)

	_, alreadyVisited := (*visited)[fieldKey]
	if alreadyVisited {
		return
	}

	plantOnField := (*plane)[y][x]

	if plantOnField != region.Plant {
		return
	}

	region.addField(x, y)
	(*visited)[fieldKey] = struct{}{}

	flood([]int{x, y - 1}, plane, region, visited)
	flood([]int{x + 1, y}, plane, region, visited)
	flood([]int{x, y + 1}, plane, region, visited)
	flood([]int{x - 1, y}, plane, region, visited)
}

func parseRegions(rows *[]string) []Region {
	allRegions := []Region{}

	visited := make(map[string]struct{})

	for y := 0; y < len(*rows); y++ {
		for x := 0; x < len((*rows)[y]); x++ {
			fieldKey := encodeFieldKey(x, y)
			_, alreadyVisited := visited[fieldKey]

			if !alreadyVisited {
				plant := (*rows)[y][x]
				region := Region{
					Plant:     plant,
					Fields:    make(map[string]Field),
					Perimeter: make(map[string]Field),
				}

				flood([]int{x, y}, rows, &region, &visited)

				allRegions = append(allRegions, region)
			}
		}

	}

	return allRegions
}

func calculateWalls(currentField []int, region *Region, walls *map[string][][]int, visited *map[string]struct{}) {
	fieldKey := encodeFieldKey(currentField[0], currentField[1])

	_, alreadyVisited := (*visited)[fieldKey]
	if alreadyVisited {
		return
	}

	(*visited)[fieldKey] = struct{}{}

	_, isInRegion := region.Fields[fieldKey]

	if !isInRegion {
		return
	}

	neighbours := region.getFieldNeighbours(currentField[0], currentField[1])

	if currentField[0] == 8 && currentField[1] == 4 {
		fmt.Printf("N: %v\n", neighbours)

	}

	if neighbours.Top == nil {
		// Y-Top
		wallKey := fmt.Sprintf("%d-top", currentField[1])
		wall, alreadyExists := (*walls)[wallKey]
		if alreadyExists {
			(*walls)[wallKey] = append(wall, currentField)
		} else {
			(*walls)[wallKey] = [][]int{currentField}
		}

	}
	if neighbours.Down == nil {
		// Y-down
		wallKey := fmt.Sprintf("%d-down", currentField[1])
		wall, alreadyExists := (*walls)[wallKey]
		if alreadyExists {
			(*walls)[wallKey] = append(wall, currentField)
		} else {
			(*walls)[wallKey] = [][]int{currentField}
		}
	}
	if neighbours.Right == nil {
		// X-right
		wallKey := fmt.Sprintf("%d-right", currentField[0])
		wall, alreadyExists := (*walls)[wallKey]
		if alreadyExists {
			(*walls)[wallKey] = append(wall, currentField)
		} else {
			(*walls)[wallKey] = [][]int{currentField}
		}
	}
	if neighbours.Left == nil {
		// X-left
		wallKey := fmt.Sprintf("%d-left", currentField[0])
		wall, alreadyExists := (*walls)[wallKey]
		if alreadyExists {
			(*walls)[wallKey] = append(wall, currentField)
		} else {
			(*walls)[wallKey] = [][]int{currentField}
		}
	}

	if len(neighbours.All) == 0 {
		// A single field, will have 4 walls
		return
	}

	calculateWalls([]int{currentField[0], currentField[1] - 1}, region, walls, visited)
	calculateWalls([]int{currentField[0] + 1, currentField[1]}, region, walls, visited)
	calculateWalls([]int{currentField[0], currentField[1] + 1}, region, walls, visited)
	calculateWalls([]int{currentField[0] - 1, currentField[1]}, region, walls, visited)
}

func Part1(rows *[]string) (int, error) {
	allRegions := parseRegions(rows)
	totalPrice := 0

	for _, region := range allRegions {
		perimeter := 0
		for _, field := range region.Perimeter {
			perimeter += field.Perimeter
		}

		price := perimeter * len(region.Fields)
		totalPrice += price
		fmt.Printf("Region %c - size %d - price %d\n", region.Plant, len(region.Fields), price)
	}

	return totalPrice, nil
}

func splitWall(wall [][]int) [][][]int {
	fmt.Printf("INPUT WALL %v\n", wall)

	if len(wall) <= 1 {
		fmt.Println("JUST ONE ELEMENT")
		return [][][]int{wall}
	}

	slices.SortStableFunc(wall, func(a []int, b []int) int {
		return a[0] - b[0] + a[1] - b[1]
	})
	fmt.Printf("SORTED     %v\n", wall)

	result := [][][]int{}
	currentWall := [][]int{}

	for i := 0; i < len(wall); i++ {
		if len(currentWall) == 0 {
			currentWall = append(currentWall, wall[i])
			continue
		}

		currentWallLastElement := currentWall[len(currentWall)-1]

		diffX := int(math.Abs(float64(currentWallLastElement[0] - wall[i][0])))
		diffY := int(math.Abs(float64(currentWallLastElement[1] - wall[i][1])))

		if diffX > 1 || diffY > 1 {
			result = append(result, currentWall)
			currentWall = [][]int{wall[i]}
		} else {
			currentWall = append(currentWall, wall[i])
		}
	}

	result = append(result, currentWall)

	fmt.Printf("WALL TO SPLIT: %v\n", wall)
	fmt.Printf("RESULT         %v\n", result)

	return result
}

func Part2(rows *[]string) (int, error) {
	allRegions := parseRegions(rows)
	totalPrice := 0

	for _, region := range allRegions {
		regionWalls := make(map[string][][]int)
		visited := make(map[string]struct{})
		var fieldToStart Field

		for _, field := range region.Perimeter {
			fieldToStart = field
			break
		}
		calculateWalls([]int{fieldToStart.X, fieldToStart.Y}, &region, &regionWalls, &visited)
		fmt.Printf("region walls: %v\n", regionWalls)

		if region.Plant == 'F' {
			fmt.Printf("VISITED %v", visited)
		}

		distinctWalls := [][][]int{}

		for wallId, walls := range regionWalls {
			fmt.Printf("Wall id %s\n", wallId)
			distinctWalls = append(distinctWalls, splitWall(walls)...)
		}

		price := len(distinctWalls) * len(region.Fields)
		totalPrice += price
		fmt.Printf("\nRegion %c - size %d - price %d, distinct walls: %d, total walls: %d\n", region.Plant, len(region.Fields), price, len(distinctWalls), len(regionWalls))
		fmt.Println(regionWalls)
		fmt.Println(distinctWalls)
		fmt.Println("---------------------------------------------")
	}

	return totalPrice, nil
}
