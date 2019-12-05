package main

import (
	"testing"
)

func TestWireBounds(t *testing.T) {
	w := wire{
		id: 1,
		path: []instruction{
			instruction{
				Direction: "U",
				Steps:     3,
			},
			instruction{
				Direction: "D",
				Steps:     5,
			},
			instruction{
				Direction: "L",
				Steps:     10,
			},
			instruction{
				Direction: "R",
				Steps:     15,
			},
		},
	}

	up, down, left, right := w.bounds()
	if up != 3 {
		t.Errorf("Expected up to be 3, got %d", up)
	}
	if down != -2 {
		t.Errorf("Expected down to be -2, got %d", down)
	}
	if left != -10 {
		t.Errorf("Expected up to be 3, got %d", left)
	}
	if right != 5 {
		t.Errorf("Expected up to be 3, got %d", right)
	}
}

func TestTrace(t *testing.T) {
	w := wire{
		id: 1,
		path: []instruction{
			instruction{
				Direction: "U",
				Steps:     3,
			},
			instruction{
				Direction: "R",
				Steps:     5,
			},
			instruction{
				Direction: "D",
				Steps:     10,
			},
			instruction{
				Direction: "L",
				Steps:     15,
			},
		},
	}

	up, down, left, right := w.bounds()

	cols := abs(left) + abs(right)
	rows := abs(up) + abs(down)

	startCol := abs(left) - 1
	startRow := abs(down)

	g := makeGrid(cols, rows)

	trace(g, w, startCol, startRow)

	g[startRow][startCol] = 7

	// for i, _ := range g {
	// 	fmt.Println(g[len(g)-1-i])
	// }
}

func TestEverything(t *testing.T) {

	t.Run("Test1", func(t *testing.T) {
		path1 := "R75,D30,R83,U83,L12,D49,R71,U7,L72"
		path2 := "U62,R66,U55,R34,D71,R55,D58,R83"

		dist := getIntersections([]string{path1, path2})
		if dist != 159 {
			t.Errorf("Expected 159, got %d", dist)
		}
	})

	t.Run("Test2", func(t *testing.T) {
		path1 := "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51"
		path2 := "U98,R91,D20,R16,D67,R40,U7,R15,U6,R7"

		dist := getIntersections([]string{path1, path2})
		if dist != 135 {
			t.Errorf("Expected 135, got %d", dist)
		}
	})
}
