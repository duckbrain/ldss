package main

import (
	"github.com/rthornton128/goncurses"
	"log"
	"ldslib"
)

type curses struct {
	args []string
	config Config
}

func (app curses) run() {
	stdscr, err := goncurses.Init()
	if err != nil {
		log.Fatal("init:", err)
	}
	defer goncurses.End()
	
	goncurses.Raw(true)
	goncurses.Echo(false)
	goncurses.Cursor(0)
	stdscr.Clear()
	stdscr.Keypad(true)
	
}

func (app *curses) cursesDisplay(item ldslib.CatalogItem) {
	switch item.(type) {
		case ldslib.Node:
			//fmt.Println(config.Library.Content(item.(ldslib.Node)))
		default:
			children, err := app.config.Library.Children(item)
			if err != nil {
				panic (err)
			}
			_ = children
	}
}