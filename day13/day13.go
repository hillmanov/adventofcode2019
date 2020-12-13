package main

import (
	"adventofcode/intcode"
	"fmt"
	"math"
	"time"
)

// Breakout
const (
	blank   = 0
	wall    = 1
	block   = 2
	paddle  = 3
	ball    = 4
	left    = -1
	right   = 1
	neutral = 0
)

type Location struct {
	Row int
	Col int
}

func main() {
	program := intcode.ReadIntcodeProgram("./input.txt")

	// Part 1
	func() {
		grid := make(map[Location]int)

		icc := intcode.NewIntCodeComputer(intcode.CopyIntcodeProgram(program))
		go icc.Run()

		for {
			select {
			case col := <-icc.OutputChannel:
				row := <-icc.OutputChannel
				id := <-icc.OutputChannel
				grid[Location{row, col}] = id
			case <-icc.DoneChannel:
				blocks := 0
				for _, v := range grid {
					if v == block {
						blocks++
					}
				}
				fmt.Printf("Part 1: %d\n", blocks)
				return
			}
		}
	}()

	// Part 2
	func() {
		p := intcode.CopyIntcodeProgram(program)

		// Command to start the game
		p[0] = 2

		grid := make(map[Location]int)

		icc := intcode.NewIntCodeComputer(p)
		icc.RequestInput = true
		go icc.Run()

		frames := 0
		ballCol := 0
		paddleCol := 0
		score := 0
		for {
			select {
			case <-icc.InputChannel:
				if ballCol > paddleCol {
					icc.InputChannel <- right
				} else if paddleCol > ballCol {
					icc.InputChannel <- left
				} else {
					icc.InputChannel <- neutral
				}
				frames++
			case col := <-icc.OutputChannel:
				row := <-icc.OutputChannel
				id := <-icc.OutputChannel

				if id == ball {
					ballCol = col
				}

				if id == paddle {
					paddleCol = col
				}

				if col == -1 && row == 0 {
					score = id
				} else {
					grid[Location{row, col}] = id
				}

				if frames > 0 {
					time.Sleep(5 * time.Millisecond)
					displayGrid(grid)
					fmt.Printf("\nScore: %d\n\n", score)
				}
			case <-icc.DoneChannel:
				fmt.Printf("Part 2: %d\n", score)
				return
			}
		}
	}()
}

func displayGrid(g map[Location]int) {
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

			switch {
			case !ok || id == blank:
				fmt.Printf("  ")
			case id == wall:
				fmt.Printf("🟥")
			case id == block:
				fmt.Printf("🟪")
			case id == paddle:
				fmt.Printf("🟩")
			case id == ball:
				fmt.Printf("⚽")
			}
		}
		fmt.Printf("\n")
	}
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
