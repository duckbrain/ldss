// +build !gui

package main

import (
	"fmt"
	"ldss/lib"

	"github.com/andlabs/ui"
)

type gui struct {
	appinfo
	pages     []*guiPage
	languages []*lib.Language

	// Controls
	tab    *ui.Tab
	window *ui.Window
}

func init() {
	app := gui{}
	app.pages = make([]*guiPage, 0)
	apps["gui"] = &app
}

func (app *gui) run() {
	if langs, err := lib.Languages(); err == nil {
		app.languages = langs
	} else {
		panic(err)
	}

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
	page := newGuiPage(app)
	app.pages = append(app.pages, page)
	app.tab.Append(fmt.Sprintf("Tab %v", app.tab.NumPages()+1), page.box)
	page.btnNewTab.OnClicked(func(btn *ui.Button) {
		app.addPage("/")
	})
	page.btnCloseTab.OnClicked(func(btn *ui.Button) {
		//TODO: Page tracks index to remove it
		for i, p := range app.pages {
			if p.btnCloseTab == btn {
				app.removePage(i)
			}
		}
	})
	/*page.contents.onItemChange = func(item lib.Item, r *guiRenderer) {
		for i, p := range app.pages {
			if p.contents == r {
				//app.tab.Delete(i)
				//app.tab.InsertAt(item.Name(), i, app.pages[i].box)
			}
		}
	}*/
	page.Lookup(path)
}

func (app *gui) removePage(i int) {
	app.pages = append(app.pages[:i], app.pages[i+1:]...)
	app.tab.Delete(i)
}
