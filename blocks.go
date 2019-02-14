package main

import (
	"github.com/nsf/termbox-go"
	"math"
	"math/rand"
	"time"
)

var (
	tetrTypes = []rune{
		'i',
		'o',
		't',
		's',
		'z',
		'j',
		'l',
	}

	blockMap = map[rune][]int{
		'i': []int{
			0x0f00,
			0x2222,
			0x00f0,
			0x4444,
		},
		'o': []int{
			0x6600,
		},
		't': []int{
			0x4e00,
			0x4640,
			0x0e40,
			0x4c40,
		},
		's': []int{
			0x6c00,
			0x4620,
			0x06c0,
			0x8c40,
		},
		'z': []int{
			0xc600,
			0x2640,
			0x0c60,
			0x4c80,
		},
		'j': []int{
			0x8e00,
			0x6440,
			0x0e20,
			0x44c0,
		},
		'l': []int{
			0x2e00,
			0x4460,
			0x0e80,
			0xc440,
		},
	}

	colorMap = map[rune]termbox.Attribute{
		'i': termbox.ColorCyan,
		'o': termbox.ColorYellow,
		't': termbox.ColorMagenta,
		's': termbox.ColorGreen,
		'z': termbox.ColorRed,
		'j': termbox.ColorBlue,
		'l': termbox.ColorWhite, //no orange :(
	}
)

//GetTetr : get new random tetrimino
func GetTetr(b *Board) *Tetrimino {
	pick := tetrTypes[rand.Intn(len(tetrTypes))]

	t := new(Tetrimino)
	t.name = pick
	t.color = colorMap[pick]
	t.blocks = blockMap[pick]
	t.facing = 0

	t.board = b

	t.locX = 3
	t.locY = -2

	t.hard = 0

	t.Draw(t.color)
	t.StartFall()
	return t
}

//Tetrimino : basic struct
type Tetrimino struct {
	name   rune
	color  termbox.Attribute
	blocks []int
	facing int
	locX   int
	locY   int
	kill   chan struct{}
	hard	 int
	board  *Board
}

//Draw : draw tetrimino on screen
//Setting color to 0 erases it
func (t *Tetrimino) Draw(color termbox.Attribute) {
	x, y := 0, 0

	for bit := 0x8000; bit > 0; bit = bit >> 1 {
		if (bit&t.blocks[t.facing]) != 0 && t.locY+y >= 0 {
			termbox.SetCell(t.board.offsetX+t.locX*2+x, t.board.offsetY+t.locY+y, ' ', coldef, color)
			termbox.SetCell(t.board.offsetX+t.locX*2+x+1, t.board.offsetY+t.locY+y, ' ', coldef, color)
		}
		if x < 6 {
			x += 2
		} else {
			x = 0
			y++
		}
	}
	termbox.Flush()
}

//Move : move tetrimino
func (t *Tetrimino) Move(dir rune) {
	t.Draw(0)
	switch dir {
	case 'r':
		if t.locX++; t.CheckCollision() {
			t.locX--
		}
	case 'l':
		if t.locX--; t.CheckCollision() {
			t.locX++
		}
	case 'd':
		if t.locY++; t.CheckCollision() {
			t.locY--
			t.Stop()
		}
	}
	t.Draw(t.color)
}

//Rotate : rotate tetrimino
func (t *Tetrimino) Rotate(dir rune) {
	t.Draw(0) //clear from board
	face := t.facing

	switch dir {
	case 'x':
		if t.facing++; t.facing == len(t.blocks) { //overflow facing
			t.facing = 0
		}
		if t.CheckCollision() {
			t.facing = face
		}
	case 'z':
		if t.facing--; t.facing < 0 {
			t.facing = len(t.blocks) - 1
		}
		if t.CheckCollision() {
			t.facing = face
		}
	}
	t.Draw(t.color)
}

//CheckCollision :check for collision with walls/placed blocks
//								returns true if collision occurs
func (t *Tetrimino) CheckCollision() bool {
	x, y := 0, 0

	for bit := 0x8000; bit > 0; bit = bit >> 1 {
		if (bit & t.blocks[t.facing]) != 0 {
			if t.locX+x == -1 || t.locX+x == 10 || t.locY+y == 20 { //check walls
				return true
			}
			if t.locY >= 0 && t.board.blocks[t.locY+y][t.locX+x] != 0 {
				return true
			}
		}
		if x < 3 {
			x++
		} else {
			x = 0
			y++
		}
	}
	return false
}

//StartFall : initiate the fall timer
//TODO: add soft dropping (20x)
func (t *Tetrimino) StartFall() {
	stop := make(chan struct{})
	level := float64(t.board.level)
	t.kill = stop

	go func() {
	loop:
		for {
			select {
			case <-stop:
				break loop
			default:
				t.Move('d')
				time.Sleep(time.Duration(math.Pow(0.8-((level-1)*0.007), (level-1))) * time.Second)
				//time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

//HardDrop : initiate hard drop
func (t *Tetrimino) HardDrop() {
	close(t.kill)
	t.hard = 1

	go func() {
		for !t.CheckCollision() {
			t.Move('d')
			time.Sleep(100 * time.Microsecond)
		}
		t.Stop()
	}()
}

//Stop : kill tetrimino and place it
func (t *Tetrimino) Stop() {
	if t.hard == 0 {
		close(t.kill)
	}
	x, y := 0, 0
	for bit := 0x8000; bit > 0; bit = bit >> 1 { //place blocks on board
		if (bit & t.blocks[t.facing]) != 0 {
			t.board.blocks[t.locY+y][t.locX+x] = t.color

		}
		if x < 3 {
			x++
		} else {
			x = 0
			y++
		}
	}

	//check for game over conditions
	//TODO

	//otherwise...
	t.board.Refresh()
	tetrQueue <- GetTetr(t.board)
}
