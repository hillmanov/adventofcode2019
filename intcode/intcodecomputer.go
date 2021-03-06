package intcode

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

type IntcodeProgram []int

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
	RequestInput  bool
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

func (icc *IntCodeComputer) Run() {
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
			// Check to see if anyone is waiting
			if icc.RequestInput {
				icc.InputChannel <- 0
			}
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

func ReadIntcodeProgram(filename string) IntcodeProgram {
	line, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var icp IntcodeProgram
	if err := json.Unmarshal([]byte("["+string(line)+"]"), &icp); err != nil {
		panic(err)
	}
	return icp
}

func CopyIntcodeProgram(source IntcodeProgram) IntcodeProgram {
	target := make([]int, len(source))
	for i, v := range source {
		target[i] = v
	}
	return target
}
