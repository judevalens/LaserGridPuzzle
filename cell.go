package main

import (
	"container/list"
	"math"
	"strings"
	"strconv"
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
	row                int
	col                int
	element            CellType
	adjacentLaser      int
	adjacentPillar     int
	failedVerification bool
	laserDependency    []Cell
	laserDependencyList list.List
	pillarNumber       int
}

func NewCell(row int, col int, element CellType) *Cell {
	var cell *Cell = new(Cell)
	cell.row  = row
	cell.col =  col
	cell.element = element
	if strings.Contains(string(Pillars),string(element)){
		if string(element) != "X" {
			n, _ := strconv.ParseInt(string(Pillars),10,16)
			cell.pillarNumber = int(n)
		}else {
			cell.pillarNumber = -1
		}
	}
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
			if len(c.laserDependency) == 0 {
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

	if action > 0 {

		if c.element == FreeSpot || c.element == Beam {
			c.element = e
			c.laserDependency = append(c.laserDependency, laser)
			r = true
		} else if c.element == Laser {
			laser.adjacentLaser = laser.adjacentLaser + 1
			c.laserDependency =append(c.laserDependency,laser)

			debug("ADJ LASER at " + string(c.row) + " " + string(c.col) + " LASER: " + string(laser.row) + " " + string(laser.
				col) + " #" + string(laser.adjacentLaser))
			c.adjacentLaser++
			r = true
		} else if strings.Contains(string(Pillars),string(e)) {

			lDistance := math.Abs(float64(c.row - laser.row))
			hDistance := math.Abs(float64(c.col - laser.col))

			if lDistance == 1 || hDistance == 1 {
				c.adjacentPillar++
			}

		}
	} else {
		if c.element == Beam {

			debug("SIZE OF LIST BEFORE REMOVING Dependency " + string(len(c.laserDependency)))
			c.laserDependency = c.laserDependency[:len(c.laserDependency)-1]
			debug("SIZE OF LIST AFTER REMOVING Dependency " + string(len(c.laserDependency)))

			if len(c.laserDependency) == 0 {
				c.element = FreeSpot
			}
			r = true
		} else if c.element == Laser {
			c.laserDependency = c.laserDependency[:len(c.laserDependency)-1]
			laser.adjacentLaser--
			c.adjacentLaser--
			r = true
		} else if strings.Contains(string(Pillars),string(c.element)) {

			lDistance := math.Abs(float64(c.row - laser.row))
			hDistance := math.Abs(float64(c.col - laser.col))

			if lDistance == 1 || hDistance == 1 {

				c.adjacentPillar--
			}
		}
	}
	return r
}
