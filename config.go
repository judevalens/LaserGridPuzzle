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

func (config *Config) printMatrix() {
	fmt.Printf("\n")

	for r := 0; r < len(config.matrix)+2; r++ {
		for c := 0; c < len(config.matrix[0])+2; c++ {
			if r > 1 {
				if c == 0 {
					fmt.Printf("%v", (r-2)%10)
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
						fmt.Printf("%v", strconv.Itoa((c-2)%10)+" ")
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

/**
 * adds lazer at the given coordinates
 *
 * @param row
 * @param col
 */
func (config *Config) addLaser(row int, col int) {

	var response string
	var status bool
	config.matrix[row][col].failedVerification = true

	if row < 0 || row > len(config.matrix) || col < 0 || col > len(config.matrix[0]) {
		status = false
	} else {
		status = config.matrix[row][col].updateElement(Laser)
	}

	if status {
		response = "Laser added at: (" + string(row) + ", " + string(col) + ")"
		config.beamRow(row, col, 0, 1, config.matrix[row][col])
		config.beamCol(row, col, 0, 1, config.matrix[row][col])
		//lastConfigStack.add(safeMatrixC[row][col])
	} else {
		response = "Error adding laser at: (" + string(row) + ", " + string(col) + ")"
	}

	data := ResponseData{
		statusMsg: response,
		action: DisplayStatus,
	}

	update(config, data)
}

/**
 * adds laser at the given coordinates
 *
 * @param row
 * @param col
 */
func (config *Config) removeLaser(row int, col int) {

	var response string
	var status bool
	config.matrix[row][col].failedVerification = true

	if row < 0 || row > len(config.matrix) || col < 0 || col > len(config.matrix[0]) {
		status = false
	} else {
		status = config.matrix[row][col].updateElement(FreeSpot)
	}

	if status {
		response = "Laser removed at: (" + string(row) + ", " + string(col) + ")"
		config.beamRow(row, col, 0, -1, config.matrix[row][col])
		config.beamCol(row, col, 0, -1, config.matrix[row][col])
		//lastConfigStack.add(safeMatrixC[row][col])
	} else {
		response = "Error removing laser at: (" + string(row) + ", " + string(col) + ")"
	}

	data := ResponseData{
		statusMsg: response,
		action: DisplayStatus,
		config: config,
	}

	update(data)
}

/**
 * when a laser is added this method adds beam on the cells that Vertically
 * adjacent to the laser It stops when a pillar is encountered
 *
 * @param row
 * @param col
 * @param direction up or down (1 or -1)
 * @param action    add or remove beams (1 or -1)
 * @param laser     the laser that emits those beams
 */
func (config *Config) beamRow(row int, col int, direction int, action int, laser *Cell) {
	if row >= 0 && row < len(config.matrix) {
		var status bool
		if direction == 0 {
			config.beamRow(row+1, col, 1, action, laser)
			config.beamRow(row-1, col, -1, action, laser)
		} else {
			status = config.matrix[row][col].propagate(Beam, laser, action)

			if status {
				config.beamRow(row+direction, col, direction, action, laser)
			}

		}

	}
}

/**
 * when a laser is added this method adds beam on the cells that horizontally
 * adjacent to the laser It stops when a pillar is encountered
 *
 * @param row       a row
 * @param col       a col
 * @param direction up or down (1 or -1)
 * @param action    add or remove beams (1 or -1)
 * @param laser     the laser that emits those beams
 */
func (config *Config) beamCol(row int, col int, direction int, action int, laser *Cell) {
	if row >= 0 && row < len(config.matrix) {
		var status bool
		if direction == 0 {
			config.beamCol(row, col+1, 1, action, laser)
			config.beamCol(row, col-1, -1, action, laser)
		} else {
			status = config.matrix[row][col].propagate(Beam, laser, action)

			if status {
				config.beamCol(row, col+direction, direction, action, laser)
			}

		}

	}
}


type Action int

const (
	Display       Action = 1
	DisplayStatus        = 2
	Solving
)

type ResponseData struct {
	statusMsg string
	action Action
	config *Config

}
