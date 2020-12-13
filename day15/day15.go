package main

import (
	"adventofcode/intcode"
	"fmt"
	"math"
)

// Movement
const (
	north = 1
	south = 2
	east  = 3
	west  = 4
)

const oxygenated = 9

// Status codes
const (
	wall  = 0
	moved = 1
	tank  = 2
)

type Location struct {
	Row int
	Col int
}

func main() {
	program := intcode.ReadIntcodeProgram("./input.txt")

	// Part 1
	grid := make(map[Location]int)
	var oxygenTankLocation Location

	path := make([]int, 0)
	currentLocation := Location{Row: 0, Col: 0}
	direction := north
	stepsToTank := 0

	icc := intcode.NewIntCodeComputer(intcode.CopyIntcodeProgram(program))
	icc.RequestInput = true
	go icc.Run()

	for direction != -2 {
		select {
		case <-icc.InputChannel:
			icc.InputChannel <- direction
		case statusCode := <-icc.OutputChannel:
			switch statusCode {
			case wall:
				wallLocation := applyDirection(currentLocation, direction)
				grid[wallLocation] = wall
			case moved:
				currentLocation = applyDirection(currentLocation, direction)
				if _, ok := grid[currentLocation]; !ok {
					grid[currentLocation] = moved
					path = append(path, direction)
				}
			case tank:
				currentLocation = applyDirection(currentLocation, direction)
				oxygenTankLocation = currentLocation
				grid[currentLocation] = tank
				path = append(path, direction)
				stepsToTank = len(path)
			}

			direction, path = getDirection(grid, currentLocation, path)
		case <-icc.DoneChannel:
			return
		}
	}

	fmt.Printf("Part 1 (Steps to tank): %d\n", stepsToTank)

	// Part 2
	currentLocation = oxygenTankLocation
	minutes := 0
	var oxygenate func(l Location, minutes int) int
	oxygenate = func(l Location, currentMinutes int) int {
		switch grid[l] {
		case wall, oxygenated:
			return currentMinutes - 1
		case moved, tank:
			grid[l] = oxygenated
			for _, oxygenateLocation := range []Location{applyDirection(l, north), applyDirection(l, south), applyDirection(l, east), applyDirection(l, west)} {
				minutes = max(minutes, oxygenate(oxygenateLocation, currentMinutes+1))
			}
		}
		return minutes
	}

	minutes = oxygenate(oxygenTankLocation, 0)
	fmt.Printf("Part 2 (Minutes to full oxygenation): %+v\n", minutes)
}

func getDirection(grid map[Location]int, l Location, path []int) (int, []int) {
	northDirection := applyDirection(l, north)
	southDirection := applyDirection(l, south)
	eastDirection := applyDirection(l, east)
	westDirection := applyDirection(l, west)

	direction := -1
	if _, ok := grid[northDirection]; !ok {
		direction = north
	} else if _, ok := grid[southDirection]; !ok {
		direction = south
	} else if _, ok := grid[eastDirection]; !ok {
		direction = east
	} else if _, ok := grid[westDirection]; !ok {
		direction = west
	}

	if direction == -1 {
		if len(path) == 0 {
			return -2, path
		}
		direction, path = getReverseDirection(path[len(path)-1]), path[:len(path)-1]
		return direction, path
	}
	return direction, path
}

func getReverseDirection(direction int) int {
	switch direction {
	case north:
		return south
	case south:
		return north
	case east:
		return west
	case west:
		return east
	}
	return -1
}

func applyDirection(l Location, direction int) Location {
	newLocation := Location{
		Row: l.Row,
		Col: l.Col,
	}
	switch direction {
	case north:
		newLocation.Row--
	case south:
		newLocation.Row++
	case east:
		newLocation.Col++
	case west:
		newLocation.Col--
	}
	return newLocation
}

func displayGrid(g map[Location]int, l Location) {
	// Get bounds of grid
	minRow := math.MaxInt64
	maxRow := math.MinInt64
	minCol := math.MaxInt64
	maxCol := math.MinInt64

	for loc, _ := range g {
		minRow = min(minRow, loc.Row)
		maxRow = max(maxRow, loc.Row)
		minCol = min(minCol, loc.Col)
		maxCol = max(maxCol, loc.Col)
	}

	// Clear the screen quickly.
	fmt.Printf("\033[0;0H")
	for row := minRow; row < maxRow+1; row++ {
		for col := minCol; col < maxCol+1; col++ {
			id, ok := g[Location{row, col}]
			if l.Row == row && l.Col == col {
				fmt.Print("*")
			} else if row == 0 && col == 0 {
				fmt.Printf("O")
			} else if !ok {
				fmt.Printf(" ")
			} else if id == moved {
				fmt.Printf(".")
			} else if id == tank {
				fmt.Printf("T")
			} else if ok && id == wall {
				fmt.Printf("#")
			} else if ok && id == oxygenated {
				fmt.Printf("9")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Println("--------------------")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
