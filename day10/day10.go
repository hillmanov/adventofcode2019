package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strings"
)

type location struct {
	Row int
	Col int
}

type directionAndAngle struct {
	Direction string
	Angle     float64
}

func (l location) distanceFrom(d location) float64 {
	return math.Abs(math.Sqrt(math.Pow(float64(d.Col)-float64(l.Col), 2) + math.Pow(float64(d.Row)-float64(l.Row), 2)))
}

func main() {
	input, _ := ioutil.ReadFile("./input.txt")
	asteroidLocations := parseAsteroidLocations(string(input))

	visibleAsteroidsByLocation := make(map[location]map[directionAndAngle][]location)
	for _, l := range asteroidLocations {
		visibleAsteroidsByLocation[l] = findVisibleAsteroids(l, asteroidLocations)
	}

	var locationWithMostVisible location
	for l, visible := range visibleAsteroidsByLocation {
		if len(visible) > len(visibleAsteroidsByLocation[locationWithMostVisible]) {
			locationWithMostVisible = l
		}
	}

	fmt.Printf("Part 1: %d\n", len(visibleAsteroidsByLocation[locationWithMostVisible]))

	clockwiseRotationOrder := []string{"U", "UR", "R", "DR", "D", "DL", "L", "UL"}
	anglesAndAsteroids := visibleAsteroidsByLocation[locationWithMostVisible]

	// Get a list of all of the angles
	angles := make([]directionAndAngle, 0)
	for da, _ := range anglesAndAsteroids {
		angles = append(angles, da)
	}

	// Sort the angles by clockwise order
	sort.Slice(angles, func(i, j int) bool {
		a := angles[i]
		b := angles[j]
		if indexOf(clockwiseRotationOrder, a.Direction) < indexOf(clockwiseRotationOrder, b.Direction) {
			return true
		} else if indexOf(clockwiseRotationOrder, a.Direction) > indexOf(clockwiseRotationOrder, b.Direction) {
			return false
		} else {
			return a.Angle < b.Angle
		}
	})

	var lastVaporized location
	for i := 0; i < 200; i++ {
		var currentAngle directionAndAngle
		j := i
		for {
			currentAngle = angles[j%len(angles)]
			if len(anglesAndAsteroids[currentAngle]) == 0 {
				j += 1
			} else {
				break
			}
		}
		lastVaporized = anglesAndAsteroids[currentAngle][0]
		anglesAndAsteroids[currentAngle] = anglesAndAsteroids[currentAngle][1:]
	}

	r := (lastVaporized.Col * 100) + lastVaporized.Row
	fmt.Printf("Part 2: %d\n", r)
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

func findVisibleAsteroids(origin location, asteroidLocations []location) map[directionAndAngle][]location {
	// Find all other asteroids. Group by "slope" (rise/run).
	asteroidsBySightAngle := make(map[directionAndAngle][]location)
	for _, al := range asteroidLocations {
		if al.Row == origin.Row && al.Col == origin.Col {
			continue
		}

		da := directionAndAngle{}
		if al.Row <= origin.Row {
			da.Direction += "U"
		} else if al.Row >= origin.Row {
			da.Direction += "D"
		}

		if al.Col > origin.Col {
			da.Direction += "R"
		} else if al.Col < origin.Col {
			da.Direction += "L"
		}

		angle := (float64(al.Row) - float64(origin.Row)) / (float64(al.Col) - float64(origin.Col))
		da.Angle = angle

		asteroidsBySightAngle[da] = append(asteroidsBySightAngle[da], al)
	}

	for _, asteroids := range asteroidsBySightAngle {
		sort.Slice(asteroids, func(i, j int) bool {
			return origin.distanceFrom(asteroids[i]) < origin.distanceFrom(asteroids[j])
		})
	}

	return asteroidsBySightAngle
}

func indexOf(haystack []string, needle string) int {
	for i, v := range haystack {
		if v == needle {
			return i
		}
	}
	return -1
}
