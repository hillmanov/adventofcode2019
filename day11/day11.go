package main

import (
	"adventofcode/intcode"
	"fmt"
	"math"
)

const white = 1
const black = 0

const left = 0
const right = 1

type Location struct {
	Row int
	Col int
}

type Robot struct {
	Grid     map[Location]int
	Location Location
	Heading  string
}

func (r *Robot) CurrentColor() int {
	if currentColor, ok := r.Grid[r.Location]; ok {
		return currentColor
	}
	return black
}

func (r *Robot) Paint(c int) {
	r.Grid[r.Location] = c
}

func (r *Robot) TurnAndMove(d int) {
	switch d {
	case left:
		switch r.Heading {
		case "U":
			r.Heading = "L"
		case "D":
			r.Heading = "R"
		case "L":
			r.Heading = "D"
		case "R":
			r.Heading = "U"
		}
	case right:
		switch r.Heading {
		case "U":
			r.Heading = "R"
		case "D":
			r.Heading = "L"
		case "L":
			r.Heading = "U"
		case "R":
			r.Heading = "D"
		}
	}
	switch r.Heading {
	case "U":
		r.Location.Row--
	case "D":
		r.Location.Row++
	case "L":
		r.Location.Col--
	case "R":
		r.Location.Col++
	}
}

func main() {
	program := intcode.ReadIntcodeProgram("./input.txt")

	// Part 1
	func() {
		grid := make(map[Location]int)
		start := Location{Row: 0, Col: 0}

		robot := Robot{
			Grid:     grid,
			Location: start,
			Heading:  "U",
		}

		icc := intcode.NewIntCodeComputer(program)
		icc.RequestInput = true
		go icc.Run()

		for {
			select {
			case <-icc.InputChannel:
				icc.InputChannel <- robot.CurrentColor()
			case color := <-icc.OutputChannel:
				direction := <-icc.OutputChannel
				robot.Paint(color)
				robot.TurnAndMove(direction)
			case <-icc.DoneChannel:
				fmt.Printf("Part 1: %d\n", len(robot.Grid))
				return
			}
		}
	}()

	// Part 2
	func() {
		grid := make(map[Location]int)
		start := Location{Row: 0, Col: 0}

		// Start on the single white panel
		grid[start] = white

		robot := Robot{
			Grid:     grid,
			Location: start,
			Heading:  "U",
		}

		icc := intcode.NewIntCodeComputer(program)
		icc.RequestInput = true
		go icc.Run()

		for {
			select {
			case <-icc.InputChannel:
				icc.InputChannel <- robot.CurrentColor()
			case color := <-icc.OutputChannel:
				direction := <-icc.OutputChannel
				robot.Paint(color)
				robot.TurnAndMove(direction)
			case <-icc.DoneChannel:
				fmt.Println("Part 2:")
				// Get bounds of grid
				minRow := math.MaxInt64
				maxRow := math.MinInt64
				minCol := math.MaxInt64
				maxCol := math.MinInt64

				for loc, _ := range robot.Grid {
					minRow = min(minRow, loc.Row)
					maxRow = max(maxRow, loc.Row)
					minCol = min(minCol, loc.Col)
					maxCol = max(maxCol, loc.Col)
				}

				for row := minRow; row < maxRow+1; row++ {
					for col := minCol; col < maxCol+1; col++ {
						c, _ := grid[Location{row, col}]
						if c == white {
							fmt.Printf("#")
						} else {
							fmt.Printf(" ")
						}
					}
					fmt.Printf("\n")
				}
				return
			}
		}
	}()
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
