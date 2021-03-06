package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

const (
	scaffold = 35
	space    = 46
	up       = 94
	down     = 62
	left     = 60
	right    = 118
	newline  = 10
)

func main() {
	input, _ := ioutil.ReadFile("./input.txt")
	var program []int
	json.Unmarshal([]byte("["+string(input)+"]"), &program)

	// Part 1
	func() {
		icc := NewIntCodeComputer(program)
		output := [][]string{}
		line := []string{}

		go icc.executeProgram()

	runProgram:
		for {
			select {
			case command := <-icc.OutputChannel:
				switch command {
				case scaffold:
					fmt.Print("#")
					line = append(line, "#")
				case space:
					fmt.Print(".")
					line = append(line, ".")
				case up:
					fmt.Print("^")
					line = append(line, "^")
				case down:
					fmt.Print("v")
					line = append(line, "v")
				case left:
					fmt.Print("<")
					line = append(line, "<")
				case right:
					fmt.Print(">")
					line = append(line, ">")
				case newline:
					fmt.Print("\n")
					if len(line) > 0 {
						output = append(output, line)
					}
					line = []string{}
				}
			case <-icc.DoneChannel:
				break runProgram
			}
		}
		fmt.Printf("output = %+v\n", output)

		alignmentParameters := findAlignmentParameters(output)

		sum := 0
		for _, pair := range alignmentParameters {
			sum += pair[0] * pair[1]
		}

		fmt.Printf("Part 1: = %+v\n", sum)

		findFullPath(output)

	}()
}

func findAlignmentParameters(output [][]string) [][]int {
	alignmentParameters := [][]int{}
	for row := 1; row < len(output)-1; row++ {
		for col := 1; col < len(output[row])-1; col++ {
			if output[row][col] == "#" &&
				output[row-1][col] == "#" &&
				output[row+1][col] == "#" &&
				output[row][col-1] == "#" &&
				output[row][col+1] == "#" {
				output[row][col] = "O"
				alignmentParameters = append(alignmentParameters, []int{row, col})
			}
		}
	}

	return alignmentParameters
}

func findFullPath(output [][]string) {
	// Find our position and headinng
	var positionRow int
	var positionCol int
	var heading string

	for row := range output {
		for col := range output[row] {
			switch output[row][col] {
			case "^", "v", "<", ">":
				positionRow = row
				positionCol = col
				heading = output[row][col]
			}
		}
	}
	fmt.Printf("positionRow = %+v\n", positionRow)
	fmt.Printf("positionCol = %+v\n", positionCol)
	fmt.Printf("heading = %+v\n", heading)
}

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

const (
	parameter = 0
	immediate = 1
	relative  = 2
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

type Param struct {
	Mode          int
	ValueOrOffset int
}

func (p *Param) Value(icc *IntCodeComputer) int {
	switch p.Mode {
	case parameter:
		return icc.MemGet(p.ValueOrOffset)
	case immediate:
		return p.ValueOrOffset
	case relative:
		return icc.MemGet(icc.RelBase + p.ValueOrOffset)
	default:
		panic("Unsupported mode")
	}
}

type IntCodeComputer struct {
	Name          string
	PhaseSetting  int
	Program       []int
	RelBase       int
	InputChannel  chan int
	OutputChannel chan int
	DoneChannel   chan bool
}

func NewIntCodeComputer(program []int) *IntCodeComputer {
	return &IntCodeComputer{
		Program:       program,
		InputChannel:  make(chan int),
		OutputChannel: make(chan int),
		DoneChannel:   make(chan bool),
	}
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
	case parameter:
		i = p.ValueOrOffset
	case relative:
		i = icc.RelBase + p.ValueOrOffset
	}

	if i >= len(icc.Program) {
		n := make([]int, i*3, i*3)
		copy(n, icc.Program)
		icc.Program = n
	}
	icc.Program[i] = v
}

func (icc *IntCodeComputer) executeProgram() {
	for programIndex := 0; programIndex < len(icc.Program); {
		value := icc.Program[programIndex]
		opcode, params := parseInstruction(value)

		for x := 0; x < opcodeArity[opcode]; x++ {
			params[x].ValueOrOffset = icc.Program[programIndex+x+1]
		}

		programIndex += opcodeArity[opcode] + 1

		switch opcode {
		case add:
			icc.MemSet(params[2], params[0].Value(icc)+params[1].Value(icc))
		case mult:
			icc.MemSet(params[2], params[0].Value(icc)*params[1].Value(icc))
		case jumpIfTrue:
			if params[0].Value(icc) != 0 {
				programIndex = params[1].Value(icc)
			}
		case jumpIfFalse:
			if params[0].Value(icc) == 0 {
				programIndex = params[1].Value(icc)
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
			icc.MemSet(params[0], <-icc.InputChannel)
		case output:
			icc.OutputChannel <- params[0].Value(icc)
		case halt:
			icc.DoneChannel <- true
			return
		}
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
