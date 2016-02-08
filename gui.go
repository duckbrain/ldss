// +build !nogui

/*
 * This file will contain a native GUI that can be run.
 */
package main

import (
	"fmt"
	"github.com/andlabs/ui"
)

type gui struct {
	appinfo
	pages []guiPage

	// Controls
	tab    *ui.Tab
	window *ui.Window
}

func init() {
	app := gui{}
	app.pages = make([]guiPage, 0)
	apps["gui"] = &app
}

func (app gui) run() {
	err := ui.Main(func() {
		app.tab = ui.NewTab()

		app.addPage("/")

		app.window = ui.NewWindow("LDS Scriptures", 200, 300, false)
		app.window.SetChild(app.tab)
		app.window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		app.window.Show()
	})
	if err != nil {
		panic(err)
	}
}

func (app *gui) addPage(path string) {
	page := newGuiPage()
	app.tab.Append(fmt.Sprintf("Tab %v", app.tab.NumPages()+1), page.box)
	page.btnNewTab.OnClicked(func(btn *ui.Button) {
		app.addPage("/")
	})
	page.btnCloseTab.OnClicked(func(btn *ui.Button) {
		//TODO: Page tracks index to remove it
	})
	page.Lookup(path)
}
