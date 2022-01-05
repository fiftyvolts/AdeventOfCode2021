package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Cave struct {
	ilen, jlen int
	nrg, acc   []int
}

type Point struct {
	i, j int
}

func NewCaveFromFile(path string) *Cave {
	c := new(Cave)

	data, err := os.ReadFile(path)

	if err != nil {
		fmt.Println(err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for {
		ok := scanner.Scan()
		if !ok {
			if scanner.Err() != nil {
				fmt.Println(scanner.Err())
			}
			break
		}
		if scanner.Text() != "" {
			c.ilen++
		}
		for _, r := range scanner.Text() {
			c.nrg = append(c.nrg, int(r-'0'))
		}
	}
	c.jlen = len(c.nrg) / c.ilen
	c.acc = make([]int, len(c.nrg))
	return c
}

func (c *Cave) P2I(p Point) int {
	return p.i*c.jlen + p.j
}

func (c *Cave) GetEnergy(p Point) int {
	return c.nrg[c.P2I(p)]
}

func (c *Cave) SetEnergy(p Point, v int) {
	c.nrg[c.P2I(p)] = v
}

func (c *Cave) GetAcc(p Point) int {
	return c.acc[c.P2I(p)]
}

func (c *Cave) SetAcc(p Point, v int) {
	c.acc[c.P2I(p)] = v
}

func (c *Cave) Adj(p Point) []Point {
	var adj []Point
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if p.i+i >= 0 && p.i+i < c.ilen && p.j+j >= 0 && p.j+j < c.jlen {
				adj = append(adj, Point{p.i + i, p.j + j})
			}
		}
	}
	return adj
}

func (c *Cave) String() string {
	b := strings.Builder{}

	b.WriteString("ilen: ")
	b.WriteString(strconv.FormatInt(int64(c.ilen), 10))
	b.WriteString(" jlen: ")
	b.WriteString(strconv.FormatInt(int64(c.jlen), 10))
	b.WriteString("\n")

	for p := range c.Points() {
		b.WriteString(strconv.FormatInt(int64(c.GetEnergy(p)), 10))
		if p.i < c.ilen-1 && p.j+1 == c.jlen {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (c *Cave) Clear() {
	c.nrg = c.nrg[:0]
}

func (c *Cave) Points() chan Point {
	channel := make(chan Point)
	go (func() {
		for i := 0; i < c.ilen; i++ {
			for j := 0; j < c.jlen; j++ {
				channel <- Point{i, j}
			}
		}
		close(channel)
	})()
	return channel
}

func (c *Cave) AccString() string {
	b := strings.Builder{}

	b.WriteString("ilen: ")
	b.WriteString(strconv.FormatInt(int64(c.ilen), 10))
	b.WriteString(" jlen: ")
	b.WriteString(strconv.FormatInt(int64(c.jlen), 10))
	b.WriteString("\n")

	for p := range c.Points() {
		b.WriteString(strconv.FormatInt(int64(c.GetAcc(p)), 10))
		if p.i < c.ilen-1 && p.j+1 == c.jlen {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (c *Cave) Size() int {
	return len(c.nrg)
}

func (c *Cave) Flash() int {
	for i := range c.acc {
		c.acc[i] = 1
	}

	flashed := make(map[Point]bool)
	for {
		again := false

		for p := range c.Points() {
			if !flashed[p] && c.GetEnergy(p)+c.GetAcc(p) > 9 {
				flashed[p] = true
				again = true
				for _, adj := range c.Adj(p) {
					c.SetAcc(adj, c.GetAcc(adj)+1)
				}
			}
		}

		if !again {
			break
		}
	}

	for p := range c.Points() {
		x := c.GetEnergy(p) + c.GetAcc(p)
		if x > 9 {
			x = 0
		}
		c.SetEnergy(p, x)
	}

	return len(flashed)
}

func main() {
	var input, output string
	flag.StringVar(&input, "input", "ex.txt", "")
	flag.StringVar(&output, "output", "out.txt", "")
	flag.Parse()

	cave := NewCaveFromFile(input)
	f, err := os.Create(output)
	if err != nil {
		fmt.Println(err)
	}

	var count, step int
	for i := 0; i < 100; i++ {
		step++
		count += cave.Flash()
	}
	fmt.Fprintf(f, "%s\nFlashes: %d\n", cave, count)

	for {
		step++
		if cave.Flash() == cave.Size() {
			break
		}
	}
	fmt.Fprintf(f, "%s\nStep: %d\n", cave, step)
}
