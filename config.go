package main

/**
Represents a safe config
*/


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
