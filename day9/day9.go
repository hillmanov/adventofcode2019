package main

import (
	"adventofcode/intcode"
	"fmt"
)

func main() {
	program := intcode.ReadIntcodeProgram("./input.txt")

	// Part 1
	func() {
		icc := intcode.NewIntCodeComputer(program)
		go icc.Run()
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
		icc := intcode.NewIntCodeComputer(program)
		go icc.Run()
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
