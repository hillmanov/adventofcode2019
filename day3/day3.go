package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type instruction struct {
	Direction string
	Steps     int
}

type wire struct {
	id   uint8
	path []instruction
}

func (w wire) bounds() (up, down, left, right int) {
	col := 0
	row := 0
	for _, instruction := range w.path {
		switch instruction.Direction {
		case "U":
			row += instruction.Steps
			up = max(up, row)
		case "D":
			row -= instruction.Steps
			down = min(down, row)
		case "L":
			col -= instruction.Steps
			left = min(left, col)
		case "R":
			col += instruction.Steps
			right = max(right, col)
		}
	}
	return
}

func main() {
	contents, _ := ioutil.ReadFile("./input.txt")
	paths := strings.Split(string(contents), "\n")

	dist := getDistanceOfClosesIntersection(paths)
	fmt.Printf("Manhattan distance of closet intersection: %d\n", dist)
}

func getDistanceOfClosesIntersection(paths []string) int {
	wire1 := wire{
		id:   1,
		path: parsePath(paths[0]),
	}

	wire2 := wire{
		id:   2,
		path: parsePath(paths[1]),
	}

	up, down, left, right := maxBounds(wire1, wire2)

	cols := abs(left) + abs(right)
	rows := abs(up) + abs(down)

	g := makeGrid(cols, rows)

	startCol, startRow := getCentralPort(wire1, wire2)

	trace(g, wire1, startCol, startRow)
	trace(g, wire2, startCol, startRow)

	intersections := findIntersections(g, startCol, startRow, wire1, wire2)

	sortIntersections(startCol, startRow, intersections)

	closest := intersections[0]

	return dist(startCol, startRow, closest[0], closest[1])
}

func trace(g [][]uint8, w wire, startCol, startRow int) {
	col := startCol
	row := startRow
	for _, instruction := range w.path {
		switch instruction.Direction {
		case "U":
			for i := 0; i < instruction.Steps; i++ {
				g[row][col] |= w.id
				row++
			}
		case "D":
			for i := 0; i < instruction.Steps; i++ {
				g[row][col] |= w.id
				row--
			}
		case "L":
			for i := 0; i < instruction.Steps; i++ {
				g[row][col] |= w.id
				col--
			}
		case "R":
			for i := 0; i < instruction.Steps; i++ {
				g[row][col] |= w.id
				col++
			}
		}
	}
}

func parsePath(line string) []instruction {
	var instructions []instruction
	for _, v := range strings.Split(line, ",") {
		direction := v[0]
		steps, _ := strconv.Atoi(string(v[1:]))
		instructions = append(instructions, instruction{
			Direction: string(direction),
			Steps:     steps,
		})
	}
	return instructions
}

func makeGrid(cols, rows int) [][]uint8 {
	grid := make([][]uint8, rows+1)
	for i := range grid {
		grid[i] = make([]uint8, cols)
	}
	return grid
}

func getCentralPort(w1, w2 wire) (startCol, startRow int) {
	_, w1Down, w1Left, _ := w1.bounds()
	_, w2Down, w2Left, _ := w2.bounds()

	startCol = max(abs(w1Left), abs(w2Left))
	startRow = max(abs(w1Down), abs(w2Down))

	return
}

func findIntersections(g [][]uint8, startCol, startRow int, wires ...wire) [][]int {
	var intersections [][]int

	intersectionFlag := intersectionOf(wires...)

	for row := range g {
		for col := range g[row] {
			if g[row][col] == intersectionFlag && (col != startCol && row != startRow) {
				intersections = append(intersections, []int{col, row})
			}
		}
	}
	return intersections
}

func sortIntersections(refCol, refRow int, intersections [][]int) {
	sort.Slice(intersections, func(i, j int) bool {
		a := intersections[i]
		b := intersections[j]

		distA := dist(refCol, refRow, a[0], a[1])
		distB := dist(refCol, refRow, b[0], b[1])

		return distA < distB
	})
}

func dist(refCol, refRow, checkCol, checkRow int) int {
	return abs(refCol-checkCol) + abs(refRow-checkRow)
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func maxBounds(w1, w2 wire) (up, down, left, right int) {
	w1Up, w1Down, w1Left, w1Right := w1.bounds()
	w2Up, w2Down, w2Left, w2Right := w2.bounds()

	return max(w1Up, w2Up), min(w1Down, w2Down), min(w1Left, w2Left), max(w1Right, w2Right)
}

func abs(v int) int {
	if v < 0 {
		return v * -1
	}
	return v
}

func intersectionOf(wires ...wire) uint8 {
	i := uint8(0)
	for _, w := range wires {
		i |= w.id
	}
	return i
}

func isIntersection(v uint8, wires ...wire) bool {
	i := intersectionOf(wires...)
	return i == v
}
