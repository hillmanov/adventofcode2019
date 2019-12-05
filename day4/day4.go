package main

import (
	"fmt"
	"strconv"
)

func main() {
	part1Matches := 0
	part2Matches := 0

	for i := 136818; i < 685979; i++ {
		d1, d2, d3, d4, d5, d6 := getDigits(i)

		// Make sure they are always increasing
		if d1 > d2 || d2 > d3 || d3 > d4 || d4 > d5 || d5 > d6 {
			continue
		}

		// Make sure there is at least 1 group of matching digits
		if d1 != d2 && d2 != d3 && d3 != d4 && d4 != d5 && d5 != d6 {
			continue
		}
		part1Matches++

		// Make sure there is an isolated 2 digit group
		if !((d1 == d2 && d2 != d3) || (d2 == d3 && d1 != d2 && d3 != d4) || (d3 == d4 && d2 != d3 && d4 != d5) || (d4 == d5 && d3 != d4 && d5 != d6) || (d5 == d6 && d5 != d4)) {
			continue
		}
		part2Matches++

	}
	fmt.Printf("Part 1 matches: %d\n", part1Matches)
	fmt.Printf("Part 2 matches: %d\n", part2Matches)
}

func getDigits(num int) (int, int, int, int, int, int) {
	var digits []int
	s := strconv.Itoa(num)

	for _, sDigit := range s {
		nDigit, _ := strconv.Atoi(string(sDigit))
		digits = append(digits, nDigit)
	}

	return digits[0], digits[1], digits[2], digits[3], digits[4], digits[5]
}
