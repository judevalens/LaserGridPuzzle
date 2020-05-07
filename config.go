package main

/**
Represents a safe config
*/
const Laser string = "L"
const Beam string = "*"
const FreeSpot string = "."
const pillars string = "01234X"

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
