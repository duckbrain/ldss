package main

import (
	"github.com/rthornton128/goncurses"
	"ldss/lib"
	"log"
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

	app.catalog = app.config.SelectedCatalog()
}

func (app *curses) cursesDisplay(item lib.Item) {
	switch item.(type) {
	case lib.Node:
		//fmt.Println(config.Library.Content(item.(lib.Node)))
	default:
		children, err := app.config.Library.Children(item)
		if err != nil {
			panic(err)
		}
		_ = children
	}
}
