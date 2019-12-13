package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
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
	PhaseSetting  int
	Program       []int
	RelBase       int
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
		icc := &IntCodeComputer{
			Name:          "IntCodeComputer",
			Program:       origP,
			InputChannel:  make(chan int),
			OutputChannel: make(chan int),
			DoneChannel:   make(chan bool),
		}

		go icc.executeProgram()
		icc.InputChannel <- 1

		for {
			select {
			case o := <-icc.OutputChannel:
				fmt.Printf("Part 1: %d\n", o)
			case <-icc.DoneChannel:
				return
			}
		}
	}()

	// Part 2
	func() {
		icc := &IntCodeComputer{
			Name:          "IntCodeComputer",
			Program:       origP,
			InputChannel:  make(chan int),
			OutputChannel: make(chan int),
			DoneChannel:   make(chan bool),
		}

		go icc.executeProgram()
		icc.InputChannel <- 2

		for {
			select {
			case o := <-icc.OutputChannel:
				fmt.Printf("Part 2: %d\n", o)
			case <-icc.DoneChannel:
				return
			}
		}
	}()
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
func permutation(xs []int) (permuts [][]int) {
	var rc func([]int, int)
	rc = func(a []int, k int) {
		if k == len(a) {
			permuts = append(permuts, append([]int{}, a...))
		} else {
			for i := k; i < len(xs); i++ {
				a[k], a[i] = a[i], a[k]
				rc(a, k+1)
				a[k], a[i] = a[i], a[k]
			}
		}
	}
	rc(xs, 0)

	return permuts
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
