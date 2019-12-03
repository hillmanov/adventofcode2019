package main

import (
	"fmt"
	"testing"
)

type spec struct {
	Input    int
	Expected int
}

func TestFuelForMassDay1(t *testing.T) {
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
			actual := getFuelForMass(spec.Input)
			if actual != spec.Expected {
				t.Errorf("Input: %d. Expected: %d. Actual: %d", spec.Input, spec.Expected, actual)
			}
		})
	}
}

func TestFuelForMassAndFuelDay1(t *testing.T) {
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
			Expected: 966,
		},
		spec{
			Input:    100756,
			Expected: 50346,
		},
	}

	for _, spec := range specs {
		t.Run(fmt.Sprintf("TestInput%d", spec.Input), func(t *testing.T) {
			actual := getFuelForMassAndFuel(spec.Input)
			if actual != spec.Expected {
				t.Errorf("Input: %d. Expected: %d. Actual: %d", spec.Input, spec.Expected, actual)
			}
		})
	}

}
