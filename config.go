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
	"time"
)

type Action int

const (
	Display       Action = 1
	DisplayStatus Action = 2
	Solving       Action = 3
)

type ResponseData struct {
	statusMsg string
	action    Action
	config    *Config
}

var solutionFound bool = false
var stackCount int

type Config struct {
	matrix       [][]*Cell
	backTrackingConfig		*Config
	pillarSet    []*Cell
	isPillarOk   bool
	currentRow   int
	currentCol   int
	successorRow int
	successorCol int
	path         []Config
	isSolutionFound bool
}

func NewConfig(path string) *Config {
	config := new(Config)
	config.currentCol = 0
	config.currentRow = 0
	config.successorCol = 0
	config.successorRow = 0
	config.isPillarOk = true
	config.createGrid(path)
	//config.backTrackingConfig = config.copyConfig()
	return config
}

func (config *Config) copyConfig() *Config {
	newConfig := new(Config)

	newConfig.matrix = make([][]*Cell, len(config.matrix))

	for r, _ := range newConfig.matrix {
		newConfig.matrix[r] = make([]*Cell, len(config.matrix[0]))
	}

	newConfig.currentCol = config.successorCol
	newConfig.currentRow = config.successorRow
	newConfig.isPillarOk = config.isPillarOk
	newConfig.pillarSet = config.pillarSet

	debug(strconv.FormatBool(newConfig.isPillarOk))

	newConfig.successorCol = newConfig.currentCol + 1
	newConfig.successorRow = newConfig.currentRow
	if newConfig.successorCol >= len(newConfig.matrix[0]) {
		debug("SWITCHING ROW")
		newConfig.successorRow++
	}
	newConfig.successorCol = newConfig.successorCol % len(newConfig.matrix[0])

	debug("Current row " + strconv.Itoa(newConfig.currentRow))
	debug("Current col " + strconv.Itoa(newConfig.currentCol))
	debug("Succ row " + strconv.Itoa(newConfig.successorRow))
	debug("Succ col " + strconv.Itoa(newConfig.successorCol))
	for r, rv := range newConfig.matrix {
		for c, rc := range rv {
			newConfig.matrix[r][c] = config.matrix[r][c].CopyCell()
			_ = rc
		}
	}
	return newConfig
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

	config.pillarSet = make([]*Cell, 0)

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
			if strings.Contains(string(Pillars),string(config.matrix[r][c2].element)){
				if string(config.matrix[r][c2].element) != "X"{
					config.pillarSet = append(config.pillarSet,config.matrix[r][c2])

				}
			}
			c2++
		}
	}

	config.printMatrix()
}

/**
 * adds a laser at the given coordinates
 *
 * @param row
 * @param col
 */
func (config *Config) addLaser(row int, col int) {

	var response string
	var status bool

	if row < 0 || row > len(config.matrix) || col < 0 || col > len(config.matrix[0]) {
		status = false
	} else {
		config.matrix[row][col].failedVerification = false

		status = config.matrix[row][col].updateElement(Laser)
	}

	if status {
		response = "Laser added at: (" + strconv.Itoa(row) + ", " + strconv.Itoa(col) + ")"
		config.beamRow(row, col, 0, 1, config.matrix[row][col])
		config.beamCol(row, col, 0, 1, config.matrix[row][col])
		//lastConfigStack.add(safeMatrixC[row][col])
	} else {
		response = "Error adding laser at: (" + strconv.Itoa(row) + ", " + strconv.Itoa(col) + ")"
	}

	data := ResponseData{
		statusMsg: response,
		action:    DisplayStatus,
		config:    config,
	}

	update(data)
}

/**
 * adds a laser at the given coordinates
 * same as addLaser but its suited to use in the backtracker
 * @param row
 * @param col
 */
func (config *Config) addLaserB(row int, col int) bool {

	var isCorrect bool

	if row < 0 || row > len(config.matrix) || col < 0 || col > len(config.matrix[0]) {
		isCorrect = false
	} else {
		isCorrect = config.matrix[row][col].updateElement(Laser)

		if isCorrect {
			config.beamRow(row, col, 0, 1, config.matrix[row][col])
			config.beamCol(row, col, 0, 1, config.matrix[row][col])
			debug("IS CORRECT")

		}


	}
	return isCorrect
}

/**
 * removes laser at the given coordinates
 *
 * @param row
 * @param col
 */
