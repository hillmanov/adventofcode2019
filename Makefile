define GO_MOD_TEMPLATE
module adventofcode/day${day}

go 1.15

replace adventofcode/utils => ../utils

require adventofcode/utils v0.0.0
endef


define GO_FILE_TEMPLATE
package main

import (
	"adventofcode/utils"
	"fmt"
)

func main() {

	entries, err := utils.ReadInts("./input.txt")
	if err != nil {
		panic(err)
	}

	fmt.Printf("entries = %+v\\n", entries)

	part1Solution, err := part1()
	if err != nil {
		panic(err)
	}

	part2Solution, err := part2()
	if err != nil {
		panic(err)
	}

  fmt.Printf("Day ${day}: Part 1: = %+v\\n", part1Solution)
	fmt.Printf("Day ${day}: Part 2: = %+v\\n", part2Solution)
}

func part1() (int, error) {
	return -1, nil
}

func part2() (int, error) {
	return -1, nil
}

endef

export GO_MOD_TEMPLATE
export GO_FILE_TEMPLATE

init:
	@mkdir day${day}
	@echo "$$GO_MOD_TEMPLATE" > day${day}/go.mod
	@echo "$$GO_FILE_TEMPLATE" > day${day}/day${day}.go
	@touch day${day}/input.txt
	@touch day${day}/README.md

