package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Node struct {
	id    string
	adj   []*Node
	small bool
}

type Graph struct {
	start, end *Node
}

func NewGraphFromFile(path string) *Graph {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))

	graph := new(Graph)
	lut := make(map[string]*Node)

	for {
		ok := scanner.Scan()
		if !ok {
			if scanner.Err() != nil {
				panic(scanner.Err())
			}
			break
		}

		pair := strings.Split(scanner.Text(), "-")

		p, ok := lut[pair[0]]
		if !ok {
			p = new(Node)
			p.id = pair[0]
			p.adj = make([]*Node, 0)
			p.small, _ = regexp.Match("[a-z]+", []byte(p.id))
			lut[pair[0]] = p
		}

		c, ok := lut[pair[1]]
		if !ok {
			c = new(Node)
			c.id = pair[1]
			c.adj = make([]*Node, 0)
			c.small, _ = regexp.Match("[a-z]+", []byte(c.id))
			lut[pair[1]] = c
		}

		p.adj = append(p.adj, c)
		c.adj = append(c.adj, p)

		if p.id == "start" && graph.start == nil {
			graph.start = p
		} else if c.id == "start" && graph.start == nil {
			graph.start = c
		}

		if p.id == "end" && graph.end == nil {
			graph.end = p
		} else if c.id == "end" && graph.end == nil {
			graph.end = c
		}
	}
	return graph
}

func (g *Graph) String() string {
	b := strings.Builder{}

	seen := make(map[*Node]bool)
	queue := []*Node{g.start}

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]

		seen[p] = true

		b.WriteString(p.id)
		b.WriteString("(")
		b.WriteString(strconv.FormatBool(p.small))
		b.WriteString("): ")
		for i, c := range p.adj {
			b.WriteString(c.id)
			if i < len(p.adj)-1 {
				b.WriteString(",")
			} else {
				b.WriteString("\n")
			}

			_, ok := seen[c]
			if !ok {
				seen[c] = true
				queue = append(queue, c)
			}
		}

	}
	return b.String()
}

func PathString(nodes *[]*Node) string {
	b := strings.Builder{}
	b.WriteString("[")
	for i, n := range *nodes {
		b.WriteString(n.id)
		if i < len(*nodes)-1 {
			b.WriteString(",")
		} else {
			b.WriteString("]")
		}
	}
	return b.String()
}

func (g *Graph) Paths(part2 bool) [][]*Node {
	var paths [][]*Node
	var stack [][]*Node

	stack = append(stack, []*Node{g.start})

	for len(stack) > 0 {
		n := len(stack) - 1
		currPath := stack[n]
		last := currPath[len(currPath)-1]
		stack = stack[:n]

		if last == g.end {
			tmp := make([]*Node, len(currPath))
			copy(tmp, currPath)
			paths = append(paths, tmp)
		} else {
		NEXTLOOP:
			for _, next := range last.adj {
				if next == g.start {
					continue
				}

				if next.small {
					smallCounts := make(map[*Node]int)
					smallCounts[next]++

					for _, check := range currPath {
						if check.small {
							smallCounts[check]++
						}
					}

					if part2 {
						if smallCounts[next] > 2 {
							continue NEXTLOOP
						}

						numTwo := 0
						for _, v := range smallCounts {
							if v == 2 {
								numTwo++
							}
						}
						if numTwo > 1 {
							continue NEXTLOOP
						}
					} else {
						if smallCounts[next] == 2 {
							continue NEXTLOOP
						}
					}
				}
				tmp := make([]*Node, len(currPath))
				copy(tmp, currPath)
				tmp = append(tmp, next)
				stack = append(stack, tmp)
			}
		}
	}

	return paths
}

func main() {
	var input, output string
	flag.StringVar(&input, "input", "ex.txt", "")
	flag.StringVar(&output, "output", "out.txt", "")
	flag.Parse()

	f, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	graph := NewGraphFromFile(input)
	fmt.Fprint(f, graph)

	paths := graph.Paths(false)
	for _, path := range paths {
		if false {
			fmt.Fprintln(f, PathString(&path))
		}
	}
	fmt.Fprintln(f, len(paths))

	paths = graph.Paths(true)
	for _, path := range paths {
		if false {
			fmt.Fprintln(f, PathString(&path))
		}
	}
	fmt.Fprintln(f, len(paths))
}
