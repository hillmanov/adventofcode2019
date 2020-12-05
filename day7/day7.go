package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

// Intcode program lifted from day5.go

type intcodeProgram []int

const (
	add         = 1
	mult        = 2
	input       = 3
	output      = 4
	jumpIfTrue  = 5
	jumpIfFalse = 6
	lessThan    = 7
	equals      = 8
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
	halt:        0,
}

type Param struct {
	Mode  int
	Value int
}

type instruction struct {
	Opcode int
	Params []Param
}

type amplifier struct {
	Name            string
	PhaseSetting    int
	Program         intcodeProgram
	InputChannel    chan int
	OutputChannel   chan int
	OutputAmplifier *amplifier
	Halted          bool
}

func (a *amplifier) getValue(p Param) int {
	if p.Value > len(a.Program) {
		return 0
	}
	if p.Mode == 0 {
		return a.Program[p.Value]
	}
	return p.Value
}

func (a *amplifier) executeProgram(finalSignal chan int) {
	for i := 0; i < len(a.Program); {
		instruction := a.Program[i]

		opcode, params := parseInstruction(instruction)

		for x := 0; x < opcodeArity[opcode]; x++ {
			params[x].Value = a.Program[i+x+1]
		}

		i += opcodeArity[opcode] + 1

		switch opcode {
		case add:
			a.Program[params[2].Value] = a.getValue(params[0]) + a.getValue(params[1])
		case mult:
			a.Program[params[2].Value] = a.getValue(params[0]) * a.getValue(params[1])
		case jumpIfTrue:
			if a.getValue(params[0]) != 0 {
				i = a.getValue(params[1])
			}
		case jumpIfFalse:
			if a.getValue(params[0]) == 0 {
				i = a.getValue(params[1])
			}
		case lessThan:
			if a.getValue(params[0]) < a.getValue(params[1]) {
				a.Program[params[2].Value] = 1
			} else {
				a.Program[params[2].Value] = 0
			}
		case equals:
			if a.getValue(params[0]) == a.getValue(params[1]) {
				a.Program[params[2].Value] = 1
			} else {
				a.Program[params[2].Value] = 0
			}
		case input:
			a.Program[params[0].Value] = <-a.InputChannel
		case output:
			outputValue := a.getValue(params[0])
			if a.OutputAmplifier != nil && a.OutputAmplifier.Halted {
				finalSignal <- outputValue
				return
			}
			a.OutputChannel <- outputValue
		case halt:
			a.Halted = true
			return
		}
	}
}

func main() {
	input, _ := ioutil.ReadFile("./input.txt")
	var origP intcodeProgram
	json.Unmarshal([]byte("["+string(input)+"]"), &origP)

	A := &amplifier{Name: "A", InputChannel: make(chan int)}
	B := &amplifier{Name: "B", InputChannel: make(chan int)}
	C := &amplifier{Name: "C", InputChannel: make(chan int)}
	D := &amplifier{Name: "D", InputChannel: make(chan int)}
	E := &amplifier{Name: "E", InputChannel: make(chan int)}
	amplifiers := []*amplifier{A, B, C, D, E}

	// Part 1
	func() {
		// Normal execution phase settings
		phases := []int{0, 1, 2, 3, 4}
		finalSignal := make(chan int)

		maxSignal := 0
		for _, phase := range permutation(phases) {
			for i, phaseSetting := range phase {
				p := make([]int, len(origP))
				copy(p, origP)
				amplifiers[i].Program = p
				amplifiers[i].Halted = false

				A.OutputChannel = B.InputChannel
				A.OutputAmplifier = B
				B.OutputChannel = C.InputChannel
				B.OutputAmplifier = C
				C.OutputChannel = D.InputChannel
				C.OutputAmplifier = D
				D.OutputChannel = E.InputChannel
				D.OutputAmplifier = E
				E.OutputChannel = finalSignal
				E.OutputAmplifier = nil

				go amplifiers[i].executeProgram(finalSignal)
				amplifiers[i].InputChannel <- phaseSetting
			}

			A.InputChannel <- 0
			maxSignal = max(maxSignal, <-finalSignal)
		}

		fmt.Printf("Part 1: Max Signal = %+v\n", maxSignal)
	}()

	// Part 2
	func() {
		// Feedback loop phase settings
		phases := []int{5, 6, 7, 8, 9}
		finalSignal := make(chan int)

		maxSignal := 0
		for _, phase := range permutation(phases) {
			for i, phaseSetting := range phase {
				p := make([]int, len(origP))
				copy(p, origP)
				amplifiers[i].Program = p
				amplifiers[i].Halted = false

				A.OutputChannel = B.InputChannel
				A.OutputAmplifier = B
				B.OutputChannel = C.InputChannel
				B.OutputAmplifier = C
				C.OutputChannel = D.InputChannel
				C.OutputAmplifier = D
				D.OutputChannel = E.InputChannel
				D.OutputAmplifier = E
				E.OutputChannel = A.InputChannel
				E.OutputAmplifier = A

				go amplifiers[i].executeProgram(finalSignal)
				amplifiers[i].InputChannel <- phaseSetting
			}

			A.InputChannel <- 0
			maxSignal = max(maxSignal, <-finalSignal)
		}

		fmt.Printf("Part 2: Max Signal = %+v\n", maxSignal)
	}()
}

func parseInstruction(instruction int) (int, []Param) {
	s := strconv.Itoa(instruction)
	if len(s) == 1 {
		return instruction, make([]Param, opcodeArity[instruction])
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
