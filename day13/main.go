package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func solve(ax, ay, bx, by, x, y, max int) int {
	d := ax*by - ay*bx
	if d == 0 {
		return 0
	}
	c := x*by - y*bx
	a := c / d
	b := (x - a * ax) / bx
	if a*ax + b*bx != x || a*ay + b*by != y {
		return 0
	}
	if a > 0 && b > 0 && a < max && b < max {
		return a * 3 + b
	}
	return 0
}

func calcTokens(filename string) (int, int) {
	file, _ := os.Open(filename)
	content, _ := io.ReadAll(file)
	file.Close()
	inputs := strings.Split(string(content), "\r\n\r\n")
	totalTokens, totalTokensHuge := 0, 0
	for _, input := range inputs {
		rows := strings.Split(input, "\r\n")
		ax,ay,bx,by,x,y := 0,0,0,0,0,0
		fmt.Sscanf(rows[0], "Button A: X%d, Y%d", &ax, &ay)
		fmt.Sscanf(rows[1], "Button B: X%d, Y%d", &bx, &by)
		fmt.Sscanf(rows[2], "Prize: X=%d, Y=%d", &x, &y)
		totalTokens += solve(ax, ay, bx, by, x, y, 100)
		x += 10000000000000
		y += 10000000000000
		totalTokensHuge += solve(ax, ay, bx, by, x, y, 10000000000000)
	}
	return totalTokens, totalTokensHuge
}

func runForFile(filename string) {
	tokens, largeTokens := calcTokens(filename)
	fmt.Printf("File: %s, SumScore: %d, Rating: %d\n", filename, tokens, largeTokens)
}

func main() {
	runForFile("test.txt")
	runForFile("input.txt")
}