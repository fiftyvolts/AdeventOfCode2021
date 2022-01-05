package main

import (
	"flag"
	"fmt"
	"os"
)

var output *os.File
var verbose bool

type Packet struct {
	version, packetType byte
	mode                byte
	value               int
	length              int
	count               int
	contents            []*Packet
}

type Point struct {
	x, y int
}
type Rect struct {
	bottomLeft, topRight Point
}

func main() {
	var inPath, outPath string
	flag.StringVar(&inPath, "input", "ex.txt", "")
	flag.StringVar(&outPath, "output", "out.txt", "")
	flag.BoolVar(&verbose, "verbose", true, "")
	flag.Parse()

	output, _ = os.Create(outPath)

	data, err := os.ReadFile(inPath)

	if err != nil {
		panic(err)
	}

	var target Rect
	fmt.Sscanf(string(data), "target area: x=%d..%d, y=%d..%d\n",
		&target.bottomLeft.x, &target.topRight.x, &target.bottomLeft.y, &target.topRight.y)

	fmt.Fprintln(output, target)
	part1(target)
	part2(target)
}

func probe(vix int, viy int, target Rect) (int, bool) {
	p := Point{0, 0}
	t := 0

	vx := vix
	vy := viy
	ymax := 0

	for {
		if p.x >= target.bottomLeft.x && p.x <= target.topRight.x &&
			p.y >= target.bottomLeft.y && p.y <= target.topRight.y {

			if verbose {
				fmt.Fprintf(output, "p=%v, vx,vy=%d,%d ymax=%d\n", p, vix, viy, ymax)
			}
			return ymax, true
		} else if (p.x < target.bottomLeft.x && vx == 0) ||
			p.x > target.topRight.x ||
			(p.y < target.bottomLeft.y && vy < 0) {

			break
		}

		t++
		p.x += vx
		p.y += vy
		if p.y > ymax {
			ymax = p.y
		}

		if vx > 0 {
			vx--
		}
		vy--
	}

	return -1, false
}

func part1(target Rect) {
	ymax := 0
	for vix := 0; vix <= 200; vix++ {
		for viy := -200; viy <= 200; viy++ {
			y, ok := probe(vix, viy, target)
			if ok && y > ymax {
				ymax = y
			}
		}
	}
	fmt.Fprintln(output, ymax)
}

func part2(target Rect) {
	ymax := 0
	total := 0

	for vix := 0; vix <= 5000; vix++ {
		for viy := -5000; viy <= 5000; viy++ {
			y, ok := probe(vix, viy, target)
			if ok {
				total++
			}

			if ok && y > ymax {
				ymax = y
			}
		}
	}

	fmt.Fprintln(output, total, ymax)
}
