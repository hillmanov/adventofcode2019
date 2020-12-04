package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type object struct {
	Name     string
	Depth    int
	orbiters []*object
	orbitee  *object
}

func (o *object) countOrbits(indirect int) int {
	orbiterCount := indirect
	for _, orbiter := range o.orbiters {
		orbiterCount += orbiter.countOrbits(indirect + 1)
	}
	return orbiterCount
}

func (o *object) hasOrbiter(name string) bool {
	for _, orbiter := range o.orbiters {
		if orbiter.Name == name || orbiter.hasOrbiter(name) {
			return true
		}
	}
	return false
}

func (o *object) distanceToOrbitee(name string) int {
	steps := 0
	for o.Name != name && o != nil {
		steps++
		o = o.orbitee
	}
	return steps
}

var nodeMap map[string]*object

func main() {
	nodeMap = make(map[string]*object)
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ")")

		orbitee := getOrCreateNode(parts[0])
		orbiter := getOrCreateNode(parts[1])
		orbitee.orbiters = append(orbitee.orbiters, orbiter)
		orbiter.orbitee = orbitee
	}

	// Part 1
	allOrbits := getOrCreateNode("COM").countOrbits(0)
	fmt.Printf("Part 1: All orbits = %+v\n", allOrbits)

	// Part 2
	you := getOrCreateNode("YOU")
	san := getOrCreateNode("SAN")

	firstCommonParent := you.orbitee
	for !firstCommonParent.hasOrbiter("SAN") {
		firstCommonParent = firstCommonParent.orbitee
	}

	fmt.Printf("Part 2: Orbital Transfers: %d\n", you.orbitee.distanceToOrbitee(firstCommonParent.Name)+san.orbitee.distanceToOrbitee(firstCommonParent.Name))
}

func getOrCreateNode(name string) *object {
	if o, ok := nodeMap[name]; ok {
		return o
	}
	o := &object{
		Name: name,
	}
	nodeMap[name] = o
	return o
}
