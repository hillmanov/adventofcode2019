package main

import (
	"adventofcode/intcode"
	"fmt"
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
	intcodeProgram := intcode.ReadIntcodeProgram("./input.txt")

	r := part1(intcode.CopyIntcodeProgram(intcodeProgram))
	fmt.Printf("Result: = %d\n", r[0])

	noun, verb := part2(intcode.CopyIntcodeProgram(intcodeProgram))
	fmt.Printf("Noun: %d Verb: %d Computed: %d\n", noun, verb, 100*noun+verb)
}

func part1(p intcode.IntcodeProgram) intcodeProgram {
	p[1] = 12
	p[2] = 2

	icc := intcode.NewIntCodeComputer(p)

	go icc.Run()
	<-icc.DoneChannel
	return icc.Program
}

func part2(p intcode.IntcodeProgram) (int, int) {
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			candidate := intcode.CopyIntcodeProgram(p)
			candidate[1] = noun
			candidate[2] = verb
			icc := intcode.NewIntCodeComputer(candidate)
			go icc.Run()
			<-icc.DoneChannel
			if icc.Program[0] == 19690720 {
				return noun, verb
			}
		}
	}
	return -1, -1
}
