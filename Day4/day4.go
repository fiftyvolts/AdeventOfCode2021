package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	numbers, boards := getNumbersBoards()
	if numbers == nil {
		return
	}
	part1(numbers, boards)
	part2(numbers, boards)
}

func getNumbersBoards() ([]int, []int) {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	b := bufio.NewScanner(f)
	ok := b.Scan()
	if !ok {
		if b.Err() != nil {
			fmt.Println(b.Err())
		}
	}

	var numbers []int

	for _, v := range strings.Split(b.Text(), ",") {
		x, _ := strconv.Atoi(v)
		numbers = append(numbers, x)
	}

	//load the boards and map to all the numbers
	var boards []int

	for {
		ok := b.Scan()
		if !ok {
			break
		}

		line := b.Text()
		if line == "" {
			continue
		}

		x := make([]int, 5)
		fmt.Sscanf(b.Text(), "%2d %2d %2d %2d %2d", &x[0], &x[1], &x[2], &x[3], &x[4])
		boards = append(boards, x...)
	}

	return numbers, boards
}

func printBoards(boards []int) {
	var i, j int
	for i = 0; i*5+j < len(boards); i++ {
		for j = 0; j < 5 && i*5+j < len(boards); j++ {
			fmt.Printf("%2d ", boards[i*5+j])
		}
		if i%5 == 4 {
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
	}

}

func getValue(boards []int, board int, row int, column int) int {
	return boards[board*25+row*5+column]
}

func checkBoard(boards []int, called map[int]bool, board int) bool {
	//rows
	for i := 0; i < 5; i++ {
		count := 0
		for j := 0; j < 5; j++ {
			if called[getValue(boards, board, i, j)] {
				count++
			} else {
				break
			}
		}
		if count == 5 {
			return true
		}
	}

	//cols
	for j := 0; j < 5; j++ {
		count := 0
		for i := 0; i < 5; i++ {
			if called[getValue(boards, board, i, j)] {
				count++
			} else {
				break
			}
		}
		if count == 5 {
			return true
		}
	}

	/*diagonal
	count := 0
	for i := 0; i < 5; i++ {
		if called[getValue(boards, board, i, i)] {
			count++
		}
	}
	if count == 5 {
		return true
	}

	count = 0
	for i := 0; i < 5; i++ {
		if called[getValue(boards, board, i, 4-i)] {
			count++
		}
	}
	if count == 5 {
		return true
	}*/

	return false
}

func checkBoards(boards []int, called map[int]bool) (int, bool) {
	numBoards := len(boards) / 25
	for i := 0; i < numBoards; i++ {
		if checkBoard(boards, called, i) {
			return i, true
		}
	}
	return -1, false
}

func scoreBoard(boards []int, called map[int]bool, board int, last int) {
	sum := 0
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			x := getValue(boards, board, i, j)
			_, ok := called[x]
			if !ok {
				sum += x
			}
		}
	}

	fmt.Printf("Score %d x %d = %d\n", sum, last, sum*last)
}

func part1(numbers []int, boards []int) {
	called := map[int]bool{}
	var board, last int

	for _, n := range numbers {
		called[n] = true
		last = n
		var bingo bool
		board, bingo = checkBoards(boards, called)
		if bingo {
			break
		}
	}

	scoreBoard(boards, called, board, last)

	return
}

func part2(numbers []int, boards []int) {
	called := map[int]bool{}
	var board, last int
	boardCount := len(boards) / 25
	winners := map[int]bool{}

CALLS:
	for _, n := range numbers {
		called[n] = true
		last = n
		var bingo bool
		for i := 0; i < boardCount; i++ {
			bingo = checkBoard(boards, called, i)
			if bingo {
				winners[i] = true
				if len(winners) == boardCount {
					board = i
					break CALLS
				}

			}
		}
	}
	scoreBoard(boards, called, board, last)

	return
}
