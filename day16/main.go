package main

import (
	"container/heap"
	"fmt"
	"io"
	"os"
	"strings"
)


type Point struct {
	x, y int
}

var directions = []Point{
	{0, -1},
	{1, 0},
	{0, 1},
	{-1, 0},
}

type Coords struct {
	x, y, dir int
	Score int
}

type PriorityQueue []*Coords

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Score < pq[j].Score }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*pq = append(*pq, x.(*Coords))
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

func (pq *PriorityQueue) GetMinValue() *Coords {
	return heap.Pop(pq).(*Coords)
}

func move(c Coords) Coords {
	return Coords{
		x: c.x + directions[c.dir].x,
		y: c.y + directions[c.dir].y,
		dir: c.dir,
		Score: c.Score +1,
	}
}

func left(c Coords) Coords {
	newDir := c.dir - 1
	if newDir < 0 {
		newDir = 3
	}
	return Coords{
		x: c.x,
		y: c.y,
		dir: newDir,
		Score: c.Score + 1000,
	}
}

func right(c Coords) Coords {
	newDir := c.dir + 1
	if newDir > 3 {
		newDir = 0
	}
	return Coords{
		x: c.x,
		y: c.y,
		dir: newDir,
		Score: c.Score + 1000,
	}
}

func PosIndex(c Coords) int {
	return c.x + 1000*c.y + 1000000*c.dir
}

type Position struct {
	Score	int
	Coords
}

func Next(current Coords, world [][]rune, visited map[int]bool) []Coords {
	r := make([]Coords, 0)
	p1 := move(current)
	if world[p1.y][p1.x] != '#' {
		r = append(r, p1)
	}
	
	p2 := left(current)
	if tmp := move(p2); world[tmp.y][tmp.x] != '#' {
		r = append(r, p2)
	}
	
	p3 := right(current)
	if tmp := move(p3); world[tmp.y][tmp.x] != '#' {
		r = append(r, p3)
	}
	return r
}

func calcScore(filename string) (int, int) {
	file, _ := os.Open(filename)
	content, _ := io.ReadAll(file)
	file.Close()
	rows := strings.Split(string(content), "\r\n")
	start := Coords {
	  dir: 1,	
	}
	finish := Coords{}
	worldMap := make([][]rune, len(rows))
	for y, row := range rows {
		worldMap[y] = make([]rune, len(row))
		for x, char := range row {
			worldMap[y][x] = char
			if char == 'S' {
				start.x = x
				start.y = y
			} else if char == 'E' {
				finish.x = x
				finish.y = y
			}
		}
	}
	fmt.Printf("Start:{%d,%d}, end{%d,%d}, %c\r\n", start.x, start.y, finish.x,finish.y, worldMap[start.y][start.x])

	for _, row := range worldMap {
		fmt.Println(string(row))
	}
	fmt.Println("\r\n")

	pq := make(PriorityQueue, 1)
	pq[0] = &start
	heap.Init(&pq)
	visited := make(map[int]bool)
	scores := make(map[int]int)
	scores[start.y*1000+start.x] = 0
	next := make([]Coords, 1)
	next[0] = start
	for len(next) > 0 {
		second := make([]Coords, 0)
		for _, crd := range next {
			idx := PosIndex(crd)
			if !visited[idx] {
				// fmt.Printf("{%2d,%2d}\r\n", crd.y, crd.x)
				visited[idx] = true
				for _, newCrd := range Next(crd, worldMap, visited) {
					idxScore := newCrd.y *1000 + newCrd.x
					if value, exists := scores[idxScore]; !exists || value > newCrd.Score {
						scores[idxScore] = newCrd.Score
						// fmt.Printf("{%2d,%2d} - %d\r\n", crd.y, crd.x, newCrd.Score)
					}
					second = append(second, newCrd)
				}
			}
		}
		next = second
	}
	
	lowestScore := scores[finish.y*1000+finish.x]

	return lowestScore, 0
}

func runForFile(filename string) {
	Score, Score2 := calcScore(filename)
	fmt.Printf("File: %s, Score#1 : %d, score#2: %d\n", filename, Score, Score2)
}


func main() {
//	runForFile("test3.txt")
//	runForFile("test.txt")
//	runForFile("test2.txt")
//	runForFile("input.txt")
	file, _ := os.Open("input.txt")
	content, _ := io.ReadAll(file)
	file.Close()
	rows := strings.Split(string(content), "\r\n")
	rows[0] = ""
	fmt.Printf("\r\n")
	day16_1(string(content))
}
