package main

import (
	"bufio"
	"fmt"
	"os"
)

type Robot struct {
	x, y, vx, vy int
}

func (r *Robot) CalcPos(t, width, height int) (int, int) {
	x := (r.x + t*r.vx) % width
	y := (r.y + t*r.vy) % height
	if x < 0 {
		x += width
	}
	if y < 0 {
		y += height
	}
	return x, y
} 

func calcRobotsFactor(robots []Robot, t, width, height int) int {
	ne, nw, se, sw := 0,0,0,0
	midX, midY := (width - 1) / 2, (height -1) / 2
	for _, r := range robots {
		x, y := r.CalcPos(t, width, height)
		if x < midX {
			if y < midY {
				ne++
			} else if y > midY {
				nw++
			}
		} else if x > midX {
			if y < midY {
				se++
			} else if y > midY {
				sw++
			}
		}
	}
	return se * sw * ne *nw
}

func processFile(filename string, width int, height int) (int, int) {
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	robots := make([]Robot, 0)
	for scanner.Scan() {
		s := scanner.Text()
		// p=0,4 v=3,-3
		r := Robot{}
		fmt.Sscanf(s, "p=%d,%d v=%d,%d", &r.x, &r.y, &r.vx, &r.vy)
		robots = append(robots, r)
	}
	sf := calcRobotsFactor(robots, 100, width, height)
	min_sf := sf
	min_t := 100
	for idx := 0; idx < width*height; idx++ {
		new_sf := calcRobotsFactor(robots, idx, width, height)
		if new_sf < min_sf {
			min_t = idx
			min_sf = new_sf
		}
	}
	return sf, min_t
}

func runForFile(filename string, width int, height int) {
	sf, min_sf := processFile(filename, width, height)
	fmt.Printf("File: %s, SumScore: %d, Rating: %d\n", filename, sf, min_sf)
}

func main() {
	runForFile("test.txt", 11, 7)
	runForFile("input.txt", 101, 103)
}