package main

import (
	"fmt"
	"testing"
)

type spec struct {
	Input    int
	Expected int
}

func TestFuelRequirement(t *testing.T) {
	specs := []spec{
		spec{
			Input:    12,
			Expected: 2,
		},
		spec{
			Input:    14,
			Expected: 2,
		},
		spec{
			Input:    1969,
			Expected: 654,
		},
		spec{
			Input:    100756,
			Expected: 33583,
		},
	}

	for _, spec := range specs {
		t.Run(fmt.Sprintf("TestInput%d", spec.Input), func(t *testing.T) {
			actual := getRequiredFuel(spec.Input)
			if actual != spec.Expected {
				t.Errorf("Input: %d. Expected: %d. Actual: %d", spec.Input, spec.Expected, actual)
			}
		})
	}

}
