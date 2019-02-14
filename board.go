package main

import (
	"github.com/nsf/termbox-go"
	"strconv"
)

//BoardState : the current Board state
type BoardState int

const (
	//StateBoot 0
	StateBoot BoardState = iota
)

//NewBoard : initialize Board state
func NewBoard() *Board {
	b := new(Board)

	b.offsetX = 6
	b.offsetY = 2
	b.level = 0

	b.DrawBoard()

	return b
}

//Board : the Board state
type Board struct {
	offsetX int
	offsetY int
	blocks  [20][10]termbox.Attribute //[y][x]
	state   BoardState
	level   int
}

//DrawBoard : draw the board outline
func (b *Board) DrawBoard() {
	drawLine(b.offsetX-1, b.offsetY-1, 'x', 22)
	drawLine(b.offsetX-1, b.offsetY+20, 'x', 22)
	drawLine(b.offsetX-2, b.offsetY-1, 'y', 22)
	drawLine(b.offsetX-1, b.offsetY-1, 'y', 22)
	drawLine(b.offsetX+20, b.offsetY-1, 'y', 22)
	drawLine(b.offsetX+21, b.offsetY-1, 'y', 22)
}

func drawLine(x int, y int, iterate rune, times int) {
	if iterate == 'x' {
		for i := 0; i < times; i++ {
			termbox.SetCell(x+i, y, ' ', coldef, termbox.ColorWhite)
		}
	} else {
		for i := 0; i < times; i++ {
			termbox.SetCell(x, y+i, ' ', coldef, termbox.ColorWhite)
		}
	}
}

//Refresh :
func (b *Board) Refresh() {
	termbox.Clear(coldef, coldef)

	for i, y := range b.blocks {
		for j, x := range y {
			termbox.SetCell(b.offsetX+j*2, b.offsetY+i, ' ', coldef, x)
			termbox.SetCell(b.offsetX+j*2+1, b.offsetY+i, ' ', coldef, x)
		}
	}

	b.DrawBoard()
	termbox.Flush()
}

//WriteDebug : duh
func WriteDebug(t *Tetrimino) {
	line := "Type: " + string(t.name) + ", facing: " + strconv.FormatInt(int64(t.facing), 10) + ", blocks: " + strconv.FormatInt(int64(t.blocks[t.facing]), 16) + ", loc: " + strconv.FormatInt(int64(t.locX), 10) + "," + strconv.FormatInt(int64(t.locY), 10) + "   "
	writeLine(line, []int{0, 0})
	termbox.Flush()
}

//WriteLine : draw contents of a string
func writeLine(line string, loc []int) {
	for i, char := range line {
		termbox.SetCell(loc[0]+i, loc[1], char, coldef, coldef)
	}
}
