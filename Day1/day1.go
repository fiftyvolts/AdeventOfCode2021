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

	prev, _, err := getnext()
	if err != nil {
		fmt.Println("file errored.")
		return
	}

	count := 0
	for {
		curr, _, err := getnext()
		if err != nil {
			break
		}

		if curr > prev {
			count++
		}

		prev = curr
	}

	fmt.Println(count)
}