func (config *Config) removeLaser(row int, col int) {

	var response string
	var status bool

	if row < 0 || row > len(config.matrix) || col < 0 || col > len(config.matrix[0]) {
		status = false
	} else {
		config.matrix[row][col].failedVerification = false

		status = config.matrix[row][col].updateElement(FreeSpot)
	}

	if status {
		response = "Laser removed at: (" + strconv.Itoa(row) + ", " + strconv.Itoa(col) + ")"
		config.beamRow(row, col, 0, -1, config.matrix[row][col])
		config.beamCol(row, col, 0, -1, config.matrix[row][col])
		//lastConfigStack.add(safeMatrixC[row][col])
	} else {
		response = "Error removing laser at: (" + strconv.Itoa(row) + ", " + strconv.Itoa(col) + ")"
	}

	data := ResponseData{
		statusMsg: response,
		action:    DisplayStatus,
		config:    config,
	}

	update(data)
}

func (config *Config) verify() {

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
	if col >= 0 && col < len(config.matrix[0]) {
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

func (config *Config) getSuccessor() []*Config {

	var successors []*Config

	if config.successorRow < len(config.matrix) {

		for _, pillar := range config.pillarSet {
			cell := config.matrix[pillar.row][pillar.col]

			if (config.successorRow - cell.row) > 1 {
				if cell.adjacentPillar != cell.pillarNumber {
					//fmt.Printf("INVALID  cell.adjacentPillar %v  cell.pillarNumber %v\n", cell.adjacentPillar,cell.pillarNumber)
					config.isPillarOk = false
					break
				}
			}
			config.isPillarOk = true
		}

			configA := config.copyConfig()
			configB := config.copyConfig()
			result := configB.addLaserB(configB.currentRow, configB.currentCol)
			debug("RESULT  " + strconv.FormatBool(result))

			if result {
				debug("aDDED LASER")
				successors = append(successors,configB)
			}

		successors = append(successors,configA)

	}else{
		debug("FAILED 2")
	}

	return successors

}


/**
verifies that this configuration is valid at its current state
 */
func (config *Config) isValid() bool {


	isValid := config.isPillarOk
	if isValid {
		 cell := config.matrix[config.currentRow][config.currentCol]
		 element := cell.element

		if element == Laser {
			debug("TOO MANY LASERS at  :" + strconv.Itoa(cell.row) + " " + strconv.Itoa(cell.col) + " # " + strconv.Itoa(cell.adjacentLaser))
			if cell.row == 8 && cell.col == 6 {

				//os.Exit(0)
			}
			if cell.adjacentLaser > 0 {
				isValid = false

			}


		}
	}
	return isValid
}

/**
verifies that the current configuration is a solution
*/
func (config Config) isGoal() bool {

	isGoal := false
	r := 0
	c := 0
	for r < len(config.matrix) {
		c = 0
		for c < len(config.matrix[0]) {
			cell := config.matrix[r][c]
			element := cell.element
			if element == (Laser) {
				if cell.adjacentLaser > 0 {
					isGoal = true
					fmt.Printf("TOO MANY LASERS at %v %v", cell.row, cell.col)

					break
				}
			} else if strings.Contains(string(Pillars), string(element)) {
				if string(element) != "X" {
					if cell.adjacentPillar != cell.pillarNumber {
						debug(string(element) + " SUPPOSED TO BE " +  strconv.Itoa(cell.pillarNumber) + " " +  strconv.Itoa(cell.adjacentPillar))
						isGoal = true
						break
					}
				}
			} else if element != Beam {
				debug("CANT BE A DOT at " + strconv.Itoa(r) + " " +strconv.Itoa(c))
				isGoal = true
				break
			}

			c++
		}
		if isGoal {
			break
		}

		r++
	}
	return !isGoal
}

func (config *Config) solve() *Config{
	debug("current config")
	//config.printMatrix()

	if config.isGoal() {
		solutionFound = true
		debug("FOUND")
		config.printMatrix()
		return  config
	}else{
		for _,c := range config.getSuccessor(){

			if c.isValid() {
				debug("valid config")
				///ssc.printMatrix()
				stackCount++
				sol := c.solve()
				stackCount--

				if solutionFound {
					if stackCount == 0 {
						solutionFound = false
					}
					return sol
				}
			}else{
				debug("invalid config")
				//c.printMatrix()
			}
		}
	}
	return nil
}

func (config *Config) getSolution(){
	debug("SOLVING " + strconv.Itoa(config.currentCol))
	start := time.Now()
	s := config.solve()
	stop := time.Now()

	fmt.Printf("DURATION %v", stop.Sub(start))
	s.printMatrix()
	debug("SOL FOUND")
}
