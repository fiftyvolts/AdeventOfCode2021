package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"flag"
	"fmt"
	"math"
	"os"
)

var output *os.File
var verbose bool
var max Point

type Point struct {
	x, y int
}

type Vertex struct {
	point     Point
	risk      int
	totalRisk int
}

type VertexHeap struct {
	Heap []*Vertex
}

func (h *VertexHeap) Len() int {
	return len(h.Heap)
}

func (h *VertexHeap) Less(i int, j int) bool {
	return h.Heap[i].totalRisk < h.Heap[j].totalRisk
}

func (h *VertexHeap) Swap(i, j int) {
	tmp := h.Heap[i]
	h.Heap[i] = h.Heap[j]
	h.Heap[j] = tmp
}

func (h *VertexHeap) Push(x interface{}) {
	h.Heap = append(h.Heap, x.(*Vertex))
}

func (h *VertexHeap) Pop() interface{} {
	n := len(h.Heap) - 1
	tmp := h.Heap[n]
	h.Heap = h.Heap[:n]
	return tmp
}

func main() {
	var inPath, outPath string
	flag.StringVar(&inPath, "input", "ex.txt", "")
	flag.StringVar(&outPath, "output", "out.txt", "")
	flag.BoolVar(&verbose, "verbose", true, "")
	flag.Parse()

	data, err := os.ReadFile(inPath)

	if err != nil {
		panic(err)
	}

	output, err = os.Create(outPath)

	scanner := bufio.NewScanner(bytes.NewReader(data))

	inputData := make(map[Point]*Vertex)
	for y := 0; ; y++ {
		ok := scanner.Scan()
		if !ok {
			if scanner.Err() != nil {
				panic(scanner.Err())
			}
			break
		}

		for x, r := range scanner.Text() {
			point := Point{y, x}
			inputData[point] = &Vertex{point, int(r - '0'), math.MaxInt}
			if x+1 > max.x {
				max.x = x + 1
			}
		}

		if y+1 > max.y {
			max.y = y + 1
		}
	}
	fmt.Fprintln(output, max)

	part1(inputData, max)
	part2(inputData, max)
}

func adj(c Point, max Point) []Point {
	ret := make([]Point, 0)

	if c.x-1 >= 0 {
		ret = append(ret, Point{c.x - 1, c.y})
	}

	if c.x+1 < max.x {
		ret = append(ret, Point{c.x + 1, c.y})
	}

	if c.y-1 >= 1 {
		ret = append(ret, Point{c.x, c.y - 1})
	}

	if c.y+1 < max.y {
		ret = append(ret, Point{c.x, c.y + 1})
	}

	return ret
}

func part1(input map[Point]*Vertex, max Point) {

	cave := make(map[Point]*Vertex)

	for _, vertex := range input {
		vcopy := new(Vertex)
		*vcopy = *vertex
		cave[vcopy.point] = vcopy
	}

	totalRisk := search(cave, max)
	fmt.Fprintln(output, totalRisk)
}

func part2(input map[Point]*Vertex, max Point) {
	cave := make(map[Point]*Vertex)

	expandedMax := Point{0, 0}
	for xshift := 0; xshift < 5; xshift++ {
		for yshift := 0; yshift < 5; yshift++ {
			for _, vertex := range input {
				vcopy := new(Vertex)
				x := vertex.point.x + (xshift * max.x)
				y := vertex.point.y + (yshift * max.y)
				vcopy.point = Point{x, y}
				vcopy.risk = ((vertex.risk + xshift + yshift - 1) % 9) + 1
				vcopy.totalRisk = math.MaxInt
				cave[vcopy.point] = vcopy
			}
		}
	}

	expandedMax = Point{max.x * 5, max.y * 5}
	totalRisk := search(cave, expandedMax)
	fmt.Fprintln(output, totalRisk)
}

func search(cave map[Point]*Vertex, max Point) int {
	visited := make(map[Point]bool)
	start := Point{0, 0}
	end := Point{max.x - 1, max.y - 1}

	if verbose {
		for x := 0; x < max.x; x++ {
			for y := 0; y < max.y; y++ {
				fmt.Fprintf(output, "%d", cave[Point{x, y}].risk)
			}
			fmt.Fprintln(output, "")
		}
	}

	cave[start].totalRisk = 0
	vHeap := new(VertexHeap)
	heap.Init(vHeap)
	heap.Push(vHeap, cave[start])

	for vHeap.Len() > 0 {
		curr := heap.Pop(vHeap).(*Vertex)
		if curr.point == end {
			break
		}

		for _, p := range adj(curr.point, max) {
			next := cave[p]
			visit := visited[next.point]
			if visit {
				var i int
				var v *Vertex
				for i, v = range vHeap.Heap {
					if next == v {
						break
					}
				}
				newTotal := curr.totalRisk + next.risk
				if newTotal < next.totalRisk {
					next.totalRisk = newTotal
					heap.Fix(vHeap, i)
				}
			} else {
				visited[next.point] = true
				next.totalRisk = curr.totalRisk + next.risk
				heap.Push(vHeap, next)
			}
		}
	}

	return cave[end].totalRisk
}
