package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type Position struct {
	X int
	Y int
	Z int
}

type Velocity struct {
	X int
	Y int
	Z int
}

type GravityDelta struct {
	X int
	Y int
	Z int
}

type Moon struct {
	Position     *Position
	Velocity     *Velocity
	GravityDelta *GravityDelta
}

func (m Moon) String() string {
	return fmt.Sprintf("\nPosition: %+v || Velocity: %+v || GravityDelta: %+v\n", m.Position, m.Velocity, m.GravityDelta)
}

func (m Moon) PotentialEnergy() int {
	return abs(m.Position.X) + abs(m.Position.Y) + abs(m.Position.Z)
}

func (m Moon) KineticEnergy() int {
	return abs(m.Velocity.X) + abs(m.Velocity.Y) + abs(m.Velocity.Z)
}

func (m Moon) TotalEnergy() int {
	return m.KineticEnergy() * m.PotentialEnergy()
}

func (m Moon) XDimension() string {
	return fmt.Sprintf("%d:%d", m.Position.X, m.Velocity.X)
}

func (m Moon) YDimension() string {
	return fmt.Sprintf("%d:%d", m.Position.Y, m.Velocity.Y)
}

func (m Moon) ZDimension() string {
	return fmt.Sprintf("%d:%d", m.Position.Z, m.Velocity.Z)
}

type Moons []Moon

func (moons Moons) Step() {
	for i := 0; i < len(moons); i++ {
		for j := i + 1; j < len(moons); j++ {
			CalculateGravityDelta(&moons[i], &moons[j])
		}
	}
	for _, moon := range moons {
		moon.ApplyGravityDelta()
		moon.Move()
	}
}

func (moons Moons) TotalEnergy() (totalEnergy int) {
	for _, moon := range moons {
		totalEnergy += moon.TotalEnergy()
	}
	return
}

func (moons Moons) XDimension() string {
	dimX := make([]string, len(moons))
	for _, moon := range moons {
		dimX = append(dimX, moon.XDimension())
	}
	return strings.Join(dimX, ",")
}

func (moons Moons) YDimension() string {
	dimY := make([]string, len(moons))
	for _, moon := range moons {
		dimY = append(dimY, moon.YDimension())
	}
	return strings.Join(dimY, ",")
}

func (moons Moons) ZDimension() string {
	dimZ := make([]string, len(moons))
	for _, moon := range moons {
		dimZ = append(dimZ, moon.ZDimension())
	}
	return strings.Join(dimZ, ",")
}

func CalculateGravityDelta(m, o *Moon) {
	if o.Position.X > m.Position.X {
		m.GravityDelta.X++
		o.GravityDelta.X--
	} else if o.Position.X < m.Position.X {
		m.GravityDelta.X--
		o.GravityDelta.X++
	}

	if o.Position.Y > m.Position.Y {
		m.GravityDelta.Y++
		o.GravityDelta.Y--
	} else if o.Position.Y < m.Position.Y {
		m.GravityDelta.Y--
		o.GravityDelta.Y++
	}

	if o.Position.Z > m.Position.Z {
		m.GravityDelta.Z++
		o.GravityDelta.Z--
	} else if o.Position.Z < m.Position.Z {
		m.GravityDelta.Z--
		o.GravityDelta.Z++
	}
}

func (m *Moon) ApplyGravityDelta() {
	m.Velocity.X += m.GravityDelta.X
	m.Velocity.Y += m.GravityDelta.Y
	m.Velocity.Z += m.GravityDelta.Z
	m.ResetGravityDelta()
}

func (m *Moon) ResetGravityDelta() {
	m.GravityDelta.X = 0
	m.GravityDelta.Y = 0
	m.GravityDelta.Z = 0
}

func (m *Moon) Move() {
	m.Position.X += m.Velocity.X
	m.Position.Y += m.Velocity.Y
	m.Position.Z += m.Velocity.Z
}

func (m *Moon) UnmarshalJSON(data []byte) error {
	m.Position = &Position{}
	m.Velocity = &Velocity{}
	m.GravityDelta = &GravityDelta{}

	if err := json.Unmarshal(data, &m.Position); err != nil {
		return err
	}
	return nil
}

func main() {
	input, _ := ioutil.ReadFile("./input.txt")
	jsonInput := convertInputToJSON(input)

	var moons Moons
	json.Unmarshal(jsonInput, &moons)

	for i := 0; i < 1000; i++ {
		moons.Step()
	}

	fmt.Printf("Part 1: %d\n", moons.TotalEnergy())

	// Start over
	json.Unmarshal(jsonInput, &moons)

	xInitial := moons.XDimension()
	yInitial := moons.YDimension()
	zInitial := moons.ZDimension()

	xPeriod := 0
	yPeriod := 0
	zPeriod := 0

	steps := 0
	for xPeriod == 0 || yPeriod == 0 || zPeriod == 0 {
		moons.Step()
		steps++
		if xPeriod == 0 {
			if moons.XDimension() == xInitial {
				xPeriod = steps
			}
		}
		if yPeriod == 0 {
			if moons.YDimension() == yInitial {
				yPeriod = steps
			}
		}
		if zPeriod == 0 {
			if moons.ZDimension() == zInitial {
				zPeriod = steps
			}
		}
	}

	fmt.Printf("Part 2: %+v\n", LCM(xPeriod, yPeriod, zPeriod))
}

// Quick and dirty. Nothing more.
func convertInputToJSON(data []byte) []byte {
	json := "[" + string(data) + "]"
	json = strings.ReplaceAll(json, "=", ":")
	json = strings.ReplaceAll(json, "<", "{")
	json = strings.ReplaceAll(json, ">", "}")
	json = strings.ReplaceAll(json, "x", "\"x\"")
	json = strings.ReplaceAll(json, "y", "\"y\"")
	json = strings.ReplaceAll(json, "z", "\"z\"")
	json = strings.ReplaceAll(json, "\n", ",")
	json = strings.ReplaceAll(json, ",]", "]")
	return []byte(json)
}

func abs(i int) int {
	if i < 0 {
		return i * -1
	}
	return i
}

// Greatest common divisor
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}
