package main

import (
	"fmt"
	"os"
)

type CodeLine struct {
	signal [10]string
	output [4]string
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	var lines []CodeLine
	for {
		var line CodeLine
		c, _ := fmt.Fscanf(
			f, "%s %s %s %s %s %s %s %s %s %s | %s %s %s %s\n",
			&line.signal[0],
			&line.signal[1],
			&line.signal[2],
			&line.signal[3],
			&line.signal[4],
			&line.signal[5],
			&line.signal[6],
			&line.signal[7],
			&line.signal[8],
			&line.signal[9],
			&line.output[0],
			&line.output[1],
			&line.output[2],
			&line.output[3])
		if c != 14 {
			break
		}

		lines = append(lines, line)
	}
	part1(lines)
	part2(lines)
}

func part1(lines []CodeLine) {
	count := 0
	for _, line := range lines {
		s2n := mapSignal(line.signal[:])
		for _, digit := range line.output {
			decoded := s2n[segments(digit)]
			if decoded == 1 || decoded == 4 || decoded == 7 || decoded == 8 {
				count++
			}
		}
	}
	fmt.Println(count)
}

func part2(lines []CodeLine) {
	sum := 0
	for _, line := range lines {
		s2n := mapSignal(line.signal[:])
		o := line.output
		num := 1000*s2n[segments(o[0])] + 100*s2n[segments(o[1])] + 10*s2n[segments(o[2])] + s2n[segments(o[3])]
		sum += num
	}
	fmt.Println(sum)
}

func segments(digit string) uint8 {
	g := uint8(0)
	for _, r := range digit {
		g |= 1 << (r - 'a')
	}
	return g
}

func mapSignal(signal []string) map[uint8]int {
	n2s := make(map[int]uint8, 10)

	//find unique numbers first
	for _, digit := range signal {
		seg := segments(digit)
		switch len(digit) {
		case 2:
			n2s[1] = seg
		case 3:
			n2s[7] = seg
		case 4:
			n2s[4] = seg
		case 7:
			n2s[8] = seg
		}
	}

	//find 3 and 6
	for _, digit := range signal {
		seg := segments(digit)
		switch len(digit) {
		case 5:
			if n2s[1]|seg == seg {
				n2s[3] = seg
			}
		case 6:
			if n2s[1]|seg != seg {
				n2s[6] = seg
			}
		}
	}

	//segment isolation
	top := 0x7F & (((^n2s[7]) | n2s[1]) ^ n2s[8])
	tright := 0x7F & (n2s[1] & (^n2s[6]))
	bright := 0x7F & (n2s[1] & (^tright))
	middle := 0x7F & (n2s[4] & (^tright) & (^bright) & (n2s[3]))
	bottom := 0x7F & (n2s[8] & (^n2s[1]) & (n2s[3] & (^top) & (^middle)))
	bleft := 0x7F & ((^n2s[4]) & (^top) & (^bottom))
	tleft := 0x7F & (n2s[8] ^ (top | middle | bottom | tright | bright | bleft))

	//construct remaining map
	n2s[2] = top | tright | middle | bleft | bottom
	n2s[5] = top | tleft | middle | bright | bottom
	n2s[9] = 0x7F ^ bleft

	s2n := make(map[uint8]int, 10)
	for k, v := range n2s {
		s2n[v] = k
	}

	return s2n
}
