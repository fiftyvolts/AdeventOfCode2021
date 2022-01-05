package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Point struct {
	x, y int
}

type Fold struct {
	dir  rune
	line int
}

type Paper struct {
	xmax, ymax int
	points     map[Point]bool
	folds      []Fold
}

func NewPaperFromBytes(b []byte) *Paper {
	paper := new(Paper)
	paper.points = make(map[Point]bool)
	scanner := bufio.NewScanner(bytes.NewReader(b))

	for {
		ok := scanner.Scan()
		if !ok {
			if scanner.Err() != nil {
				panic(scanner.Err())
			}
			break
		}

		line := scanner.Text()
		if strings.Index(line, "fold") == 0 {
			var fold Fold
			fmt.Sscanf(line, "fold along %c=%d", &fold.dir, &fold.line)
			paper.folds = append(paper.folds, fold)
		} else if line == "" {
			continue
		} else {
			var point Point
			fmt.Sscanf(line, "%d,%d", &point.x, &point.y)
			paper.points[point] = true
			if point.x > paper.xmax {
				paper.xmax = point.x
			}

			if point.y > paper.ymax {
				paper.ymax = point.y
			}
		}
	}

	return paper
}

func (paper *Paper) Fold() bool {
	if len(paper.folds) == 0 {
		return false
	}

	newPoints := make(map[Point]bool)
	var xmax, ymax int
	fold := paper.folds[0]
	paper.folds = paper.folds[1:]
	for point := range paper.points {
		var tmp Point
		if fold.dir == 'x' {
			if point.x > fold.line {
				tmp = Point{fold.line - (point.x - fold.line), point.y}
			} else if point.x < fold.line {
				tmp = point
			} else {
				continue
			}
		} else {
			if point.y > fold.line {
				tmp = Point{point.x, fold.line - (point.y - fold.line)}
			} else if point.y < fold.line {
				tmp = point
			} else {
				continue
			}
		}

		newPoints[tmp] = true
		if tmp.x > xmax {
			xmax = point.x
		}

		if tmp.y > ymax {
			ymax = point.y
		}

	}
	paper.points = newPoints
	paper.xmax = xmax
	paper.ymax = ymax
	return true
}

func (paper *Paper) Size() int {
	return len(paper.points)
}

func (paper *Paper) String() string {
	b := strings.Builder{}

	for y := 0; y <= paper.ymax; y++ {
		for x := 0; x <= paper.xmax; x++ {
			_, ok := paper.points[Point{x, y}]
			if ok {
				b.WriteRune('#')
			} else {
				b.WriteRune('.')
			}
		}
		if y < paper.ymax {
			b.WriteRune('\n')
		}
	}
	return b.String()
}

var outfile *os.File

func main() {
	var input, output string
	var printPoints bool
	flag.StringVar(&input, "input", "ex.txt", "")
	flag.StringVar(&output, "output", "out.txt", "")
	flag.BoolVar(&printPoints, "points", true, "")
	flag.Parse()

	f, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	outfile = f

	var b []byte
	b, err = os.ReadFile(input)
	if err != nil {
		panic(err)
	}
	paper := NewPaperFromBytes(b)

	for paper.Fold() {
	}
	if printPoints {
		fmt.Fprintln(outfile, paper)
	}
	fmt.Fprintln(outfile, paper.Size())

}
