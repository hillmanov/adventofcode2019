package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Input struct {
	Chemical string
	Amount   int
}

type Output struct {
	Chemical string
	Amount   int
}

type Reaction struct {
	Inputs []Input
	Output Output
}

func main() {
	input, _ := ioutil.ReadFile("./input.txt")
	reactions := parseReactions(string(input))

	reactionsByOuputChemical := make(map[string]Reaction)
	for _, reaction := range reactions {
		reactionsByOuputChemical[reaction.Output.Chemical] = reaction
	}

	ore := "ORE"
	target := Output{
		Amount:   1,
		Chemical: "FUEL",
	}

	amountAvailable := make(map[string]int)

	var calculateOreAmount func(chemical string, amountNeeded int) int
	calculateOreAmount = func(chemical string, amountNeeded int) int {
		if _, ok := amountAvailable[chemical]; !ok {
			amountAvailable[chemical] = 0
		}

		if amountAvailable[chemical] >= amountNeeded {
			amountAvailable[chemical] = amountAvailable[chemical] - amountNeeded
			return 0
		}

		if chemical == ore {
			return amountNeeded
		}

		producingReaction := reactionsByOuputChemical[chemical]
		baseReactionAmount := producingReaction.Output.Amount

		multiplier := int(math.Max(math.Ceil((float64(amountNeeded-amountAvailable[chemical]) / float64(baseReactionAmount))), 1))

		amountAvailable[chemical] = amountAvailable[chemical] + (producingReaction.Output.Amount * multiplier) - amountNeeded

		oreNeeded := 0
		for _, input := range producingReaction.Inputs {
			oreNeeded += calculateOreAmount(input.Chemical, input.Amount*multiplier)
		}
		return oreNeeded
	}

	oreNeeded := calculateOreAmount(target.Chemical, target.Amount)
	fmt.Printf("Part 1: %d\n", oreNeeded)

	// Part 2
	// I did the "bottom up" approach for part 1, which makes part 2 trickier.
	// Taking a more intelligent brute force approach for part 2.
	// Will try to zero in on the exact amount of fuel needed
	oreNeeded = 0
	oreAvailable := 1000000000000
	fuelAmount := 1

	delta := oreAvailable / 10
	direction := 1

	for delta != 0 {
		// Reset inbetween runs
		amountAvailable = make(map[string]int)
		oreNeeded = calculateOreAmount("FUEL", fuelAmount)
		if oreNeeded < oreAvailable {
			if direction == -1 {
				direction = 1
				delta = delta / 10
			}
		}

		if oreNeeded > oreAvailable {
			if direction == 1 {
				direction = -1
				delta = delta / 10
			}
		}

		fuelAmount += (delta * direction)
	}
	fmt.Printf("Part 2: %d\n", fuelAmount)
}

func parseReactions(data string) []Reaction {
	amountAndChemicalRe := regexp.MustCompile(`(\d+)\s(\w+)`)
	lines := strings.Split(data, "\n")

	reactions := make([]Reaction, 0)
	for _, line := range lines {
		matches := amountAndChemicalRe.FindAllStringSubmatch(line, -1)
		if len(matches) == 0 {
			break
		}

		inputMatches := matches[:len(matches)-1]
		outputMatch := matches[len(matches)-1]

		inputs := make([]Input, 0)
		for _, inputMatch := range inputMatches {
			amount, _ := strconv.Atoi(inputMatch[1])
			chemical := inputMatch[2]
			inputs = append(inputs, Input{Chemical: chemical, Amount: amount})
		}

		outputAmount, _ := strconv.Atoi(outputMatch[1])
		outputChemical := outputMatch[2]

		reactions = append(reactions, Reaction{
			Inputs: inputs,
			Output: Output{Chemical: outputChemical, Amount: outputAmount},
		})
	}

	return reactions
}
