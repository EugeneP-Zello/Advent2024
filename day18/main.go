package main

import (
	"bufio"
	"container/heap"
	"fmt"
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

func Next(current *Coords, world [][]rune, maxIdx int, track bool) []*Coords {
	r := make([]*Coords, 0)
	p1 := move(current, track)
	if p1.x >= 0 && p1.x < maxIdx && p1.y >= 0 && p1.y < maxIdx {
		if world[p1.y][p1.x] != '#' {
			r = append(r, p1)
		}
	}
	
	p2 := left(current, track)
	if p2.x >= 0 && p2.x < maxIdx && p2.y >= 0 && p2.y < maxIdx {
		if world[p2.y][p2.x] != '#' {
			r = append(r, p2)
		}
	}
	
	p3 := right(current, track)
	if p3.x >= 0 && p3.x < maxIdx && p3.y >= 0 && p3.y < maxIdx {
		if world[p3.y][p3.x] != '#' {
			r = append(r, p3)
		}
	}

	return r
}

func pathExists(obstacles []Point, gridSize int, fallenBlock int) bool {
	worldMap := make([][]rune, gridSize)
	for y := 0; y < gridSize; y++ {
		worldMap[y] = make([]rune, gridSize)
		for x := 0; x < gridSize; x++ {
			worldMap[y][x] = '.'
		}
	}
	for idx := 0; idx < fallenBlock; idx++ {
		worldMap[obstacles[idx].y][obstacles[idx].x] = '#'
	}

	pq := make(PriorityQueue, 1)
	start := Coords{
		x: 0,
		y: 0,
		dir: 1,
		Score: 0,
		Path: make([]*Coords,0),
	}
	finish := Coords{
		x: gridSize - 1,
		y: gridSize - 1,
		dir: 0,
		Score: 0,
		Path: make([]*Coords,0),
	}
	pq[0] = &start
	heap.Init(&pq)
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
		if prevScore, exists := visited[idx]; exists && prevScore <= crd.Score {
			continue
		}
		if crd.x == finish.x && crd.y == finish.y {
			return true
		}
		visited[idx] = crd.Score
		for _, newCrd := range Next(crd, worldMap, gridSize, false) {			
			newIdx := PosIndex(newCrd)
			if prevScore, exists := visited[newIdx]; exists && prevScore < newCrd.Score {
				continue
			}
			pq.Push(newCrd)
		}
	}
	return false
}


func processGrid(filename string, gridSize int, fallenBlock int) (int, int) {
	worldMap := make([][]rune, gridSize)
	for y := 0; y < gridSize; y++ {
		worldMap[y] = make([]rune, gridSize)
		for x := 0; x < gridSize; x++ {
			worldMap[y][x] = '.'
		}
	}
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	obstacles := make([]Point, 0)
	p := Point{}
	for scanner.Scan() {
		s := scanner.Text()
		fmt.Sscanf(s, "%d,%d", &p.x, &p.y)
		obstacles = append(obstacles, p)
	}
	for idx := 0; idx < fallenBlock; idx++ {
		worldMap[obstacles[idx].y][obstacles[idx].x] = '#'
	}

	for _, row := range worldMap {
		fmt.Println(string(row))
	}
	fmt.Println("\r\n")

	pq := make(PriorityQueue, 1)
	start := Coords{
		x: 0,
		y: 0,
		dir: 1,
		Score: 0,
		Path: make([]*Coords,0),
	}
	finish := Coords{
		x: gridSize - 1,
		y: gridSize - 1,
		dir: 0,
		Score: 0,
		Path: make([]*Coords,0),
	}
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
			for _, v := range crd.Path {
				worldMap[v.y][v.x] = 'O'
			}
			continue
		}
		visited[idx] = crd.Score
		for _, newCrd := range Next(crd, worldMap, gridSize, true) {			
			newIdx := PosIndex(newCrd)
			if prevScore, exists := visited[newIdx]; exists && prevScore < newCrd.Score {
				continue
			}
			pq.Push(newCrd)
		}
	}

	fmt.Println("\r\n")
	for _, row := range worldMap {
		fmt.Println(string(row))
	}
	
	okBlocks := fallenBlock
	failedBlocks := len(obstacles)
	for okBlocks+1 != failedBlocks {
		blockCount := okBlocks + (failedBlocks - okBlocks) / 2
		if blockCount == okBlocks {
			blockCount+=1
		}
		if pathExists(obstacles, gridSize, blockCount) {
			okBlocks = blockCount
			fmt.Printf("Check %d - ok\r\n", blockCount)
		} else {
			failedBlocks = blockCount
			fmt.Printf("Check %d (%d,%d)- failure\r\n", blockCount, obstacles[blockCount].x, obstacles[blockCount].y)
		}
	}
	return score, failedBlocks//fmt.Sprintf("%d,%d", obstacles[failedBlocks].x, obstacles[failedBlocks].y) 
}

func runForFile(filename string, gridSize int, fallenBlock int) {
	Steps, BlocksCount := processGrid(filename, gridSize, fallenBlock)
	fmt.Printf("File: %s, Steps : %d, critical tiles count: %d\n", filename, Steps, BlocksCount)
}


func main() {
	//runForFile("test.txt", 7, 12)
	//runForFile("input.txt", 71, 1024)
	runForFile("input.txt", 71, 2956)
}
