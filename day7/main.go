package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type test struct {
	expected int
	arguments []int
}

func (t test) verify() bool {
	if len(t.arguments) == 1 {
		return t.expected == t.arguments[0]	
	}
	opCount := len(t.arguments) - 1
	optionsTotal := 1 << opCount
	for i := 0; i < optionsTotal; i++ {
		result := t.arguments[0]
		for op := 0; op < opCount; op++ {
			if i & (1 << op) != 0 {
				result += t.arguments[op + 1]
			} else {
				result *= t.arguments[op + 1]
			}
		}
		if result == t.expected {
			return true
		}
	}
	return false
}

func newTest(s string) test {
	s2 := strings.Split(s, ":")
	expected, _ := strconv.Atoi(s2[0])
	s3 := strings.Split(strings.TrimLeft(s2[1], " "), " ")
	t := test{
		expected: expected,
		arguments: make([]int, len(s3)),
	}
	for i, v := range s3 {
		t.arguments[i], _ = strconv.Atoi(v)
	}
	return t
}

func checkEquastion(s string) int {
	t := newTest(s)
	if t.verify() {
		return t.expected
	}
	return 0
}

func calc(filename string) (int, int) {
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	total := 0
	for scanner.Scan() {
		total += checkEquastion(scanner.Text())
	}
	return total, 0
}

func runForFile(filename string) {
	c, lc := calc(filename)
	fmt.Printf("File: %s, count: %d, loop count: %d\n", filename, c, lc)
}

func main() {
	runForFile("test.txt")
	runForFile("input.txt")
}
