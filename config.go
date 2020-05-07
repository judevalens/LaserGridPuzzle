package main

/**
Represents a safe config
*/

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Config struct {
	matrix       [][]*Cell
	pillarSet    []Cell
	isPillarOk   bool
	currentRow   int
	currentCol   int
	successorRow int
	successorCol int
	path         []Config
}

func NewConfig(path string) *Config {
	config := new(Config)

	config.currentCol = 0
	config.currentRow = 0
	config.successorCol = 0
	config.successorRow = 0
	config.createGrid(path)

	return config
}

func (config *Config) createGrid(path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	lineReader := bufio.NewReader(bytes.NewReader(file))

	var row, col int
	line, err := lineReader.ReadBytes('\n')
	lineSt := string(line)
	if err == io.EOF {

	}

	lineSt = strings.TrimSpace(lineSt)
	size := strings.Split(lineSt, " ")
	row, _ = strconv.Atoi(size[0])
	col, _ = strconv.Atoi(size[1])

	config.matrix = make([][]*Cell, row)

	for y, _ := range config.matrix {
		config.matrix[y] = make([]*Cell, col)
	}

	for r := 0; r < row; r++ {
		line, _ := lineReader.ReadBytes('\n')
		lineSt := string(line)
		c2 := 0
		for c := 0; c < col+col-1; c += 2 {
			config.matrix[r][c2] = NewCell(r, c2, CellType(lineSt[c:c+1]))
			c2++
		}
	}

	config.printMatrix()
}

func (config *Config) printMatrix()  {
	fmt.Printf("\n")

	for r := 0; r < len(config.matrix) + 2; r++ {
		for c := 0; c < len(config.matrix[0]) + 2; c++ {
			if r > 1 {
				if c == 0 {
					fmt.Printf("%v", (r - 2) % 10)
				} else if c == 1 {
					fmt.Printf("|")
				} else {
					fmt.Printf(string(config.matrix[r-2][c-2].element + " "))
				}
			} else {
				if r == 0 {
					if c < 2 {
						fmt.Printf(" ")
					} else {
						fmt.Printf("%v",strconv.Itoa((c - 2) % 10) + " ")
					}
				} else {
					if c < 2 {
						fmt.Printf(" ")
					} else {
						if c < (len(config.matrix[0]) + 1) {
							fmt.Printf("--")
						} else {
							fmt.Printf("-")

						}
					}
				}
			}

		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")

}

