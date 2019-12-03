package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type intcodeProgram []int

const (
	add  = 1
	mult = 2
	halt = 99
)

type instruction struct {
	Opcode      int
	Param1      int
	Param2      int
	Destination int
}

func main() {
	line, err := ioutil.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}

	var origP intcodeProgram
	if err := json.Unmarshal([]byte("["+string(line)+"]"), &origP); err != nil {
		panic(err)
	}

	p := make([]int, len(origP))

	// Part1
	copy(p, origP)
	r := part1(p)
	fmt.Printf("Result: = %d\n", r[0])

	// Part2
	copy(p, origP)
	noun, verb := part2(p)
	fmt.Printf("Noun: %d Verb: %d Computed: %d\n", noun, verb, 100*noun+verb)
}

func part1(p intcodeProgram) intcodeProgram {
	// Initialize according to instructions
	p[1] = 12
	p[2] = 2

	return executeIntcodeProgram(p)
}

func part2(p intcodeProgram) (int, int) {
	candidateP := make([]int, len(p))
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			copy(candidateP, p)
			candidateP[1] = noun
			candidateP[2] = verb
			r := executeIntcodeProgram(candidateP)
			if r[0] == 19690720 {
				return noun, verb
			}
		}
	}
	return -1, -1
}

func executeIntcodeProgram(p intcodeProgram) intcodeProgram {
	var f instruction
	for i := 0; f.Opcode != halt; i += 4 {
		if len(p[i:]) >= 4 {
			f = instruction{
				Opcode:      p[i],
				Param1:      p[i+1],
				Param2:      p[i+2],
				Destination: p[i+3],
			}
		} else {
			f = instruction{
				Opcode: p[i],
			}
		}

		switch f.Opcode {
		case add:
			p[f.Destination] = p[f.Param1] + p[f.Param2]
		case mult:
			p[f.Destination] = p[f.Param1] * p[f.Param2]
		case halt:
			break
		}
	}

	return p
}
