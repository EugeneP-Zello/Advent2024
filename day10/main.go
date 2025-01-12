package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var directions = [][2]int{
	{-1, 0},
	{1, 0},
	{0, -1},
	{0, 1},
}
	
func nextStep(x, y, elevation int, world [][]int, finishes map[[2]int]bool) int {
	if x < 0 || y < 0 || x >= len(world[0]) || y >= len(world) {
		return 0;
	}
	if world[x][y] != elevation {
		return 0;
	}
	if elevation == 9 {
		finishes[[2]int{x, y}] = true
		return 1
	}
	trails := 0
	for _, direction := range directions {
		nextX, nextY := x + direction[0], y + direction[1]
		trails += nextStep(nextX, nextY, elevation + 1, world, finishes)
	}
	return trails
}

func calc(filename string) (int, int) {
	file, _ := os.Open(filename)
	content, _ := io.ReadAll(file)
	file.Close()
	rows := strings.Split(string(content), "\n")
	worldMap := make([][]int, len(rows))
	for i, row := range rows {
		row = strings.TrimSuffix(row, "\r")
		worldMap[i] = make([]int, len(row))
		for j, char := range row {
			worldMap[i][j] = int(char - '0')
		}
	}

	count := 0
	rating := 0 
	// now iterate through the world, set obstacles and count loops
	for i, row := range worldMap {
		for j, elevation := range row {
			if elevation == 0 {
				finishes := make(map[[2]int]bool)
				rating += nextStep(i, j, 0, worldMap, finishes)
				for range finishes {
					count++
				}
			}
		}
	}
	return count, rating
}

func runForFile(filename string) {
	ssc, rating := calc(filename)
	fmt.Printf("File: %s, SumScore: %d, Rating: %d\n", filename, ssc, rating)
}

func main() {
	runForFile("test.txt")
	runForFile("input.txt")
}
