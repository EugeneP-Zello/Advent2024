package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var directions = [4][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

func visit(x, y int, plant rune, world [][]rune, visited map[[2]int]bool) int {
	if x < 0 || y < 0 || x >= len(world[0]) || y >= len(world) {
		return 0;
	}
	perimeter := 0
	area := 0
	stack := [][2]int{{x, y}}
	visited[[2]int{x, y}] = true
	for len(stack) > 0 {
		pos := stack[0]
		stack = stack[1:]
		area++
		for _, direction := range directions {
			nextX, nextY := pos[0] + direction[0], pos[1] + direction[1]
			if nextX < 0 || nextY < 0 || nextX >= len(world[0]) || nextY >= len(world) {
				perimeter++
				continue
			}
			if world[nextX][nextY] != plant {
				perimeter++
				continue
			}
			if !visited[[2]int{nextX, nextY}] {
				visited[[2]int{nextX, nextY}] = true
				stack = append(stack, [2]int{nextX, nextY})
			}
		}
	}
	return area * perimeter
}

func visit2(x, y int, plant rune, world [][]rune, visited map[[2]int]bool) int {
	if x < 0 || y < 0 || x >= len(world[0]) || y >= len(world) {
		return 0;
	}
	perimeter := 0
	area := 0
	stack := [][2]int{{x, y}}
	visited[[2]int{x, y}] = true
	for len(stack) > 0 {
		pos := stack[0]
		stack = stack[1:]
		area++
		for _, direction := range directions {
			nextX, nextY := pos[0] + direction[0], pos[1] + direction[1]
			if nextX < 0 || nextY < 0 || nextX >= len(world[0]) || nextY >= len(world) {
				perimeter++
				continue
			}
			if world[nextX][nextY] != plant {
				perimeter++
				continue
			}
			if !visited[[2]int{nextX, nextY}] {
				visited[[2]int{nextX, nextY}] = true
				stack = append(stack, [2]int{nextX, nextY})
			}
		}
	}
	return area * perimeter
}

func calc(filename string) (int, int) {
	file, _ := os.Open(filename)
	content, _ := io.ReadAll(file)
	file.Close()
	rows := strings.Split(string(content), "\n")
	worldMap := make([][]rune, len(rows))
	for i, row := range rows {
		row = strings.TrimSuffix(row, "\r")
		worldMap[i] = make([]rune, len(row))
		for j, char := range row {
			worldMap[i][j] = char
		}
	}

	total := 0
	visited := make(map[[2]int]bool)
	for i, row := range worldMap {
		for j, plant := range row {
			if !visited[[2]int{i, j}] {
				total += visit(i, j, plant, worldMap, visited)
			}
		}
	}

	discount := 0
	visited = make(map[[2]int]bool)
	for i, row := range worldMap {
		for j, plant := range row {
			if !visited[[2]int{i, j}] {
				discount += visit2(i, j, plant, worldMap, visited)
			}
		}
	}

	return total, discount
}

func runForFile(filename string) {
	price, discountPrice := calc(filename)
	fmt.Printf("File: %s, SumScore: %d, Rating: %d\n", filename, price, discountPrice)
}

func main() {
	runForFile("test.txt")
	runForFile("input.txt")
}
