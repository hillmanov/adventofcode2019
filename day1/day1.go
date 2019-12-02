package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}

	totalFuelRequirements := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		mass, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}
		totalFuelRequirements = totalFuelRequirements + getRequiredFuel(mass)
	}

	fmt.Printf("Total = %+v\n", totalFuelRequirements)
}

func getRequiredFuel(mass int) int {
	return int(mass/3) - 2
}
