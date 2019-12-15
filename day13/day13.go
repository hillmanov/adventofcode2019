package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
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

const (
	add         = 1
	mult        = 2
	input       = 3
	output      = 4
	jumpIfTrue  = 5
	jumpIfFalse = 6
	lessThan    = 7
	equals      = 8
	setRelBase  = 9
	halt        = 99
)

var opcodeArity = map[int]int{
	add:         3,
	mult:        3,
	input:       1,
	output:      1,
	jumpIfTrue:  2,
	jumpIfFalse: 2,
	lessThan:    3,
	equals:      3,
	setRelBase:  1,
	halt:        0,
}

type Location struct {
	Row int
	Col int
}

type Param struct {
	Mode          int
	ValueOrOffset int
}

func (p *Param) Value(icc *IntCodeComputer) int {
	switch p.Mode {
	case 0:
		return icc.MemGet(p.ValueOrOffset)
	case 1:
		return p.ValueOrOffset
	case 2:
		return icc.MemGet(icc.RelBase + p.ValueOrOffset)
	default:
		panic("Should not get here")
	}
}

type IntCodeComputer struct {
	Name          string
	Program       []int
	RelBase       int
	InputRequest  chan bool
	InputChannel  chan int
	OutputChannel chan int
	DoneChannel   chan bool
}

func (icc *IntCodeComputer) MemGet(i int) int {
	if i > len(icc.Program) {
		return 0
	}
	return icc.Program[i]
}

func (icc *IntCodeComputer) MemSet(p Param, v int) {
	var i int
	switch p.Mode {
	case 0:
		i = p.ValueOrOffset
	case 2:
		i = icc.RelBase + p.ValueOrOffset
	}

	if i > len(icc.Program) {
		n := make([]int, i*2, i*2)
		copy(n, icc.Program)
		icc.Program = n
	}
	icc.Program[i] = v
}

func (icc *IntCodeComputer) executeProgram() {
	for i := 0; i < len(icc.Program); {
		value := icc.Program[i]
		opcode, params := parseInstruction(value)

		for x := 0; x < opcodeArity[opcode]; x++ {
			params[x].ValueOrOffset = icc.Program[i+x+1]
		}
		// This might be overwritten
		i += opcodeArity[opcode] + 1

		switch opcode {
		case add:
			icc.MemSet(params[2], params[0].Value(icc)+params[1].Value(icc))
		case mult:
			icc.MemSet(params[2], params[0].Value(icc)*params[1].Value(icc))
		case jumpIfTrue:
			if params[0].Value(icc) != 0 {
				i = params[1].Value(icc)
			}
		case jumpIfFalse:
			if params[0].Value(icc) == 0 {
				i = params[1].Value(icc)
			}
		case lessThan:
			if params[0].Value(icc) < params[1].Value(icc) {
				icc.MemSet(params[2], 1)
			} else {
				icc.MemSet(params[2], 0)
			}
		case equals:
			if params[0].Value(icc) == params[1].Value(icc) {
				icc.MemSet(params[2], 1)
			} else {
				icc.MemSet(params[2], 0)
			}
		case setRelBase:
			icc.RelBase += params[0].Value(icc)
		case input:
			icc.InputRequest <- true
			icc.MemSet(params[0], <-icc.InputChannel)
		case output:
			icc.OutputChannel <- params[0].Value(icc)
		case halt:
			icc.DoneChannel <- true
			return
		}
	}
}

func main() {
	input, _ := ioutil.ReadFile("./input.txt")
	var origP []int
	json.Unmarshal([]byte("["+string(input)+"]"), &origP)

	// Part 1
	func() {
		p := make([]int, len(origP))
		copy(p, origP)
		icc := &IntCodeComputer{
			Name:          "IntCodeComputer",
			Program:       p,
			InputRequest:  make(chan bool),
			InputChannel:  make(chan int),
			OutputChannel: make(chan int),
			DoneChannel:   make(chan bool),
		}

		grid := make(map[Location]int)

		go icc.executeProgram()

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
		p := make([]int, len(origP))
		copy(p, origP)
		icc := &IntCodeComputer{
			Name:          "IntCodeComputer",
			Program:       p,
			InputRequest:  make(chan bool),
			InputChannel:  make(chan int),
			OutputChannel: make(chan int),
			DoneChannel:   make(chan bool),
		}

		// Command to start the game
		p[0] = 2

		grid := make(map[Location]int)

		go icc.executeProgram()

		frames := 0
		ballCol := 0
		paddleCol := 0
		score := 0
		for {
			select {
			case <-icc.InputRequest:
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

func parseInstruction(value int) (int, []Param) {
	s := strconv.Itoa(value)
	if len(s) == 1 {
		return value, make([]Param, opcodeArity[value])
	}

	// Opcode is the last two digits of the instruction
	opcode, _ := strconv.Atoi(string(s[len(s)-2:]))

	// Params are the first N-2 digits of the instruction in reverse order
	paramModes := reverse(s[:len(s)-2])

	// We only gather the parameter modes at this point. Values will be gathered later.
	params := make([]Param, opcodeArity[opcode])
	for i, v := range paramModes {
		m, _ := strconv.Atoi(string(v))
		params[i].Mode = m
	}

	return opcode, params
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
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
