package main

import (
	"fmt"
	"testing"
)

type spec struct {
	Input    []int
	Expected []int
}

func TestGravityAssistDay2(t *testing.T) {
	specs := []spec{
		spec{
			Input:    []int{1, 0, 0, 0, 99},
			Expected: []int{2, 0, 0, 0, 99},
		},

		spec{
			Input:    []int{2, 3, 0, 3, 99},
			Expected: []int{2, 3, 0, 6, 99},
		},

		spec{
			Input:    []int{2, 4, 4, 5, 99, 0},
			Expected: []int{2, 4, 4, 5, 99, 9801},
		},

		spec{
			Input:    []int{1, 1, 1, 4, 99, 5, 6, 0, 99},
			Expected: []int{30, 1, 1, 4, 2, 5, 6, 0, 99},
		},
	}

	for i, spec := range specs {
		t.Run(fmt.Sprintf("TestInput%d", i), func(t *testing.T) {
			actual := executeIntcodeProgram(spec.Input)
			if !Equal(actual, spec.Expected) {
				t.Errorf("Input: %d. Expected: %d. Actual: %d", spec.Input, spec.Expected, actual)
			}
		})
	}
}

func Equal(a, b intcodeProgram) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
