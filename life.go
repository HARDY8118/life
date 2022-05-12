package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func initData(rows, cols int) [][]bool {

	dataRows := make([][]bool, rows+2)

	for i := 0; i < rows+2; i++ {
		dataRows[i] = make([]bool, cols+2)
	}

	return dataRows
}

func drawScreen(data [][]bool, markings bool) {
	block, _ := utf8.DecodeRuneInString("▇")
	tlcorner, _ := utf8.DecodeRuneInString("┌")
	trcorner, _ := utf8.DecodeRuneInString("┐")
	brcorner, _ := utf8.DecodeRuneInString("┘")
	blcorner, _ := utf8.DecodeRuneInString("└")
	vline, _ := utf8.DecodeRuneInString("│")
	hline, _ := utf8.DecodeRuneInString("─")

	// fmt.Println("\033[2J")
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()

	if markings {
		fmt.Print("   ")
		for i := 0; i < (len(data[0]) - 2); i++ {
			fmt.Printf("%02d%c", i, vline)
		}
		fmt.Println()
		for i := 1; i < len(data)-1; i++ {
			fmt.Printf("%02d%c", i-1, vline)
			for j := 1; j < (len(data[i]) - 1); j++ {
				if data[i][j] {
					fmt.Printf("%c%c%c", block, block, vline)
				} else {
					fmt.Printf("  %c", vline)
				}
			}
			fmt.Println()
		}
		fmt.Println()
	} else {
		fmt.Printf("%c", tlcorner)
		for i := 2; i < len(data[0]); i++ {
			fmt.Printf("%c%c", hline, hline)
		}
		fmt.Printf("%c\n", trcorner)
		for r := 1; r < (len(data) - 1); r++ {
			fmt.Printf("%c", vline)
			for c := 1; c < (len(data[r]) - 1); c++ {
				if data[r][c] {
					fmt.Printf("%c%c", block, block)
				} else {
					fmt.Printf("  ")
				}
			}
			fmt.Printf("%c\n", vline)
		}
		fmt.Printf("%c", blcorner)
		for i := 2; i < len(data[0]); i++ {
			fmt.Printf("%c%c", hline, hline)
		}
		fmt.Printf("%c\n", brcorner)
	}
}

func toggleCell(data *[][]bool, row, col int) {
	if row == 0 || row >= len(*data)-1 || col == 0 || col >= len((*data)[row])-1 {
		return
	}
	(*data)[row][col] = !((*data)[row][col])
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func neighbourCount(data [][]bool, r, c int) int {
	n := 0
	if r == 0 || r == len(data)-1 || c == 0 || c == len(data[r])-1 {
		return 0
	}

	n = n + boolToInt(data[r-1][c])
	n = n + boolToInt(data[r-1][c+1])
	n = n + boolToInt(data[r][c+1])
	n = n + boolToInt(data[r+1][c+1])
	n = n + boolToInt(data[r+1][c])
	n = n + boolToInt(data[r+1][c-1])
	n = n + boolToInt(data[r][c-1])
	n = n + boolToInt(data[r-1][c-1])

	return n
}

/*
Any live cell with fewer than two live neighbours dies, as if by underpopulation.
Any live cell with two or three live neighbours lives on to the next generation.
Any live cell with more than three live neighbours dies, as if by overpopulation.
Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
*/
func tick(data *[][]bool) {

	newData := initData(len(*data)-2, len((*data)[0])-2)
	for i := 1; i < (len(*data))-1; i++ {
		for j := 1; j < (len((*data)[i]) - 1); j++ {
			neighbours := neighbourCount(*data, i, j)
			// fmt.Printf("%d %d %d\n", i, j, neighbours)
			if (*data)[i][j] {
				if neighbours > 1 && neighbours < 4 {
					// toggleCell(data, i, j)
					newData[i][j] = true
				}
			} else {
				if neighbours == 3 {
					// toggleCell(data, i, j)
					newData[i][j] = true
				}
			}
		}
	}
	(*data) = newData
}

func debugCells(data [][]bool) {
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			if data[i][j] {
				fmt.Print(" 1")
			} else {
				fmt.Print(" 0")
			}
		}
		fmt.Println()
	}
}

func input() (i string, r, c int) {
	// var inputCmd string
	// var row, col int
	fmt.Print("Toggle cells [row col] | 'start' | 'quit': ")
	inputCmd, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatal("Error: ", err)
	}

	inputArgs := strings.Split(strings.TrimSpace(inputCmd), " ")

	if len(inputArgs) > 2 || len(inputArgs) < 1 {
		log.Fatal("Invalid command. Enter row and column separated by space or a single command")
	}

	row, err := strconv.Atoi(inputArgs[0])
	if err != nil {
		return inputArgs[0], -1, -1
	}

	col, err := strconv.Atoi(inputArgs[1])
	if err != nil {
		log.Fatal("Error: ", err)
	}

	return "", row, col

}

func main() {
	const rows = 31
	const cols = 40
	const fps = 2

	data := initData(rows, cols)
	drawScreen(data, true)

	i := ""
	r := 0
	c := 0
	for len(i) == 0 {
		drawScreen(data, true)
		i, r, c = input()
		if len(i) == 0 {
			// fmt.Println("Toggle ", r, " ", c)
			toggleCell(&data, r+1, c+1)
		} else {
			if strings.Compare(i, "quit") == 0 {
				return
			} else if strings.Compare(i, "start") == 0 {
				break
			} else {
				fmt.Println("Invalid command")
			}
		}
	}

	for true {
		drawScreen(data, false)
		tick(&data)
		time.Sleep(time.Second / fps)
	}

}
