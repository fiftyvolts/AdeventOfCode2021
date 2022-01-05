package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type Point struct {
	i, j int
}

type Geo struct {
	basin, height int
	point         Point
}

type Grid [][]Geo

func (grid *Grid) getAdj(i int, j int) []Point {
	row := (*grid)[i]
	var points []Point

	if i > 0 {
		points = append(points, Point{i - 1, j})
	}

	if j+1 < len(row) {
		points = append(points, Point{i, j + 1})
	}

	if i < len(*grid)-1 {
		points = append(points, Point{i + 1, j})
	}

	if j > 0 {
		points = append(points, Point{i, j - 1})
	}
	return points
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	bf := bufio.NewScanner(f)

	var grid Grid
	i := 0
	for {
		ok := bf.Scan()

		if !ok {
			if bf.Err() != nil {
				fmt.Println(bf.Err())
			}
			break
		}

		grid = append(grid, make([]Geo, 0))
		row := len(grid) - 1
		for j, r := range bf.Text() {
			grid[row] = append(grid[row], Geo{math.MaxInt, int(r - '0'), Point{i, j}})
		}
	}

	lowPoints := part1(grid)
	part2b(grid, lowPoints)
}

func part1(grid Grid) []Point {
	var sum int
	var lowPoints []Point

	for i, row := range grid {
		for j, v := range row {

			isLow := true
			a := grid[i][j]
			for _, adj := range grid.getAdj(i, j) {
				b := grid[adj.i][adj.j]
				if a.height >= b.height {
					isLow = false
					break
				}
			}
			if isLow {
				sum += v.height + 1
				lowPoints = append(lowPoints, Point{i, j})
			}
		}
	}
	fmt.Println(sum)
	return lowPoints
}

func part2b(grid Grid, lowPoints []Point) {
	nextBasin := 0

	for _, lowPoint := range lowPoints {
		grid[lowPoint.i][lowPoint.j].basin = nextBasin
		nextBasin++

		pointStack := []Point{lowPoint}

		for len(pointStack) > 0 {
			n := len(pointStack) - 1
			point := pointStack[n]
			pointStack = pointStack[:n]
			a := &grid[point.i][point.j]

			for _, adj := range grid.getAdj(point.i, point.j) {
				b := &grid[adj.i][adj.j]

				if b.height != 9 && b.basin == math.MaxInt {
					b.basin = a.basin
					pointStack = append(pointStack, adj)
				}
			}
		}
	}

	basinSizeMap := make(map[int]int, 0)
	for _, row := range grid {
		for _, v := range row {
			if v.basin != math.MaxInt {
				basinSizeMap[v.basin]++
			}
		}
	}

	var basinSizes sort.IntSlice
	for _, v := range basinSizeMap {
		basinSizes = append(basinSizes, v)
	}

	l := len(basinSizes)
	fmt.Println(l)
	basinSizes.Sort()
	fmt.Printf("%d x %d x %d = %d\n",
		basinSizes[l-1], basinSizes[l-2], basinSizes[l-3],
		basinSizes[l-1]*basinSizes[l-2]*basinSizes[l-3])

	//printGrid(&grid, &lowPoints)
	fmt.Printf("lowpoints: %d basins %d\n", len(lowPoints), len(basinSizes))
}

const Alnum = "abcdefghijklmnopqrstuvqxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func printGrid(grid *Grid, lowPoints *[]Point) {
	pMap := make(map[Point]bool)
	for _, p := range *lowPoints {
		pMap[p] = true
	}

	for i, row := range *grid {
		for j, v := range row {
			_, ok := pMap[Point{i, j}]
			if ok {
				fmt.Print("%")
			} else if v.height != 9 {
				fmt.Print(string(Alnum[v.basin%len(Alnum)]))
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println("")
	}
}
