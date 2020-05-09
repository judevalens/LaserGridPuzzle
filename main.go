package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const DEBUG bool = false
var safeConfig *Config
func main() {
	safeConfig = NewConfig(os.Args[1])

	fmt.Printf("The cell %v is at %v\n", safeConfig.currentRow, safeConfig.currentCol)

	readInput()

}


func readInput() {
	scanner := bufio.NewScanner(os.Stdin)

	text := ""
	fmt.Println("")

		fmt.Printf("> ")

		scanner.Scan()
		text = scanner.Text()
		commands := strings.Split(text, " ")
		exec(commands)
		fmt.Printf("%s\n------\n", text)

}

func exec(commands []string) {
	cmtList := "arhqvsd"
	isCommandValid := true

	command := commands[0]
	command = strings.ToLower(command)

	if !strings.Contains(cmtList, command) {
		isCommandValid = false
	} else if command == "a" || command == "r" {
		if len(commands) != 3 {
			isCommandValid = false
		}
	}

	if isCommandValid {

		switch command {
		case "a":
			r, _ := strconv.Atoi(commands[1])
			c, _ := strconv.Atoi(commands[2])
			safeConfig.addLaser(r,c)
		case "r":
			r, _ := strconv.Atoi(commands[1])
			c, _ := strconv.Atoi(commands[2])
			safeConfig.removeLaser(r,c)
		case "c":
			safeConfig.verify()
		case "d":
			safeConfig.printMatrix()
			readInput()
		case "s":
			safeConfig.getSolution()
		case "v":

			debug("VERIFIED ? " + strconv.FormatBool(safeConfig.isGoal()))

		case "q":
			os.Exit(0)
		}
	}

}

func update(response ResponseData)  {
	if response.action == DisplayStatus{
		fmt.Printf("%s",response.statusMsg)
		response.config.printMatrix()
	}else if response.action == Display {
		safeConfig.printMatrix()
	}

	readInput()
}

func debug(s string)  {
	if DEBUG {
		fmt.Println(s)
	}
}
