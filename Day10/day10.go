package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
)

func main() {
	data, err := os.ReadFile(os.Args[1])

	if err != err {
		fmt.Println(err)
	}

	part1(&data)
}

var ErrorScore = map[rune]int{
	')': 3,
	']': 57,
	'}': 1197,
	'>': 25137}

var CompletionScore = map[rune]int{
	')': 1,
	']': 2,
	'}': 3,
	'>': 4}

var Opens = map[rune]bool{
	'(': true, '[': true, '{': true, '<': true}

var Pairs = map[rune]rune{
	'(': ')', '[': ']', '{': '}', '<': '>'}

func part1(data *[]byte) {
	scanner := bufio.NewScanner(bytes.NewReader(*data))

	errScore := 0
	var completionScores []int
	var corrupt bool
	for {
		corrupt = false
		ok := scanner.Scan()

		if !ok {
			if scanner.Err() != nil {
				fmt.Println(scanner.Err())
			}
			break
		}

		var stack []rune
		for _, r := range scanner.Text() {
			_, ok := Opens[r]
			if ok {
				stack = append(stack, r)
			} else {
				n := len(stack) - 1
				top := stack[n]
				if Pairs[top] == r {
					stack = stack[:n]
				} else {
					score := ErrorScore[r]
					errScore += score
					corrupt = true
					break
				}
			}
		}

		if corrupt {
			continue
		}

		completionScore := 0

		for len(stack) > 0 {
			n := len(stack) - 1
			r := stack[n]
			stack = stack[:n]
			completionScore = (completionScore * 5) + CompletionScore[Pairs[r]]
		}

		if completionScore != 0 {
			completionScores = append(completionScores, completionScore)
		}
	}

	sort.IntSlice(completionScores).Sort()

	midScore := completionScores[(len(completionScores)-1)/2]

	fmt.Printf("ErrorScore: %d CompletionScore: %d %v\n", errScore, midScore, completionScores)
}
