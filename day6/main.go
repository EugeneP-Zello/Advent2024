package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

var change [4][2]int = [4][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

type state struct {
	x, y       int
	maxX, maxY int
	world      [][]rune
	direction  int
	visited    mapset.Set[int]
}

func (s *state) calcVisited() int {
	return s.maxX*s.maxY*s.direction + s.maxX*s.y + s.x
}

func newState(x, y, direction int, world [][]rune) state {
	s := state{
		x:         x,
		y:         y,
		maxX:      len(world[0]),
		maxY:      len(world),
		direction: direction,
		world:     make([][]rune, len(world)),
		visited:   mapset.NewSet[int](),
	}
	for idx, row := range world {
		s.world[idx] = make([]rune, len(row))
		copy(s.world[idx], row)
	}
	return s
}

func detectLoop(s state) bool {
	for {
		currentVisited := s.calcVisited()
		if s.visited.Contains(currentVisited) {
			return true
		}
		s.visited.Add(currentVisited)
		nextX, nextY := s.x+change[s.direction][0], s.y+change[s.direction][1]
		if nextY < 0 || nextY >= s.maxY || nextX < 0 || nextX >= s.maxX {
			return false
		}
		if s.world[nextY][nextX] == '#' {
			s.direction = (s.direction + 1) % 4
		} else {
			s.x, s.y = nextX, nextY
		}
	}
}

func calc(filename string) (int, int) {
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var world [][]rune
	startX, startY := 0, 0

	changeDirection := 0
	row := 0
	for scanner.Scan() {
		s := scanner.Text()
		if startIndex := strings.LastIndexByte(s, '^'); startIndex >= 0 {
			startY = row
			startX = startIndex
		}
		runeSlice := []rune(s)
		world = append(world, runeSlice)
		row++
	}
	file.Close()
	fmt.Printf("World size: %d x %d\n", len(world[0]), len(world))
	fmt.Printf("Start position: (%d, %d) %c\n", startY, startX, world[startY][startX])
	x, y := startX, startY
	maxWidth, maxHeight := len(world[0]), len(world)
	for {
		world[y][x] = 'X'
		nextX, nextY := x+change[changeDirection][0], y+change[changeDirection][1]
		if nextY < 0 || nextY >= maxHeight || nextX < 0 || nextX >= maxWidth {
			break
		}
		if world[nextY][nextX] == '#' {
			changeDirection = (changeDirection + 1) % 4
		} else {
			x, y = nextX, nextY
		}
	}

	count := 0
	for _, row := range world {
		for _, r := range row {
			if r == 'X' {
				count++
			}
		}
	}
	loopCount := 0
	// now iterate through the world, set obstacles and count loops
	for iy, row := range world {
		for ix, r := range row {
			if r == 'X' && (ix != startX || iy != startY) {
				s := newState(startX, startY, 0, world)
				s.world[iy][ix] = '#'
				if detectLoop(s) {
					loopCount++
				}
			}
		}
	}
	return count, loopCount
}

func runForFile(filename string) {
	c, lc := calc(filename)
	fmt.Printf("File: %s, count: %d, loop count: %d\n", filename, c, lc)
}

func main() {
	runForFile("test.txt")
	runForFile("input.txt")
}
