package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func Execute(opCode int, arg int, a, b, c, ip *int) string {
	operand := arg
	if opCode == 0 || opCode == 2 || opCode == 5 || opCode == 6 || opCode == 7 {
		if arg == 4 {
			operand = *a
		} else if arg == 5 {
			operand = *b
		} else if arg == 6 {
			operand = *c
		}	
	}
	switch opCode {
	case 0: //adv
		*a = *a / int(math.Pow(2, float64(operand)))
	case 1: //bxl
		*b = *b ^ operand
	case 2: //bst
		*b = operand % 8
	case 3: //jnz
		if *a != 0 {
			*ip = operand
		}
	case 4: //bxc
		*b = *b ^ *c
	case 5:
		return strconv.Itoa(operand % 8)
	case 6:
		*b = *a / int(math.Pow(2, float64(operand)))
	case 7:
		*c = *a / int(math.Pow(2, float64(operand)))
	}
	return ""
}
//Program: 2,4,1,1,7,5,4,4,1,4,0,3,5,5,3,0

// B = A % 8
// B = B ^ 1
// C = A >> B
// B = B ^ C
// B = B ^ 4
// A = A >> 3
// OUT B % 8
// RESTART while A 

func FindSolution(program []int) int {
	A, B, C := 0, 0, 0
	possibleValues := make([]int,1)
	possibleValues[0] = 0
	for pos := len(program) - 1; pos >= 0; pos-- {
		next := make([]int,0)
		for _ , Ai := range possibleValues { 
			A = Ai << 3
			for i:=0;i<8;i++ {
				if slices.Equal(program[pos:], ExecuteProgram(program, A+i, B, C)) {
					next = append(next, A+i)
					fmt.Printf("%d %d: %+v\r\n", pos, A+i, program[pos:])
				}
			}
		}
		possibleValues = next
	}
	slices.Sort(possibleValues)
	return possibleValues[0]
}

func ExecuteProgram(program []int, a, b, c int) []int {
	ip := 0
	outputs := make([]int, 0)
	for ip < len(program) {
		opCode := program[ip]
		arg := program[ip + 1]
		ip += 2
		if res := Execute(opCode, arg, &a, &b, &c, &ip); res != "" {
			v, _ := strconv.Atoi(res)
			outputs = append(outputs, v)
		}
	}
	return outputs
}

func getOutput(filename string) (string, int) {
	file, _ := os.Open(filename)
	content, _ := io.ReadAll(file)
	file.Close()
	parts := strings.Split(string(content), "\r\n\r\n")
	rows := strings.Split(parts[0], "\r\n")
	var a,b,c int
	fmt.Sscanf(rows[0],"Register A: %d", &a)
	fmt.Sscanf(rows[0],"Register B: %d", &b)
	fmt.Sscanf(rows[0],"Register C: %d", &c)
	programText := strings.Split(parts[1], " ")
	commandTexts := strings.Split(programText[1], ",")
	commands := make([]int, len(commandTexts))
	for i,cmd := range commandTexts {
		commands[i], _ = strconv.Atoi(cmd)
	}
	ip := 0
	outputs := make([]string, 0)
	for ip < len(commands) {
		opCode := commands[ip]
		arg := commands[ip + 1]
		ip += 2
		if res := Execute(opCode, arg, &a, &b, &c, &ip); res != "" {
			outputs = append(outputs, res)
		}
	}
	Amin := FindSolution(commands)

	return strings.Join(outputs, ","), Amin
}

func runForFile(filename string) {
	output, Amin := getOutput(filename)
	fmt.Printf("File: %s, Program output: %s Amin=%d\n", filename, output, Amin)
}

func main() {
	/*runForFile("test.txt")
	runForFile("test2.txt")
	runForFile("test3.txt")*/
	runForFile("input.txt")
}
