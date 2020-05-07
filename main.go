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

	for text != "stop" {
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

func debug(s string)  {
	if DEBUG {
		fmt.Println(s)
	}
}
