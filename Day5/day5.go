package main

import (
	"fmt"
	"os"
)

type Point struct {
	x, y int
}
type Vent struct {
	p1, p2 Point
}

func main() {
	vents := readVents()
	part1(vents)
	part2(vents)
}

func readVents() []Vent {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return nil
	}
	vents := []Vent{}
	for {
		var v Vent
		c, _ := fmt.Fscanf(f, "%d,%d -> %d,%d",
			&v.p1.x, &v.p1.y, &v.p2.x, &v.p2.y)

		if c != 4 {
			break
		}
		vents = append(vents, v)
	}
	return vents
}

func part1(vents []Vent) {
	var vert, horiz []Vent
	var xdim, ydim int

	for _, v := range vents {
		if v.p1.x > xdim {
			xdim = v.p1.x
		}
		if v.p2.x > xdim {
			xdim = v.p2.x
		}
		if v.p1.y > ydim {
			ydim = v.p1.y
		}
		if v.p2.y > ydim {
			ydim = v.p2.y
		}

		if v.p1.x == v.p2.x {
			if v.p1.y <= v.p2.y {
				vert = append(vert, v)
			} else {
				vert = append(vert, Vent{v.p2, v.p1})
			}
		} else if v.p1.y == v.p2.y {
			if v.p1.x <= v.p2.x {
				horiz = append(horiz, v)
			} else {
				horiz = append(horiz, Vent{v.p2, v.p1})
			}
		}
	}

	xdim++
	ydim++

	points := make([]int, xdim*ydim)
	for _, v := range vert {
		for i := v.p1.y; i <= v.p2.y; i++ {
			points[i*xdim+v.p1.x]++
		}
	}

	for _, v := range horiz {
		for i := v.p1.x; i <= v.p2.x; i++ {
			points[v.p1.y*xdim+i]++
		}
	}

	//printPoints(points, xdim, ydim)
	printPointCount(points)
}

func part2(vents []Vent) {
	var vert, horiz, diagD, diagU []Vent
	var xdim, ydim int

	for _, v := range vents {
		if v.p1.x > xdim {
			xdim = v.p1.x
		}
		if v.p2.x > xdim {
			xdim = v.p2.x
		}
		if v.p1.y > ydim {
			ydim = v.p1.y
		}
		if v.p2.y > ydim {
			ydim = v.p2.y
		}

		if v.p1.x == v.p2.x {
			if v.p1.y <= v.p2.y {
				vert = append(vert, v)
			} else {
				vert = append(vert, Vent{v.p2, v.p1})
			}
		} else if v.p1.y == v.p2.y {
			if v.p1.x <= v.p2.x {
				horiz = append(horiz, v)
			} else {
				horiz = append(horiz, Vent{v.p2, v.p1})
			}
		} else {
			nv := v
			if v.p1.x > v.p2.x {
				nv = Vent{v.p2, v.p1}
			}

			if nv.p1.y > nv.p2.y {
				diagU = append(diagU, nv)
			} else {
				diagD = append(diagD, nv)
			}
		}
	}

	xdim++
	ydim++

	points := make([]int, xdim*ydim)
	for _, v := range vert {
		for i := v.p1.y; i <= v.p2.y; i++ {
			points[i*xdim+v.p1.x]++
		}
	}

	for _, v := range horiz {
		for i := v.p1.x; i <= v.p2.x; i++ {
			points[v.p1.y*xdim+i]++
		}
	}

	for _, v := range diagD {
		i := v.p1.y
		for j := v.p1.x; j <= v.p2.x; j++ {
			points[i*xdim+j]++
			i++
		}
	}

	for _, v := range diagU {
		i := v.p1.y
		for j := v.p1.x; j <= v.p2.x; j++ {
			points[i*xdim+j]++
			i--
		}
	}

	//printPoints(points, xdim, ydim)
	printPointCount(points)
}

func printPointCount(points []int) {
	count := 0
	for _, v := range points {
		if v > 1 {
			count++
		}
	}
	fmt.Println(count)
}

func printPoints(points []int, xdim int, ydim int) {
	fmt.Printf("%d x %d\n", xdim, ydim)
	for i := 0; i < len(points); i += xdim {
		for j := i; j < i+xdim-1; j++ {
			fmt.Printf("%d ", points[j])
		}
		fmt.Printf("%d\n", points[i+xdim-1])
	}
}
