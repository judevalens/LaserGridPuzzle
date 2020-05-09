package main

import (
	"container/list"
	"math"
	"strconv"
	"strings"
)

/*
Represent a cell
*/

type CellType string

const (
	Laser    CellType = "L"
	Beam     CellType = "*"
	FreeSpot CellType = "."
	Pillars  CellType = "01234X"
)

type Cell struct {
	row                 int
	col                 int
	element             CellType
	adjacentLaser       int
	adjacentPillar      int
	failedVerification  bool
	laserDependency     int
	laserDependencyList list.List
	pillarNumber        int
}

func NewCell(row int, col int, element CellType) *Cell {
	var cell *Cell = new(Cell)
	cell.row = row
	cell.col = col
	cell.element = element
	if strings.Contains(string(Pillars), string(element)) {
		if string(element) != "X" {
			n, _ := strconv.ParseInt(string(cell.element), 10,64)
			cell.pillarNumber = int(n)
			debug(string(cell.element) + " PILLAR N " + strconv.Itoa(cell.pillarNumber))
		} else {
			cell.pillarNumber = -1
		}
	}
	return cell
}

func (c *Cell) CopyCell() *Cell{
	var cell  = new(Cell)

	cell.row = c.row
	cell.col = c.col
	cell.element = c.element
	cell.pillarNumber = c.pillarNumber
	cell.adjacentLaser = c.adjacentLaser
	cell.adjacentPillar = c.adjacentPillar
	cell.laserDependency = c.laserDependency

	return cell
}

func (c *Cell) AdjacentLaser() int {
	return c.adjacentLaser
}

func (c *Cell) SetAdjacentLaser(adjacentLaser int) {
	c.adjacentLaser = adjacentLaser
}

func (c *Cell) updateElement(e CellType) bool {
	r := false
	if e == Laser {
		if c.element == FreeSpot || c.element == Beam {
			c.element = e
			r = true
		}
	} else if e == FreeSpot {
		if c.element == Laser {
			if c.laserDependency == 0 {
				c.element = e
			} else {
				c.element = Beam
			}

			r = true
		}
	}

	return r
}

func (c *Cell) propagate(e CellType, laser *Cell, action int) bool {

	r := false
	// when its a free spot or beam/ we add the new laser dependency and change the
	// symbol to a bea
	// when its a laser we just add the dependency and check the adjacent cells
	// only a pillar can stop a laser's trajectory
	//fmt.Printf("CELL TYPE %v\n" , string(e))
	if action > 0 {

		if c.element == FreeSpot || c.element == Beam {
			c.element = e
			c.laserDependency++
			r = true
		} else if c.element == Laser {
			laser.adjacentLaser++
			c.laserDependency++
			debug("ADJ LASER at " +  strconv.Itoa(c.row) + " " +  strconv.Itoa(c.col) + " LASER: " +  strconv.Itoa(laser.row) + " " +  strconv.Itoa(laser.
				col) + " #" +  strconv.Itoa(laser.adjacentLaser))
			c.adjacentLaser++
			r = true
		} else if strings.Contains(string(Pillars), string(c.element)) {

			lDistance := math.Abs(float64(c.row - laser.row))
			hDistance := math.Abs(float64(c.col - laser.col))
			///fmt.Printf("DISTANCE %v %v", lDistance,hDistance)
			if lDistance == 1 || hDistance == 1 {
				c.adjacentPillar++
			}

		}
	} else {
		if c.element == Beam {

			debug("SIZE OF LIST BEFORE REMOVING Dependency " + string(c.laserDependency))
			c.laserDependency--
			debug("SIZE OF LIST AFTER REMOVING Dependency " + string(c.laserDependency))

			if c.laserDependency == 0 {
				c.element = FreeSpot
			}
			r = true
		} else if c.element == Laser {
			c.laserDependency--
			laser.adjacentLaser--
			c.adjacentLaser--
			r = true
		} else if strings.Contains(string(Pillars), string(c.element)) {

			lDistance := math.Abs(float64(c.row - laser.row))
			hDistance := math.Abs(float64(c.col - laser.col))

			if lDistance == 1 || hDistance == 1 {

				c.adjacentPillar--
			}
		}
	}
	return r
}
