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
	var data []int
	var s string

	_, err = fmt.Fscanf(f, "%s\n", &s)

	if err != nil {
		fmt.Println(err)
		return
	}

	size := len(s)
	msb := size - 1
	counts := make([]int, size)
	f.Seek(0, 0)

	for {
		var x int
		_, err := fmt.Fscanf(f, "%b\n", &x)

		if err != nil {
			break
		}

		rows++

		for i := 0; i < size; i++ {
			if x&(1<<i) != 0 {
				counts[i] += 1
			}
		}
		data = append(data, x)
	}

	gamma := 0
	epsilon := 0
	for i, v := range counts {
		if v >= rows/2.0 {
			gamma |= 1 << i
		} else {
			epsilon |= 1 << i
		}
	}

	fmt.Printf("%5d(%012b) %5d(%012b) = %d\n", gamma, gamma, epsilon, epsilon, gamma*epsilon)

	filter := func(rows int, data []int, most bool) int {
		set := make(map[int]bool, rows)

		for _, v := range data {
			set[v] = true
		}
		for i := msb; i >= 0; i-- {
			count, items := 0, 0

			for v := range set {
				items++
				if v&(1<<i) != 0 {
					count++
				}
			}
			var keep int
			if most {
				if count*2 >= items {
					keep = 1
				} else {
					keep = 0
				}
			} else {
				if count*2 >= items {
					keep = 0
				} else {
					keep = 1
				}
			}

			for v := range set {
				if v&(1<<i) != keep<<i {
					delete(set, v)
				}
			}

			if len(set) == 1 {
				break
			}
		}

		for v := range set {
			return v
		}
		return -1
	}

	o2 := filter(rows, data, true)
	co2 := filter(rows, data, false)
	fmt.Printf("%5d(%012b) %5d(%012b) = %d\n", o2, o2, co2, co2, o2*co2)
}
