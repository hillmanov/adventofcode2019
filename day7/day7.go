package main

import (
	"adventofcode/intcode"
	"fmt"
)

func main() {
	input := intcode.ReadIntcodeProgram("./input.txt")

	// Part 1
	func() {
		A := intcode.NewIntCodeComputer(nil)
		B := intcode.NewIntCodeComputer(nil)
		C := intcode.NewIntCodeComputer(nil)
		D := intcode.NewIntCodeComputer(nil)
		E := intcode.NewIntCodeComputer(nil)

		amplifiers := []*intcode.IntCodeComputer{A, B, C, D, E}

		// Normal execution phase settings
		phases := []int{0, 1, 2, 3, 4}
		finalSignal := make(chan int)

		maxSignal := 0
		for _, phase := range permutation(phases) {
			for i, phaseSetting := range phase {
				amplifiers[i].Program = intcode.CopyIntcodeProgram(input)

				A.OutputChannel = B.InputChannel
				B.OutputChannel = C.InputChannel
				C.OutputChannel = D.InputChannel
				D.OutputChannel = E.InputChannel
				E.OutputChannel = finalSignal

				go amplifiers[i].Run()
				amplifiers[i].InputChannel <- phaseSetting
			}

			A.InputChannel <- 0
			maxSignal = max(maxSignal, <-finalSignal)
		}

		fmt.Printf("Part 1: Max Signal = %+v\n", maxSignal)
	}()

	// Part 2
	func() {
		A := intcode.NewIntCodeComputer(nil)
		B := intcode.NewIntCodeComputer(nil)
		C := intcode.NewIntCodeComputer(nil)
		D := intcode.NewIntCodeComputer(nil)
		E := intcode.NewIntCodeComputer(nil)
		amplifiers := []*intcode.IntCodeComputer{A, B, C, D, E}

		// Feedback loop phase settings
		phases := []int{5, 6, 7, 8, 9}
		finalSignal := make(chan int)

		maxSignal := 0
		for _, phase := range permutation(phases) {
			for i, phaseSetting := range phase {
				amplifiers[i].Program = intcode.CopyIntcodeProgram(input)

				A.OutputChannel = B.InputChannel
				B.OutputChannel = C.InputChannel
				C.OutputChannel = D.InputChannel
				D.OutputChannel = E.InputChannel
				E.OutputChannel = A.InputChannel

				go amplifiers[i].Run()
				amplifiers[i].InputChannel <- phaseSetting
			}

			A.InputChannel <- 0

		run:
			for {
				select {
				case <-A.DoneChannel:
					E.OutputChannel = finalSignal
				case signal := <-finalSignal:
					maxSignal = max(maxSignal, signal)
				case <-E.DoneChannel:
					break run
				}
			}
		}

		fmt.Printf("Part 2: Max Signal = %+v\n", maxSignal)
	}()
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
