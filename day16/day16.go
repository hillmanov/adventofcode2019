package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	input, _ := ioutil.ReadFile("./input.txt")
	originalSignal := make([]int, 0)
	basePattern := []int{0, 1, 0, -1}

	for _, sDigit := range strings.Trim(string(input), "\n") {
		digit, _ := strconv.Atoi(string(sDigit))
		originalSignal = append(originalSignal, digit)
	}

	verifySignal := make([]int, len(originalSignal))
	copy(verifySignal, originalSignal)

	for i := 0; i < 100; i++ {
		verifySignal = runPhase(verifySignal, basePattern)
	}

	output := ""
	for _, digit := range verifySignal[:8] {
		output += fmt.Sprintf("%d", digit)
	}

	fmt.Printf("Part 1: %+v\n", output)

	//	Part 2
	realSignal := make([]int, 0)
	for i := 0; i < 10000; i++ {
		realSignal = append(realSignal, originalSignal...)
	}

	offset := realSignal[0]*1000000 + realSignal[1]*100000 + realSignal[2]*10000 + realSignal[3]*1000 + realSignal[4]*100 + realSignal[5]*10 + realSignal[6]

	verifySignal = make([]int, len(realSignal))
	copy(verifySignal, realSignal)
	for i := 0; i < 100; i++ {
		verifySignal[len(realSignal)-1] = realSignal[len(realSignal)-1]
		for j := len(realSignal) - 2; j >= 0; j-- {
			verifySignal[j] = (realSignal[j] + verifySignal[j+1]) % 10
		}

		realSignal = verifySignal
	}

	fmt.Printf("realSignal[:8] = %+v\n", realSignal[offset:offset+8])
}

func runPhase(signal []int, basePattern []int) []int {
	newSignal := make([]int, len(signal))
	for i := 0; i < len(signal); i++ {
		iterationPattern := getPatternForPhase(basePattern, i, len(signal))
		iterationSum := 0
		for j := 0; j < len(signal); j++ {
			iterationSum += signal[j] * iterationPattern[j]
		}
		newSignal[i] = abs(iterationSum) % 10
	}

	return newSignal
}

func getPatternForPhase(basePattern []int, iteration int, length int) []int {
	iterationPattern := make([]int, length+1)
	for i := 0; i < length+1; i++ {
		iterationPattern[i] = basePattern[((i / (iteration + 1)) % len(basePattern))]
	}
	return iterationPattern[1:]
}

func abs(a int) int {
	if a < 0 {
		return a * -1
	}
	return a
}
