package main

/**
Represents a safe config
*/

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)


type Config struct {
	matrix       [][]Cell
	pillarSet    []Cell
	isPillarOk   bool
	currentRow   int
	currentCol   int
	successorRow int
	successorCol int
	path         []Config
}


func (c *Config)createGrid(path string){
	  file ,  err  := ioutil.ReadFile(path)

	 if err != nil {
		log.Fatal(err)
	 }
	 lineReader := bufio.NewReader(bytes.NewReader(file))

	for{
		
		line, _ := lineReader.ReadBytes('\n')

		if line == nil {
			break
		}

		fmt.Println(line)
	
	}

}
