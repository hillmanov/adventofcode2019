package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

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

func main() {
	input, _ := ioutil.ReadFile("./input.txt")
	var origP intcodeProgram
	json.Unmarshal([]byte("["+string(input)+"]"), &origP)

	p := make([]int, len(origP))

	fmt.Println("Part 1")
	copy(p, origP)
	executeIntcodeProgram(p, 1)

	fmt.Println("Part 2")
	copy(p, origP)
	executeIntcodeProgram(p, 5)
}

func executeIntcodeProgram(ip intcodeProgram, inputValue int) intcodeProgram {
	for i := 0; i < len(ip); {
		value := ip[i]

		opcode, params := parseInstruction(value)

		for x := 0; x < opcodeArity[opcode]; x++ {
			params[x].Value = ip[i+x+1]
		}
		// This might be overwritten
		i += opcodeArity[opcode] + 1

		switch opcode {
		case add:
			ip[params[2].Value] = getValue(params[0], ip) + getValue(params[1], ip)
		case mult:
			ip[params[2].Value] = getValue(params[0], ip) * getValue(params[1], ip)
		case jumpIfTrue:
			if getValue(params[0], ip) != 0 {
				i = getValue(params[1], ip)
			}
		case jumpIfFalse:
			if getValue(params[0], ip) == 0 {
				i = getValue(params[1], ip)
			}
		case lessThan:
			if getValue(params[0], ip) < getValue(params[1], ip) {
				ip[params[2].Value] = 1
			} else {
				ip[params[2].Value] = 0
			}
		case equals:
			if getValue(params[0], ip) == getValue(params[1], ip) {
				ip[params[2].Value] = 1
			} else {
				ip[params[2].Value] = 0
			}
		case input:
			ip[params[0].Value] = inputValue
		case output:
		case halt:
			return ip
		}
	}

	return ip
}

func getValue(p Param, ip intcodeProgram) int {
	if p.Mode == 0 {
		return ip[p.Value]
	}
	return p.Value
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
