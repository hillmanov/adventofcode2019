package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
)

const rows = 6
const cols = 25

func main() {
	input, _ := ioutil.ReadFile("./input.txt")

	encImage := getEncodedSpaceImage(string(input), cols, rows)

	// Part 1
	func() {
		var layerWithFewest0s int
		fewest0s := math.MaxInt64

		for i, layer := range encImage {
			amountOf0s := 0
			for _, row := range layer {
				amountOf0s += count(row, 0)
			}
			if amountOf0s < fewest0s {
				layerWithFewest0s = i
				fewest0s = amountOf0s
			}
		}

		var amountOf1s int
		var amountOf2s int
		for _, row := range encImage[layerWithFewest0s] {
			amountOf1s += count(row, 1)
			amountOf2s += count(row, 2)
		}

		result := amountOf1s * amountOf2s
		fmt.Printf("Part 1: %d\n", result)
	}()

	// Part 2

	func() {
		var image [rows][cols]int

		for row := 0; row < rows; row++ {
			for col := 0; col < cols; col++ {
				for layer := 0; layer < len(encImage); layer++ {
					value := encImage[layer][row][col]
					if layer == 0 {
						image[row][col] = value
					} else {
						switch image[row][col] {
						case 2:
							image[row][col] = value
						}
					}
				}
			}
		}

		fmt.Println("Part 2:")
		for _, row := range image {
			for _, col := range row {
				if col == 1 {
					fmt.Printf("0")
				} else {
					fmt.Printf(" ")
				}
			}
			fmt.Printf("\n")
		}
	}()

}

func getEncodedSpaceImage(s string, width, height int) [] /*layer*/ [] /*row*/ [] /*col*/ int {
	layers := len(s) / (width * height)
	image := make([][][]int, layers)

	for layer := 0; layer < layers; layer++ {
		image[layer] = make([][]int, height)
		for row := 0; row < height; row++ {
			image[layer][row] = make([]int, width)
			for col := 0; col < width; col++ {
				index := col + row*width + layer*width*height
				num, _ := strconv.Atoi(string(s[index]))
				image[layer][row][col] = num
			}
		}
	}

	return image
}

func count(haystack []int, needle int) (count int) {
	for _, v := range haystack {
		if v == needle {
			count++
		}
	}
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
