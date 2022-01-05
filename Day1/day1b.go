package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("./input.txt")

	if err != nil {
		fmt.Println(err)
		return
	}

	getnext := func() (int, int, error) {
		var x int
		n, err := fmt.Fscanf(f, "%d\n", &x)
		return x, n, err
	}

	win := [4]int{}
	count := 0

	rot := func(x int, y int) int {
		return (x + y) % len(win)
	}

	for pos, seen := 0, 0; ; pos = rot(pos, 1) {
		var err error
		win[pos], _, err = getnext()
		if err != nil {
			break
		}

		seen++
		if seen < 4 {
			continue
		}

		if win[rot(pos, 1)]+win[rot(pos, 2)]+win[rot(pos, 3)] <
			win[rot(pos, 2)]+win[rot(pos, 3)]+win[rot(pos, 4)] {
			count++
		}
	}

	fmt.Println(count)
}
