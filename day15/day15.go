package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
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
		var oxygenTankLocation Location

		go icc.executeProgram()
		path := make([]int, 0)
		currentLocation := Location{Row: 0, Col: 0}
		direction := north
		stepsToTank := 0

		for direction != -2 {
			select {
			case <-icc.InputRequest:
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

		// displayGrid(grid, currentLocation)
		fmt.Printf("Part 1 (Steps to tank): %d\n", stepsToTank)

		currentLocation = oxygenTankLocation

		minutes := 0
		var oxygenate func(l Location, minutes int) int
		oxygenate = func(l Location, currentMinutes int) int {
			switch grid[l] {
			case wall, oxygenated:
				return currentMinutes - 1
			case moved, tank:
				grid[l] = oxygenated
				// displayGrid(grid, l)
				// time.Sleep(1 * time.Millisecond)
				// fmt.Printf("minutes = %+v\n", minutes)
				for _, oxygenateLocation := range []Location{applyDirection(l, north), applyDirection(l, south), applyDirection(l, east), applyDirection(l, west)} {
					minutes = max(minutes, oxygenate(oxygenateLocation, currentMinutes+1))
				}
			}
			return minutes
		}

		minutes = oxygenate(oxygenTankLocation, 0)
		fmt.Printf("Part 2 (Minutes to full oxygenation): %+v\n", minutes)
	}()
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
