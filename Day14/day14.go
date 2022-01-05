package main

import (
	"bufio"
	"bytes"
	"container/list"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
)

var verbose bool
var output *os.File
var reps int

func PrettyList(list *list.List) string {
	b := strings.Builder{}
	b.Reset()
	for e := list.Front(); e != nil; e = e.Next() {
		b.WriteRune(e.Value.(rune))
	}
	return b.String()
}

func main() {
	var inPath, outPath string
	flag.StringVar(&inPath, "input", "ex.txt", "")
	flag.StringVar(&outPath, "output", "out.txt", "")
	flag.BoolVar(&verbose, "verbose", true, "")
	flag.IntVar(&reps, "reps", 10, "")
	flag.Parse()

	data, err := os.ReadFile(inPath)

	if err != nil {
		panic(err)
	}

	output, err = os.Create(outPath)

	scanner := bufio.NewScanner(bytes.NewReader(data))

	ok := scanner.Scan()

	if !ok {
		if scanner.Err() != nil {
			panic(scanner.Err())
		}
	}

	polymer := new(list.List)

	for _, r := range scanner.Text() {
		polymer.PushBack(r)
	}

	//skip a line
	ok = scanner.Scan()
	if !ok {
		if scanner.Err() != nil {
			panic(scanner.Err())
		}
	}

	pairs := make(map[[2]rune]rune)
	for {
		ok = scanner.Scan()
		if !ok {
			if scanner.Err() != nil {
				panic(scanner.Err())
			}
			break
		}

		var p [2]rune
		var insertion rune
		fmt.Sscanf(scanner.Text(), "%c%c -> %c", &(p[0]), &(p[1]), &insertion)
		pairs[p] = insertion
	}

	//part1(polymer, pairs)
	part2(polymer, pairs)

}

func part1(polymer *list.List, pairs map[[2]rune]rune) {

	if verbose {
		fmt.Fprintln(output, PrettyList(polymer))
	}

	for i := 0; i < reps; i++ {
		first := polymer.Front()
		second := polymer.Front().Next()

		for second != nil {
			polymer.InsertAfter(pairs[[2]rune{first.Value.(rune), second.Value.(rune)}], first)
			first = second
			second = second.Next()
		}

		if verbose {
			fmt.Fprintln(output, PrettyList(polymer))
		}
	}

	freq := make(map[rune]int)
	for e := polymer.Front(); e != nil; e = e.Next() {
		freq[e.Value.(rune)]++
	}

	low := math.MaxInt
	high := 0
	for _, i := range freq {
		if i > high {
			high = i
		}

		if i < low {
			low = i
		}
	}
	fmt.Fprintln(output, high, low, high-low)
}

type Partial struct {
	r    rune
	reps int
}

func PrettyPartials(partials []Partial) string {
	var rs []rune
	for i := len(partials) - 1; i >= 0; i-- {
		rs = append(rs, partials[i].r)
	}
	return string(rs)
}

type DescentData struct {
	first, last rune
	reps        int
}

func part2(polymer *list.List, pairs map[[2]rune]rune) {

	freq := make(map[rune]int)
	memo := make(map[DescentData]map[rune]int)
	memoVerb := make(map[DescentData]string)

	var descend func(dd DescentData) map[rune]int
	descend = func(dd DescentData) map[rune]int {
		memoret, ok := memo[dd]
		if ok {
			if verbose {
				fmt.Fprintf(output, memoVerb[dd])
			}
			return memoret
		}

		ret := make(map[rune]int)

		if dd.reps == 0 {
			if verbose {
				memoVerb[dd] = string(dd.first)
				fmt.Fprintf(output, "%c", dd.first)
			}
			ret[dd.first] = 1
			memo[dd] = ret
			return ret
		} else {
			next := pairs[[2]rune{dd.first, dd.last}]
			for k, v := range descend(DescentData{dd.first, next, dd.reps - 1}) {
				ret[k] += v
			}

			for k, v := range descend(DescentData{next, dd.last, dd.reps - 1}) {
				ret[k] += v
			}
		}

		memo[dd] = ret
		return ret
	}

	for e := polymer.Front(); e.Next() != nil; e = e.Next() {
		dd := DescentData{e.Value.(rune), e.Next().Value.(rune), reps}
		for k, v := range descend(dd) {
			freq[k] += v
		}

	}

	if verbose {
		fmt.Fprintf(output, "%c", polymer.Back().Value.(rune))
	}
	freq[polymer.Back().Value.(rune)]++

	low := math.MaxInt
	high := 0
	for _, i := range freq {
		if i > high {
			high = i
		}

		if i < low {
			low = i
		}
	}
	if verbose {
		fmt.Fprintln(output, "")
	}
	fmt.Fprintln(output, high, low, high-low)
}
