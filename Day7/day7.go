package main

import (
	"container/list"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func printList(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Print(e.Value)
		if e.Next() != nil {
			fmt.Print(" ")
		} else {
			fmt.Print("\n")
		}
	}
}

func main() {
	d, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
	var crabPos []int
	for _, v := range strings.Split(string(d), ",") {
		i, _ := strconv.Atoi(strings.TrimFunc(v, unicode.IsSpace))
		crabPos = append(crabPos, i)
	}

	sort.Sort(sort.IntSlice(crabPos))
	part1(crabPos)
	part2(crabPos)
}

func part1(crabPos []int) {
	leftFuel := 0
	rightFuel := 0
	rightCount := len(crabPos)
	leftCount := 0

	for _, p := range crabPos {
		rightFuel += p
	}

	lastFuel := rightFuel
	currPos := 0
	i := 0

	for i < len(crabPos) {
		currCount := 0
		delta := crabPos[i] - currPos
		currPos = crabPos[i]
		for i < len(crabPos) && crabPos[i] == currPos {
			currCount++
			i++
		}

		rightFuel -= rightCount * delta
		leftFuel += leftCount * delta

		if lastFuel < (rightFuel + leftFuel) {
			break
		}

		leftCount += currCount
		rightCount -= currCount
		lastFuel = leftFuel + rightFuel
	}
	fmt.Printf("%d\n", lastFuel)
}

type CrabCost struct {
	pos, cost int
}

func part2(crabPos []int) {
	lastCost := CrabCost{-1, math.MaxInt}
	maxPos := crabPos[len(crabPos)-1]

	for currPos := 0; currPos < maxPos; currPos++ {
		cost := CrabCost{currPos, 0}
		for j := 0; j < len(crabPos); j++ {
			pos := crabPos[j]
			delta := pos - currPos
			if delta < 0 {
				delta = -delta
			}
			if delta&1 != 0 {
				half := (delta - 1) / 2
				cost.cost += delta * (half + 1)
			} else {
				half := delta / 2
				cost.cost += delta*half + half
			}
		}

		if lastCost.cost < cost.cost {
			break
		}
		lastCost = cost
	}

	fmt.Println(lastCost)
}
