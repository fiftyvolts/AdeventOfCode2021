package main

import (
	"fmt"
	"os"
)

func main() {
	part1(os.Args[1])
}

func part1(path string) {
	f, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		return
	}

	rows := 0
	var counts []int
	var bits int
	for {
		var s string
		_, err := fmt.Fscanf(f, "%s\n", &s)

		if err != nil {
			break
		}

		bits = len(s)
		if counts == nil {
			counts = make([]int, bits)
		}
		rows++
		for i, v := range s {
			counts[bits-i-1] += (map[rune]int{'1': 1, '0': 0})[v]
		}
	}

	gamma := 0
	epsilon := 0
	for i, v := range counts {
		if v >= rows/2.0 {
			gamma |= 1 << i
		} else {
			epsilon |= 1 << i
		}

		fmt.Printf("%5d %5d %5d = %016b\n", i, v, rows, gamma)
	}

	fmt.Printf("%d(%012b) %d(%012b) = %d\n", gamma, gamma, epsilon, epsilon, gamma*epsilon)

}
