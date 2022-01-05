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

var output *os.File
var verbose bool

type Pair struct {
	a, b, parent *Pair
	value        int
}

func (root *Pair) explode() bool {
	var checkDepth func(curr *Pair, depth int) bool
	checkDepth = func(curr *Pair, depth int) bool {
		if curr.a == nil && curr.b == nil {
			return false
		}

		if depth == 4 {
			return true
		}

		return checkDepth(curr.a, depth+1) || checkDepth(curr.b, depth+1)
	}

	shouldExplode := checkDepth(root, 1)

	return shouldExplode
}

func (root *Pair) Copy() *Pair {
	copyRoot := new(Pair)
	if root.a != nil {
		copyRoot.a = root.a.Copy()
		copyRoot.a.parent = copyRoot
	}

	if root.b != nil {
		copyRoot.b = root.b.Copy()
		copyRoot.b.parent = copyRoot
	}

	if root.a == nil && root.b == nil {
		copyRoot.value = root.value
	}

	return copyRoot
}

func (root *Pair) String() string {

	if root.a == nil && root.b == nil {
		return strconv.Itoa(root.value)
	}

	b := strings.Builder{}
	b.WriteString("[")
	b.WriteString(root.a.String())
	b.WriteString(",")
	b.WriteString(root.b.String())
	b.WriteString("]")
	return b.String()
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

	scanner := bufio.NewScanner(bytes.NewReader(data))
	snailnums := make([]*Pair, 0)
	i := 0
	for {
		ok := scanner.Scan()
		if !ok {
			if scanner.Err() != nil {
				panic(scanner.Err())
			}
			break
		}

		line := scanner.Text()
		j := 1
		snailnums = append(snailnums, new(Pair))
		curr := snailnums[i]

		for j < len(line) {
			switch line[j] {
			case '[':
				child := new(Pair)
				child.parent = curr
				if curr.a == nil {
					curr.a = child
				} else {
					curr.b = child
				}
				curr = child
			case ',':
				//do nothing

			case ']':
				curr = curr.parent

			case '\n':
				break

			default: // it's a number
				child := new(Pair)
				child.parent = curr
				child.value = int(line[j] - '0')
				if curr.a == nil {
					curr.a = child
				} else {
					curr.b = child
				}
			}
			j++
		}

		if verbose {
			fmt.Fprintln(output, snailnums[i])
		}
	}

	//part1(snailnums)
}

func part1(snailnums []*Pair) {
	/*depth := 0
	var stack []*Pair
	totalsum := 0

	for _, root := range snailnums {

	}*/

	fmt.Fprintln(output, snailnums[0].explode())
}
