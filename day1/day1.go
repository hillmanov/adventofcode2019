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

	fuelForMass := 0
  fuelForMassAndFuel := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		mass, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}
		fuelForMass += getFuelForMass(mass)
    fuelForMassAndFuel += getFuelForMassAndFuel(mass)
	}

	fmt.Printf("Fuel for mass: = %+v\n", fuelForMass)
	fmt.Printf("Fuel for mass and fuel: = %+v\n", fuelForMassAndFuel)
}

func getFuelForMassAndFuel(mass int) int {
	fuel := getFuelForMass(mass)
	if fuel > 0 {
		fuel += getFuelForMassAndFuel(fuel)
		return fuel
	}
	return 0
}

func getFuelForMass(mass int) int {
	return int(mass/3) - 2
}
