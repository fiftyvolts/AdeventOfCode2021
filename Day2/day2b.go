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

	x := 0
	y := 0
	aim := 0

	for {
		var cmd string
		var v int
		_, err = fmt.Fscanf(f, "%s %d\n", &cmd, &v)

		if err != nil {
			break
		}

		switch cmd {
		case "forward":
			x += v
			y += aim * v
		case "down":
			aim += v
		case "up":
			aim -= v
		}
	}

	fmt.Printf("%d, %d = %d\n", x, y, x*y)
}
