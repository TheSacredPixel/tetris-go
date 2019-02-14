package main

import (
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

const coldef = termbox.ColorDefault

var tetrQueue = make(chan *Tetrimino)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.HideCursor()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	b := NewBoard()
	t := GetTetr(b)

	for {
		select { //for the different channels
		case event := <-eventQueue:
			switch event.Type {
			case termbox.EventKey:
				if event.Ch == 'z' || event.Ch == 'x' {
					t.Rotate(event.Ch)
				}
				switch event.Key { //inputs
				case termbox.KeyEsc:
					return

				case termbox.KeyArrowRight:
					t.Move('r')
				case termbox.KeyArrowLeft:
					t.Move('l')
				case termbox.KeyArrowUp:
					if t.hard != 1 {
						t.HardDrop()
					}

				case termbox.KeyEnter:
					t.Stop()

				}

			case termbox.EventResize:

			}
		case t = <-tetrQueue:
		}
		WriteDebug(t)
	}
}
