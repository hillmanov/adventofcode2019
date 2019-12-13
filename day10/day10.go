package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type location struct {
	Row int
	Col int
}

func main() {
	input, _ := ioutil.ReadFile("./input.txt")
	asteroidLocations := parseAsteroidLocations(string(input))

	// Part 1
	func() {
		visibleAsteroidsByLocation := make(map[location]int)
		for _, l := range asteroidLocations {
			visibleAsteroidsByLocation[l] = findVisibleAsteroidCount(l, asteroidLocations)
		}

		var locationWithMostVisible location
		for l, visible := range visibleAsteroidsByLocation {
			if visible > visibleAsteroidsByLocation[locationWithMostVisible] {
				locationWithMostVisible = l
			}
		}

		fmt.Printf("Part 1: %d\n", visibleAsteroidsByLocation[locationWithMostVisible])
	}()
}

func parseAsteroidLocations(s string) (asteroidField []location) {
	lines := strings.Split(s, "\n")
	for row, line := range lines {
		for col, object := range line {
			if string(object) == "#" {
				asteroidField = append(asteroidField, location{row, col})
			}
		}
	}
	return
}

func findVisibleAsteroidCount(origin location, asteroidLocations []location) int {
	// Find all other asteroids. Group by "slope" (rise/run).
	asteroidsBySightAngle := make(map[string][]location)
	for _, al := range asteroidLocations {
		if al.Row == origin.Row && al.Col == origin.Col {
			continue
		}

		direction := ""
		if al.Row > origin.Row {
			direction += "D"
		} else if al.Row < origin.Row {
			direction += "D"
		}

		if al.Col > origin.Col {
			direction += "R"
		} else if al.Col < origin.Col {
			direction += "L"
		}

		angle := (float64(al.Row) - float64(origin.Row)) / (float64(al.Col) - float64(origin.Col))
		key := fmt.Sprintf("%s:%f", direction, angle)

		asteroidsBySightAngle[key] = append(asteroidsBySightAngle[key], al)
	}

	// Sort all asteroids by distance

	return len(asteroidsBySightAngle)
}
