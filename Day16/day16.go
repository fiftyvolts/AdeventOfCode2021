package main

import (
	"bytes"
	"flag"
	"fmt"
	"math"
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

	reader := bytes.NewReader(data)
	var input []byte

	for {
		var b byte
		var n int
		n, err = fmt.Fscanf(reader, "%2x", &b)

		if err != nil {
			break
		}

		if n != 1 {
			break
		}
		input = append(input, b)
	}
	part1(input)
	part2(input)
}

func mask(width int) int {
	m := 0
	for i := 0; i < width; i++ {
		m <<= 1
		m |= 1
	}
	return m
}

func unpackint(data []byte, width int, offset int) int {
	if width > 16 {
		return -1
	}

	align := offset % 8
	idx := offset / 8
	var value int

	size := 8 - align
	m := mask(size)
	if width+align > 8 {
		value = (int(data[idx]) & m) << (width - size)
	} else {
		value = (int(data[idx]) & m) >> (8 - align - width)
	}

	idx++
	if width-size >= 8 {
		value |= int(data[idx]) << (width - size - 8)
		size += 8
		idx++
	}

	remainder := width - size
	if remainder > 0 {
		value |= int(data[idx]) >> (8 - remainder) & mask(remainder)
	}

	return value
}

func unpack(data []byte, offset int) (*Packet, int) {
	p := new(Packet)

	p.version = byte(unpackint(data, 3, offset))
	offset += 3

	p.packetType = byte(unpackint(data, 3, offset))
	offset += 3

	if p.packetType == 4 {
		for {
			more := unpackint(data, 1, offset)
			offset++

			p.value = p.value<<4 | unpackint(data, 4, offset)
			offset += 4
			if more == 0 {
				break
			}
		}
		return p, offset
	}

	p.mode = byte(unpackint(data, 1, offset))
	offset++

	if p.mode == 1 {
		p.count = unpackint(data, 11, offset)
		offset += 11
		for i := 0; i < p.count; i++ {
			var subp *Packet
			subp, offset = unpack(data, offset)

			p.contents = append(p.contents, subp)
		}
	} else {
		p.length = unpackint(data, 15, offset)
		offset += 15
		for left := p.length; left > 0; {
			var subp *Packet
			start := offset
			subp, offset = unpack(data, offset)
			p.contents = append(p.contents, subp)
			left -= (offset - start)
		}
	}

	switch p.packetType {
	case 0: //sum
		p.value = 0
		for _, subp := range p.contents {
			p.value += subp.value
		}
	case 1: //product
		p.value = 1
		for _, subp := range p.contents {
			p.value *= subp.value
		}
	case 2: //min
		p.value = math.MaxInt
		for _, subp := range p.contents {
			if subp.value < p.value {
				p.value = subp.value
			}
		}
	case 3: //max
		p.value = 0
		for _, subp := range p.contents {
			if subp.value > p.value {
				p.value = subp.value
			}
		}
	case 5: // gt
		if p.contents[0].value > p.contents[1].value {
			p.value = 1
		} else {
			p.value = 0
		}
	case 6: // lt
		if p.contents[0].value < p.contents[1].value {
			p.value = 1
		} else {
			p.value = 0
		}
	case 7: // eq

		if p.contents[0].value == p.contents[1].value {
			p.value = 1
		} else {
			p.value = 0
		}
	}

	return p, offset
}

func part1(data []byte) {
	p, _ := unpack(data, 0)

	var sum int
	stack := []*Packet{p}

	for len(stack) > 0 {
		n := len(stack) - 1
		next := stack[n]
		stack = stack[:n]

		sum += int(next.version)
		for _, sub := range next.contents {
			stack = append(stack, sub)
		}
	}

	fmt.Fprintln(output, sum)
}

func part2(data []byte) {
	p, _ := unpack(data, 0)
	fmt.Fprintln(output, p.value)
}
