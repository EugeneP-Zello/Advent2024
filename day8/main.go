package main

import (
	"bufio"
	"fmt"
	"os"
)
type Pos struct {
	x, y int
}

func calc(filename string) (int, int) {
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	total1, total2 := 0, 0
	var world [][]rune
	for scanner.Scan() {
		s := scanner.Text()
		world = append(world, []rune(s))
	}
	antennas := make(map[rune][]Pos)
	for y, row := range world {
		for x, cell := range row {
			if cell != '.' {
				antennas[cell] = append(antennas[cell], Pos{x, y})
			}
		}
	}
	world1 := make([][]rune, len(world))
	world2 := make([][]rune, len(world))

	for idx, row := range world {
		world1[idx] = make([]rune, len(row))
		copy(world1[idx], row)
		world2[idx] = make([]rune, len(row))
		copy(world2[idx], row)
	}

	for _, antenna := range antennas {
		for i := 0; i < len(antenna); i++ {
			for j := i + 1; j < len(antenna); j++ {
				deltaX := antenna[i].x - antenna[j].x
				deltaY := antenna[i].y - antenna[j].y
				nodeX := antenna[i].x + deltaX
				nodeY := antenna[i].y + deltaY
				if nodeX >= 0 && nodeX < len(world[0]) && nodeY >= 0 && nodeY < len(world) {
					world1[nodeY][nodeX] = '#'
				}
				nodeX = antenna[j].x - deltaX
				nodeY = antenna[j].y - deltaY
				if nodeX >= 0 && nodeX < len(world[0]) && nodeY >= 0 && nodeY < len(world) {
					world1[nodeY][nodeX] = '#'
				}
			}
		}
	}

	for _, antenna := range antennas {
		for i := 0; i < len(antenna); i++ {
			for j := i + 1; j < len(antenna); j++ {
				deltaX := antenna[i].x - antenna[j].x
				deltaY := antenna[i].y - antenna[j].y
				nodeX := antenna[i].x
				nodeY := antenna[i].y
				for nodeX >= 0 && nodeX < len(world[0]) && nodeY >= 0 && nodeY < len(world) {
					world2[nodeY][nodeX] = '#'
					nodeX += deltaX
					nodeY += deltaY
				}
				nodeX = antenna[j].x
				nodeY = antenna[j].y
				for nodeX >= 0 && nodeX < len(world[0]) && nodeY >= 0 && nodeY < len(world) {
					world2[nodeY][nodeX] = '#'
					nodeX -= deltaX
					nodeY -= deltaY
				}
			}
		}
	}

	for _, row := range world1 {
		for _, cell := range row {
			if cell == '#' {
				total1++
			}
		}
	}

	for _, row := range world2 {
		for _, cell := range row {
			if cell == '#' {
				total2++
			}
		}
	}

	return total1, total2
}

func runForFile(filename string) {
	c, lc := calc(filename)
	fmt.Printf("File: %s, simple antinodes: %d, resonanse antinodes: %d\n", filename, c, lc)
}

func main() {
	runForFile("test2.txt")
	runForFile("test.txt")
	runForFile("input.txt")
}
