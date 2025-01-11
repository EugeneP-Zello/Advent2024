package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type test struct {
	expected int
	arguments []int
}

func concatenate(a int, b int) int {
	result, _ := strconv.Atoi(fmt.Sprintf("%d%d", a, b))
	return result
}

func (t test) verify1() bool {
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
		if result > t.expected {
			continue
		}
	}
	return false
}


func (t test) verify2() bool {
	if len(t.arguments) == 1 {
		return t.expected == t.arguments[0]	
	}
	opCount := len(t.arguments) - 1
	optionsTotal := int(math.Pow(3, float64(opCount)))
	for i := 0; i < optionsTotal; i++ {
		operations := strconv.FormatInt(int64(i), 3)
		if zerosNeeded := opCount - len(operations); zerosNeeded > 0 {
			operations = strings.Repeat("0", zerosNeeded) + operations
		}
		result := t.arguments[0]
		for op := 0; op < opCount; op++ {
			if operations[op] == '0'{
				result += t.arguments[op + 1]
			} else if operations[op] == '1' {
				result *= t.arguments[op + 1]
 			} else {
				result = concatenate(result, t.arguments[op + 1])
			}
		}
		if result == t.expected {
			return true
		}
		if result > t.expected {
			continue
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

func checkEquastion(s string) (int, int) {
	t := newTest(s)
	result1, result2 := 0, 0
	if t.verify2() {
		result2 = t.expected
	}
	if t.verify1() {
		result1 = t.expected
	}
	return result1, result2
}

func calc(filename string) (int, int) {
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	total1, total2 := 0, 0
	for scanner.Scan() {
		r1, r2 := checkEquastion(scanner.Text())
		total1 += r1
		total2 += r2
	}
	return total1, total2
}

func runForFile(filename string) {
	c, lc := calc(filename)
	fmt.Printf("File: %s, count: %d, loop count: %d\n", filename, c, lc)
}

func main() {
	runForFile("test.txt")
	runForFile("input.txt")
}
