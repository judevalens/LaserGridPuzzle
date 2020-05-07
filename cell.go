package main

/*
Represent a cell
 */
type Cell struct {
	row                int
	col                int
	Element            string
	adjacentLaser      int
	adjacentPillar     int
	failedVerification bool
	laserDependency    []Cell
}

func NewCell(row int, col int, element string) *Cell {
	return &Cell{row: row, col: col, Element: element}
}

func (c *Cell) AdjacentLaser() int {
	return c.adjacentLaser
}

func (c *Cell) SetAdjacentLaser(adjacentLaser int) {
	c.adjacentLaser = adjacentLaser
}
