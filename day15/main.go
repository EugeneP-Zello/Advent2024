package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Pos struct {
	x, y int
}

type Content struct {
	newValue rune
	Pos
}

func CheckMove(current Pos, dir Pos, world [][]rune, boxesToMove *[]Pos) bool {
	if world[current.y][current.x] == '.' {
		return true
	} else if world[current.y][current.x] == '#' {
		return false
	} else if world[current.y][current.x] == 'O' {
		*boxesToMove = append(*boxesToMove, current)
		current.x += dir.x
		current.y += dir.y
		return CheckMove(current, dir, world, boxesToMove)
	}
	return false
}

func CheckMove2(current Pos, dir Pos, world [][]rune, boxesToMove *[]Pos) bool {
	if world[current.y][current.x] == '.' {
		return true
	} else if world[current.y][current.x] == '#' {
		return false
	} else if world[current.y][current.x] == '[' || world[current.y][current.x] == ']' {
		if dir.y == 0 {
			*boxesToMove = append(*boxesToMove, current)
			current.x += dir.x
			return CheckMove2(current, dir, world, boxesToMove)
		} else {
			*boxesToMove = append(*boxesToMove, current)
			buddyBox := Pos {x:current.x,y:current.y}
			if world[current.y][current.x] == '[' {
				buddyBox.x += 1
			} else {
				buddyBox.x -= 1
			}
			*boxesToMove = append(*boxesToMove, buddyBox)
			current.y += dir.y
			buddyBox.y +=dir.y
			return CheckMove2(current, dir, world, boxesToMove) && CheckMove2(buddyBox, dir, world, boxesToMove)			
		}
	}
	return false
}


func DirectionFromRune(move rune) Pos {
	dir := Pos{x:0,y:0,}
	if move == '^' {
		dir.y = -1
	} else if move == 'v' {
		dir.y = 1
	} else if move == '>' {
		dir.x = 1
	} else if move == '<' {
		dir.x = -1
	}
	return dir
}

func MakeMove(move rune, current Pos, world [][]rune) Pos {
	dir := DirectionFromRune(move)
	boxesToMove := make([]Pos, 0)
	next := Pos{
		x: current.x + dir.x,
		y: current.y + dir.y,
	}
	if CheckMove(next, dir, world, &boxesToMove) {
		for _, box := range boxesToMove {
			world[box.y + dir.y][box.x + dir.x] = 'O'
		}
		world[current.y][current.x] = '.'
		world[next.y][next.x] = '@'

		//fmt.Println("\r\n\r\n")
		//for _, row := range world {
		//	fmt.Println(string(row))
		//}
		return next
	}
	return current
}

func MakeMove2(move rune, current Pos, world [][]rune) Pos {
	dir := DirectionFromRune(move)
	boxesToMove := make([]Pos, 0)
	next := Pos{
		x: current.x + dir.x,
		y: current.y + dir.y,
	}
	if CheckMove2(next, dir, world, &boxesToMove) {
		//fmt.Printf("\r\nMove %c %v\r\n", move, dir)
		if len(boxesToMove) > 0 {
			boxesMap := make(map[int]bool)
			newValues := make([]Content, 0)
			//fmt.Printf("Moving boxes %+v\r\n", boxesToMove)
			for _, box := range boxesToMove {
				if boxesMap[box.y*1000+box.x] {
					continue
				}
				boxesMap[box.y*1000+box.x] = true
				c := Content {
					newValue: world[box.y][box.x],
				}
				c.x = box.x + dir.x
				c.y = box.y + dir.y
				newValues = append(newValues, c)
				world[box.y][box.x] = '.'
			}
			for _, box := range newValues {
				world[box.y][box.x] = box.newValue
			}
		}
		world[current.y][current.x] = '.'
		world[next.y][next.x] = '@'

		
		//for _, row := range world {
		//	fmt.Println(string(row))
		//}
		return next
	} else {
		//fmt.Printf("\r\nMove %c %v -- hit the wall\r\n", move, dir)
	}
	return current
}


func calcGpsCoords(filename string) (int, int) {
	file, _ := os.Open(filename)
	content, _ := io.ReadAll(file)
	file.Close()
	parts := strings.Split(string(content), "\r\n\r\n")
	rows := strings.Split(parts[0], "\r\n")
	current := Pos{x:0,y:0,}
	worldMap := make([][]rune, len(rows))
	worldMap2 := make([][]rune, len(rows))
	for i, row := range rows {
		worldMap[i] = make([]rune, len(row))
		worldMap2[i] = make([]rune, len(row)*2)
		for j, char := range row {
			worldMap[i][j] = char
			if char == '@' {
				current.x = j
				current.y = i
				worldMap2[i][j*2] = '@'
				worldMap2[i][j*2 + 1] = '.'
			} else if char == '#' {
				worldMap2[i][j*2] = '#'
				worldMap2[i][j*2 + 1] = '#'
			} else if char == '.' {
				worldMap2[i][j*2] = '.'
				worldMap2[i][j*2 + 1] = '.'
			} else if char == 'O' {
				worldMap2[i][j*2] = '['
				worldMap2[i][j*2 + 1] = ']'
			}
		}
	}
	current2 := Pos {
		x: current.x*2,
		y: current.y,
	}
	fmt.Println("\r\n\r\n")
	for _, row := range worldMap2 {
		fmt.Println(string(row))
	}
	fmt.Println("\r\n\r\n")

	moves := strings.ReplaceAll(parts[1], "\r\n", "")
	for _, move := range moves {
		current = MakeMove(move, current, worldMap)
		current2 = MakeMove2(move, current2, worldMap2)
	}

	fmt.Println("\r\n\r\n")
	for _, row := range worldMap2 {
		fmt.Println(string(row))
	}
	fmt.Println("\r\n\r\n")

	totalGps, totalGps2 := 0, 0
	for i, row := range worldMap {
		for j, object := range row {
			if object == 'O' {
				totalGps += i*100 + j
			}
		}
	}
	for i, row := range worldMap2 {
		for j, object := range row {
			if object == '[' {
				totalGps2 += i*100 + j
			}
		}
	}

	return totalGps, totalGps2
}

func runForFile(filename string) {
	gpsCoords, gpsCoordsLarge := calcGpsCoords(filename)
	fmt.Printf("File: %s, GPS coordinates: %d, large GPS coordinates: %d\n", filename, gpsCoords, gpsCoordsLarge)
}

func main() {
	//runForFile("test.txt")
	runForFile("test3.txt")
	runForFile("input.txt")
}
