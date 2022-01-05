package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	bytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	buf := strings.Builder{}
	buf.Grow(len(bytes))
	buf.Write(bytes)

	var input []int
	for _, v := range strings.Split(buf.String(), ",") {
		i, _ := strconv.Atoi(strings.Trim(v, " \n"))
		input = append(input, i)
	}

	part1(input)
}

func part1(input []int) {
	cycles := [9]int{}

	for _, c := range input {
		cycles[c]++
	}

	for i := 0; i < 256; i++ {
		zero := cycles[0]
		for j := 1; j < len(cycles); j++ {
			cycles[j-1] = cycles[j]
		}
		cycles[6] += zero
		cycles[8] = zero
	}

	sum := 0
	for _, c := range cycles {
		sum += c
	}
	fmt.Println(sum)
}
