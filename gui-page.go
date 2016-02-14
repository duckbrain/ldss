// +build !nogui

package main

import (
	"fmt"
	"ldss/lib"

	"github.com/andlabs/ui"
	"github.com/duckbrain/uidoc"
)

type guiPage struct {
	app                                                                     *gui
	item                                                                    lib.Item
	lang                                                                    *lib.Language
	box, toolbar                                                            *ui.Box
	contents                                                                *uidoc.UIDoc
	address                                                                 *ui.Entry
	title, status                                                           *ui.Label
	btnUp, btnNext, btnPrevious, btnNewTab, btnCloseTab                     *ui.Button
	titleFont, subtitleFont, summaryFont, verseFont, contentFont, errorFont *ui.Font
}

func newGuiPage() *guiPage {
	var err error
	p := &guiPage{}

	//p.childMap = make(map[uintptr]string)

	p.box = ui.NewVerticalBox()
	p.toolbar = ui.NewHorizontalBox()
	p.contents = uidoc.New()

	p.btnUp = ui.NewButton("")
	p.btnUp.OnClicked(func(btn *ui.Button) {
		p.SetItem(p.item.Parent(), true)
	})
	p.btnNext = ui.NewButton("")
	p.btnNext.OnClicked(func(btn *ui.Button) {
		p.SetItem(p.item.Next(), true)
	})
	p.btnPrevious = ui.NewButton("")
	p.btnPrevious.OnClicked(func(btn *ui.Button) {
		p.SetItem(p.item.Previous(), true)
	})
	p.btnNewTab = ui.NewButton("")
	p.btnCloseTab = ui.NewButton("")

	p.lang, err = lib.DefaultLanguage()

	if err != nil {
		panic(err)
	}

	p.title = ui.NewLabel("LDS Scriptures")
	p.status = ui.NewLabel("")
	p.address = ui.NewEntry()

	p.address.OnChanged(p.onPathChanged)

	p.toolbar.Append(p.btnPrevious, false)
	p.toolbar.Append(p.btnUp, false)
	p.toolbar.Append(p.address, true)
	p.toolbar.Append(p.btnNext, false)
	p.toolbar.Append(p.btnNewTab, false)
	p.toolbar.Append(p.btnCloseTab, false)
	p.box.Append(p.title, false)
	p.box.Append(p.toolbar, false)
	p.box.Append(p.contents, true)
	p.box.Append(p.status, false)

	p.contentFont = ui.LoadClosestFont(&ui.FontDescriptor{
		Family: "Deja Vu",
		Size:   12,
	})
	p.titleFont = ui.LoadClosestFont(&ui.FontDescriptor{
		Family: "Deja Vu",
		Size:   12,
		Weight: ui.TextWeightHeavy,
	})
	p.titleFont = p.contentFont
	p.subtitleFont = p.contentFont
	p.summaryFont = ui.LoadClosestFont(&ui.FontDescriptor{
		Family: "Deja Vu",
		Size:   12,
		Italic: ui.TextItalicItalic,
	})
	p.verseFont = p.titleFont
	p.errorFont = ui.LoadClosestFont(&ui.FontDescriptor{
		Family: "Deja Vu",
		Size:   12,
	})

	return p
}

func (p *guiPage) Lookup(s string) {
	p.handleMessages(lib.Lookup(p.lang, s), true)
}

func toggleBtn(btn *ui.Button, item interface{}) {
	if item == nil {
		btn.Disable()
	} else {
		btn.Enable()
	}
}

func (p *guiPage) SetItem(item lib.Item, setText bool) {
	if item == nil {
		p.title.SetText("")
		p.btnUp.Disable()
		p.btnNext.Disable()
		p.btnPrevious.Disable()
		p.contents.SetDocument(nil)
	} else {
		toggleBtn(p.btnUp, item.Parent())
		toggleBtn(p.btnNext, item.Next())
		toggleBtn(p.btnPrevious, item.Previous())
		p.title.SetText(item.String())
		if setText {
			p.address.SetText(item.Path())
		}

		root := uidoc.NewGroup(make([]uidoc.Element, 0))

		defer func() {
			r := recover()
			if err, ok := r.(error); ok {
				root.Append(uidoc.NewText(err.Error(), p.errorFont))
			}
		}()

		root.Append(uidoc.NewText(item.Name(), p.titleFont))

		if children, err := item.Children(); err == nil {
			//TODO Add children
			for _, child := range children {
				func(child lib.Item) {
					text := uidoc.NewText(child.Name(), p.contentFont)
					text.PaddingLeft = 5
					text.PaddingRight = 5
					text.PaddingTop = 5
					text.PaddingBottom = 5
					text.MarginTop = 3
					text.MarginRight = 3
					text.MarginBottom = 3
					text.MarginLeft = 3
					text.LayoutMode = uidoc.LayoutInline
					text.Wrap = false
					button := uidoc.NewButton(text, func() {
						p.SetItem(child, true)
					})
					root.Append(button)
				}(child)
			}
		}

		if node, ok := item.(*lib.Node); ok {
			if content, err := node.Content(); err == nil {
				if page, err := content.Page(); err == nil {
					if len(page.Subtitle) > 0 {
						root.Append(uidoc.NewText(page.Subtitle, p.subtitleFont))
					}
					if len(page.Summary) > 0 {
						root.Append(uidoc.NewText(page.Summary, p.summaryFont))
					}
					for _, v := range page.Verses {
						verse := uidoc.NewText(fmt.Sprintf("%v", v.Number), p.verseFont)
						verse.LayoutMode = uidoc.LayoutInline
						verse.MarginRight = 5
						root.Append(verse)
						root.Append(uidoc.NewText(v.Text, p.contentFont))
					}
				}
			}
		}

		root.MarginTop = 20
		root.MarginLeft = 20
		root.MarginRight = 20
		root.MarginBottom = 100

		p.contents.SetDocument(root)
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
	p.handleMessages(lib.Lookup(p.lang, sender.Text()), false)
}
