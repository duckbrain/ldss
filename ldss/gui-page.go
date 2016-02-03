package main

import (
	"github.com/andlabs/ui"
	"ldss/lib"
)

type guiPage struct {
	app                    *gui
	item                   lib.Item
	lang                   *lib.Language
	box, toolbar, contents *ui.Box
	address                *ui.Entry
	title, status          *ui.Label
	btnUp                  *ui.Button
	contentCount           int
	childMap               map[uintptr]string
}

func newGuiPage() *guiPage {
	var err error
	p := &guiPage{}

	p.childMap = make(map[uintptr]string)

	p.box = ui.NewVerticalBox()
	p.toolbar = ui.NewHorizontalBox()
	p.contents = ui.NewVerticalBox()

	p.btnUp = ui.NewButton("Up")
	p.btnUp.OnClicked(func(btn *ui.Button) {
		p.SetItem(p.item.Parent(), true)
	})

	p.lang, err = lib.DefaultLanguage()

	if err != nil {
		panic(err)
	}

	p.title = ui.NewLabel("LDS Scriptures")
	p.status = ui.NewLabel("")
	p.address = ui.NewEntry()

	p.address.OnChanged(p.onPathChanged)

	p.toolbar.Append(p.btnUp, false)
	p.toolbar.Append(p.address, true)
	p.box.Append(p.title, false)
	p.box.Append(p.toolbar, false)
	p.box.Append(p.contents, true)
	p.box.Append(p.status, false)

	return p
}

func (p *guiPage) Lookup(s string) {
	p.handleMessages(lib.LookupPath(p.lang, s), true)
}

func (p *guiPage) SetItem(item lib.Item, setText bool) {
	for ; p.contentCount > 0; p.contentCount-- {
		p.contents.Delete(0)
	}
	p.childMap = make(map[uintptr]string)
	if item == nil {
		p.title.SetText("")
		p.btnUp.Disable()
	} else {
		if item.Parent() == nil {
			p.btnUp.Disable()
		} else {
			p.btnUp.Enable()
		}
		p.title.SetText(item.String())
		if setText {
			p.address.SetText(item.Path())
		}
		children, err := item.Children()
		if err != nil {
			p.ShowError(err)
			return
		}
		for _, c := range children {
			btn := ui.NewButton(c.Name())
			btn.OnClicked(func(btn *ui.Button) {
				path := p.childMap[btn.Handle()]
				p.Lookup(path)
			})
			p.childMap[btn.Handle()] = c.Path()
			p.contents.Append(btn, false)
			p.contentCount++
		}
	}
	p.item = item
}

func (p *guiPage) ShowError(err error) {
	p.status.Show()
	p.status.SetText(err.Error())
}

func (p *guiPage) handleMessages(c <-chan lib.Message, setText bool) {
	for m := range c {
		switch m.(type) {
		case lib.MessageDone:
			item := m.(lib.MessageDone).Item().(lib.Item)
			p.SetItem(item, setText)
			p.status.Hide()
		default:
			p.status.Show()
			p.status.SetText(m.String())
		}
	}

}

func (p *guiPage) onPathChanged(sender *ui.Entry) {
	p.handleMessages(lib.LookupPath(p.lang, sender.Text()), false)
}
