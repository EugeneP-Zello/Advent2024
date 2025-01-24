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

func move(c *Coords) *Coords {
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

func left(c *Coords) *Coords {
	newDir := c.dir - 1
	if newDir < 0 {
		newDir = 3
	}
	return &Coords{
		x: c.x,
		y: c.y,
		dir: newDir,
		Score: c.Score + 1000,
		Path: c.Path,
	}
}

func right(c *Coords) *Coords {
	newDir := c.dir + 1
	if newDir > 3 {
		newDir = 0
	}
	return &Coords{
		x: c.x,
		y: c.y,
		dir: newDir,
		Score: c.Score + 1000,
		Path: c.Path,
	}
}

func PosIndex(c *Coords) int {
	return c.x + 1000*c.y + 1000000*c.dir
}

func Next(current *Coords, world [][]rune) []*Coords {
	r := make([]*Coords, 0)
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

	for _, row := range worldMap {
		fmt.Println(string(row))
	}
	fmt.Println("\r\n")

	pq := make(PriorityQueue, 1)
	pq[0] = &start
	heap.Init(&pq)
	bestPositions := make(map[int]bool)
	visited := make(map[int]int)
	next := make([]Coords, 1)
	next[0] = start
	score := 0
	for pq.Len() > 0 {
		crd := pq.GetMinValue()
		if score > 0 && score < crd.Score {
			continue
		}
		idx := PosIndex(crd)
		if prevScore, exists := visited[idx]; exists && prevScore < crd.Score {
			continue
		}
		if crd.x == finish.x && crd.y == finish.y {
			if score > 0 && score < crd.Score {
				break
			}
			score = crd.Score
			bestPositions[crd.y *1000 + crd.x] = true
			for _, v := range crd.Path {
				idxScore := v.y *1000 + v.x
				bestPositions[idxScore] = true
				worldMap[v.y][v.x] = 'O'
			}
			continue
		}
		visited[idx] = crd.Score
		for _, newCrd := range Next(crd, worldMap) {
			pq.Push(newCrd)
		}
	}

	fmt.Println("\r\n")
	for _, row := range worldMap {
		fmt.Println(string(row))
	}
	
	count := len(bestPositions)
	return score, count 
}

func runForFile(filename string) {
	Score, BestTilesCount := calcScore(filename)
	fmt.Printf("File: %s, Score : %d, best tiles count: %d\n", filename, Score, BestTilesCount)
}


func main() {
	runForFile("test3.txt")
	runForFile("test.txt")
	runForFile("test2.txt")
	runForFile("input.txt")
}
