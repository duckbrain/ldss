// +build nogui

/*
 * This file will contain a native GUI that can be run.
 */
package main

import (
    "github.com/andlabs/ui"
)

type gui struct {
	appinfo
}

func (app gui) run () {
    err := ui.Main(func() {
        box := ui.NewVerticalBox()
	toolbar := ui.NewHorizontalBox()
	contents := ui.NewVerticalBox()

	title := ui.NewLabel("LDS Scriptures")
	address := ui.NewEntry()

	address.OnChanged(func(sender *ui.Entry) {
		_, _ := app.config.Library.Lookup(sender.Text())
	})

	toolbar.Append(address, false)
	box.Append(title, false)
	box.Append(toolbar, false)
	box.Append(contents, false)

        window := ui.NewWindow("Hello", 200, 100, false)
        window.SetChild(box)
        window.OnClosing(func(*ui.Window) bool {
            ui.Quit()
            return true
        })
        window.Show()
    })
    if err != nil {
        panic(err)
    }
}
