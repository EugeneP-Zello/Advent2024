package main

import (
	"container/heap"
	"fmt"
	"io"
	"strings"
	"os"
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
	Path []*Coords
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

func move(c *Coords, track bool) *Coords {
	if !track {
		return &Coords{
			x: c.x + directions[c.dir].x,
			y: c.y + directions[c.dir].y,
			dir: c.dir,
			Score: c.Score + 1,
		}
	}
	Path := make([]*Coords,len(c.Path))
	copy(Path, c.Path)
	return &Coords{
		x: c.x + directions[c.dir].x,
		y: c.y + directions[c.dir].y,
		dir: c.dir,
		Score: c.Score + 1,
		Path: append(Path, c),
	}
}

func left(c *Coords, track bool) *Coords {
	newDir := c.dir - 1
	if newDir < 0 {
		newDir = 3
	}
	c2 := Coords{
		x: c.x,
		y: c.y,
		dir: newDir,
		Score: c.Score,
		Path: c.Path,
	}
	return move(&c2, track)
}

func right(c *Coords, track bool) *Coords {
	newDir := c.dir + 1
	if newDir > 3 {
		newDir = 0
	}
	c2 := Coords{
		x: c.x,
		y: c.y,
		dir: newDir,
		Score: c.Score,
		Path: c.Path,
	}
	return move(&c2, track)
}

func PosIndex(c *Coords) int {
	return c.x + 1000*c.y
}

func Next(current *Coords, world [][]rune, track bool) []*Coords {
	r := make([]*Coords, 0)
	p1 := move(current, track)
	if world[p1.y][p1.x] != '#' {
		r = append(r, p1)
	}
	
	p2 := left(current, track)
	if world[p2.y][p2.x] != '#' {
		r = append(r, p2)
	}
	
	p3 := right(current, track)
	if world[p3.y][p3.x] != '#' {
		r = append(r, p3)
	}

	return r
}

func calcDistance(node1 *Coords, node2 *Coords) int {
	distance := 0
	if node1.x > node2.x {
		distance += node1.x - node2.x
	} else {
		distance += node2.x - node1.x
	}
	if node1.y > node2.y {
		distance += node1.y - node2.y
	} else {
		distance += node2.y - node1.y
	}
	return distance
}

func findAllCheats(filename string, minWin int, cheatLengthMax int) int { 
	file, _ := os.Open(filename)
	content, _ := io.ReadAll(file)
	file.Close()
	rows := strings.Split(string(content), "\r\n")
	start := Coords {
	  dir: 0,
	  Path: make([]*Coords, 0),	
	}
	finish := Coords{Path: make([]*Coords, 0),}
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


	//for _, row := range worldMap {
	//	fmt.Println(string(row))
	//}
	//fmt.Println("\r\n")

	pq := make(PriorityQueue, 1)
	pq[0] = &start
	heap.Init(&pq)
	visited := make(map[int]int)
	next := make([]Coords, 1)
	next[0] = start
	score := 0
	cnt := 0
	for pq.Len() > 0 {
		crd := pq.GetMinValue()
		if cnt++; cnt == 300 {
			cnt = 0
			fmt.Printf("V=%d, pq=%d crd=(%d,%d)\r\n", len(visited), pq.Len(), crd.x, crd.y)
		}
		
		if score > 0 && score < crd.Score {
			continue
		}
		idx := PosIndex(crd)
		if prevScore, exists := visited[idx]; exists && prevScore <= crd.Score {
			continue
		}
		if crd.x == finish.x && crd.y == finish.y {
			score = crd.Score
			finish.Score = crd.Score
			finish.Path = append(crd.Path, crd)
			for _, v := range crd.Path {
				worldMap[v.y][v.x] = 'O'
			}
			continue
		}
		visited[idx] = crd.Score
		for _, newCrd := range Next(crd, worldMap, true) {			
			newIdx := PosIndex(newCrd)
			if prevScore, exists := visited[newIdx]; exists && prevScore < newCrd.Score {
				continue
			}
			pq.Push(newCrd)
		}
	}

	//fmt.Println("\r\n")
	//for _, row := range worldMap {
	//	fmt.Println(string(row))
	//}
	count := 0
	max := len(finish.Path)
	for idx, node1 := range finish.Path[:max-2] {
		for _, node2 := range finish.Path[idx+2:] {
			if dist := calcDistance(node1, node2); dist <= cheatLengthMax {
				if node2.Score - node1.Score >= minWin + dist {
					count += 1
				}
			}
		}
	}

	return count
}

func runForFile(filename string, minWin int, cheatLengthMax int) {
	cheatCount := findAllCheats(filename, minWin, cheatLengthMax)
	fmt.Printf("File: %s, cheat(%d) count %d or above: %d\r\n", filename, cheatLengthMax, minWin, cheatCount)
}


func main() {
	runForFile("test.txt", 20, 2)
	runForFile("test.txt", 12, 2)
	runForFile("test.txt", 10, 2)
	runForFile("test.txt", 74, 20)
	runForFile("test.txt", 70, 20)
	runForFile("input.txt", 100, 2)
	runForFile("input.txt", 100, 20)
}
