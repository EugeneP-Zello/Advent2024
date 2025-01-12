package main

import (
	"fmt"
	"io"
	"os"
)

func sumDigits(s string) int {
	sum := 0
	for _, char := range s {
		sum += int(char - '0')
	}
	return sum
}

func calc(filename string) (int, int) {
	file, _ := os.Open(filename)
	content, _ := io.ReadAll(file)
	s := string(content)
	size := sumDigits(s)
	blocks := make([]int, size)
	pos := 0
	isEmptySpace := false
	for idx, char := range s {
		blockLen := int(char - '0')
		fileNum := idx / 2
		if isEmptySpace {
			fileNum = -1
		}
		for i := 0; i < blockLen; i++ {
			blocks[pos] = fileNum
			pos++
		}
		isEmptySpace = !isEmptySpace
	}
	
	blocks2 := make([]int, size)
	copy(blocks2, blocks)

	emptyPos := 0
	fileBlockPos := size - 1
	for blocks[emptyPos] != -1 {
		emptyPos++
	}
	for blocks[fileBlockPos] == -1 {
		fileBlockPos--
	}
	for emptyPos < fileBlockPos {
		blocks[emptyPos] = blocks[fileBlockPos]
		blocks[fileBlockPos] = -1
		for blocks[emptyPos] != -1 {
			emptyPos++
		}
		for blocks[fileBlockPos] == -1 {
			fileBlockPos--
		}
	}
// part 2 starts here
	fileBlockPos = size - 1
	for blocks2[fileBlockPos] != 0 && fileBlockPos >= 0 {
		for blocks2[fileBlockPos] == -1 && fileBlockPos > 0 {
			fileBlockPos--
		}
		fileEndPos := fileBlockPos
		fileId := blocks2[fileBlockPos]
		for blocks2[fileBlockPos] == fileId && fileBlockPos > 0 {
			fileBlockPos--
		}
		if fileBlockPos <= 0 {
			break;
		}
		fileStartPos := fileBlockPos + 1
		fileLength := fileEndPos - fileStartPos + 1

		spaceStartPos := -1
		spaceLength := 0
		for pos = 0; pos < fileStartPos; pos++ {
			if blocks2[pos] != -1 {
				spaceStartPos = -1
				spaceLength = 0
			} else {
				if spaceStartPos == -1 {
					spaceStartPos = pos
				}
				spaceLength++
				if spaceLength == fileLength {
					copy(blocks2[spaceStartPos:spaceStartPos + spaceLength + 1], blocks2[fileStartPos:fileEndPos+1])
					for j := fileStartPos; j <= fileEndPos; j++ {
						blocks2[j] = -1
					}
					break
				}
			}
		}
	}
	pos = 0
	total1, total2 := 0, 0
	for blocks[pos] >= 0 {
		total1 += blocks[pos] * pos
		pos++
	}
	for idx, val := range blocks2 {
		if val > 0 {
			total2 += val * idx
		}
	}

	return total1, total2
}

func runForFile(filename string) {
	chcksum, c2 := calc(filename)
	fmt.Printf("File: %s, checksum: %d, : %d\n", filename, chcksum, c2)
}

func main() {
	runForFile("test0.txt")
	runForFile("test.txt")
	runForFile("input.txt")
}