package main

import (
	"adventofcode/intcode"
	"fmt"
)

func main() {
	intcodeProgram := intcode.ReadIntcodeProgram("./input.txt")

	// Part 1
	icc := intcode.NewIntCodeComputer(intcode.CopyIntcodeProgram(intcodeProgram))
	go icc.Run()
	icc.InputChannel <- 1
	var finalOutput int

run:
	for {
		select {
		case o := <-icc.OutputChannel:
			finalOutput = o
		case <-icc.DoneChannel:
			fmt.Printf("Part 1: %+v\n", finalOutput)
			break run
		}
	}

	// Part 2
	icc = intcode.NewIntCodeComputer(intcode.CopyIntcodeProgram(intcodeProgram))
	go icc.Run()
	icc.InputChannel <- 5
	fmt.Printf("Part 2: %+v\n", <-icc.OutputChannel)
}
