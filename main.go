package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const DEBUG bool = true

func main() {
	var c *Config = NewConfig(os.Args[1])

	fmt.Printf("The cell %v is at %v\n", c.currentRow, c.currentCol)

	readInput()

}



func readInput() {
	scanner := bufio.NewScanner(os.Stdin)

	text := ""
	fmt.Println("")

	for text != "stop" {
		fmt.Printf("> ")

		scanner.Scan()
		text = scanner.Text()
		commands := strings.Split(text, " ")
		exec(commands)
		fmt.Printf("%s\n------\n", text)

	}
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

	}

}

func update(response ResponseData)  {
	if response.action == DisplayStatus{
		fmt.Printf("%s",response.statusMsg)
		response.config.printMatrix()
	}else if response.action == Display {
		response.config.printMatrix()
	}

	readInput()
}

func debug(s string)  {
	if DEBUG {
		fmt.Println(s)
	}
}
