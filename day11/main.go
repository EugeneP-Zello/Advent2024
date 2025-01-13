package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func calcMap(initStones []int, steps int) int {
	stones := make(map[int]int)
	for _, value := range initStones {
		stones[value] = 1
	}
	for i := 0; i < steps; i++ {
		newMap := make(map[int]int)
		for stone, count := range stones {
			newStones, valid := split(stone)
			if valid {
				newMap[newStones[0]] += count
				newMap[newStones[1]] += count
			} else if stone == 0 {
				newMap[1] += count
			} else {
				newMap[stone * 2024] += count
			}
		}
		stones = newMap
	}
	total := 0
	for _, count := range stones {
		total += count
	}
	return total
}

func split(num int) (stones []int, valid bool) {
	converted := strconv.Itoa(num)
	valid = len(converted) % 2 == 0
	if !valid {
		return
	}
	middle := len(converted) / 2
	stones = make([]int, 2)
	stones[0], _ = strconv.Atoi(converted[:middle])
	stones[1], _ = strconv.Atoi(converted[middle:])
	return
}

func blink(stones []int) []int {
	nextStones := make([]int, 0)
	for _, stone := range stones {
		newStones, valid := split(stone)
		if valid {
			nextStones = append(nextStones, newStones...)
		} else if stone == 0 {
			nextStones = append(nextStones, 1)
		} else {
			nextStones = append(nextStones, stone * 2024)
		}
	}
	return nextStones
}

func calc(filename string) (int, int) {
	file, _ := os.Open(filename)
	content, _ := io.ReadAll(file)
	initStones := strings.Split(string(content), " ")
	stones := make([]int, len(initStones))
	for idx, value := range initStones {
		stones[idx], _ = strconv.Atoi(value)
	}
	cMap := calcMap(stones, 75)
	for i := 0; i < 25; i++ {
		stones = blink(stones)
	}
	return len(stones), cMap
}

func runForFile(filename string) {
	count, cv2 := calc(filename)
	fmt.Printf("File: %s, Stones count: %d, v2 : %d\n", filename, count, cv2)
}

func main() {
	runForFile("test.txt")
	runForFile("input.txt")
}