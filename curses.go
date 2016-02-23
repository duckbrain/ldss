// +build !nocurses,broken

package main

import (
	"ldss/lib"
	"log"

	"github.com/rthornton128/goncurses"
)

type curses struct {
	appinfo
	catalog *lib.Catalog
}

func (app *curses) run() {
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

func (app *curses) cursesDisplay(item lib.Item) {
	switch item.(type) {
	case lib.Node:
		//fmt.Println(config.Library.Content(item))
	default:
		children, err := app.config.Library.Children(item)
		if err != nil {
			panic(err)
		}
		_ = children
	}
}
