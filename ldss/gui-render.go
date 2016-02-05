package main

import (
	"fmt"
	"ldss/lib"

	"github.com/andlabs/ui"
)

type guiRenderer struct {
	area                                                         *ui.Area
	item                                                         lib.Item
	page                                                         *lib.Page
	elements                                                     []guiRenderElement
	titleFont, subtitleFont, summaryFont, verseFont, contentFont *ui.Font
	width, height, scrollY                                       float64
}

type guiRenderElement struct {
	layout *ui.TextLayout
	x, y   float64
}

func newGuiRenderer() *guiRenderer {
	r := &guiRenderer{}
	//TODO Scrolling Area
	r.area = ui.NewArea(r)
	r.contentFont = ui.LoadClosestFont(&ui.FontDescriptor{
		Family: "Ubuntu",
		Size:   12,
	})
	r.titleFont = ui.LoadClosestFont(&ui.FontDescriptor{
		Family: "Ubuntu",
		Size:   12,
		Weight: ui.TextWeightBold,
	})
	r.titleFont = r.contentFont
	r.subtitleFont = r.contentFont
	r.summaryFont = ui.LoadClosestFont(&ui.FontDescriptor{
		Family: "Ubuntu",
		Size:   12,
		Italic: ui.TextItalicItalic,
	})
	r.verseFont = r.titleFont
	return r
}

func (r *guiRenderer) SetItem(item lib.Item) error {
	r.item = item
	if node, ok := item.(*lib.Node); ok {
		content, err := node.Content()
		if err != nil {
			return err
		}
		r.page, err = content.Page()
		if err != nil {
			r.item = nil
			r.page = nil
			return err
		}
	}
	
	// Add the elements
	if r.elements != nil {
		for _, ele := range r.elements {
			ele.layout.Free()
		}
	}
	if r.page == nil {
		return nil
	}
	
	r.elements = make([]guiRenderElement, 3+len(r.page.Verses))
	r.elements[0] = guiRenderElement{
		layout: ui.NewTextLayout(r.page.Title, r.titleFont, r.width),
	}
	r.elements[1] = guiRenderElement{
		layout: ui.NewTextLayout(r.page.Subtitle, r.subtitleFont, r.width),
	}
	r.elements[2] = guiRenderElement{
		layout: ui.NewTextLayout(r.page.Summary, r.summaryFont, r.width),
	}
	for i, v := range r.page.Verses {
		//TODO Add verse number with float left
		r.elements[i+3] = guiRenderElement{
			layout: ui.NewTextLayout(v.Text, r.contentFont, r.width),
		}
	}

	r.width = 0
	r.area.QueueRedrawAll()
	return nil
}

func (r *guiRenderer) layout(width float64) {
	if r.width == width {
		return
	}
	

	y := 0.0
	for i, ele := range r.elements {
		ele.y = y
		_, h := ele.layout.Extents()
		y += h
		r.elements[i] = ele
	}

	r.width = width
	r.height = y
}

func (r *guiRenderer) Draw(a *ui.Area, dp *ui.AreaDrawParams) {
	//fmt.Printf("Area Size: %v, %v Clip Box: %v, %v, %v, %v\n", dp.AreaHeight, dp.AreaWidth, dp.ClipX, dp.ClipY, dp.ClipWidth, dp.ClipHeight)
	r.layout(dp.AreaWidth)
	if r.elements == nil {
		return
	}
	for _, e := range r.elements {
		dp.Context.Text(e.x, e.y - r.scrollY, e.layout)
	}
}

func (r *guiRenderer) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
	fmt.Printf("Up/Down:%v/%v Count:%v Modifiers:%v Held:%v \n", me.X, me.Y, me.AreaWidth, me.AreaHeight, me.Up, me.Down, me.Count, me.Modifiers, me.Held);
}

func (r *guiRenderer) MouseCrossed(a *ui.Area, left bool) {

}

func (r *guiRenderer) DragBroken(a *ui.Area) {

}

func (r *guiRenderer) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	handled = true
	switch ke.Key {
		case 'j':
			r.scrollY -= 10
			r.area.QueueRedrawAll()
		case 'k':
			r.scrollY += 10
			r.area.QueueRedrawAll()
		default:
			return false
	}
	return
}

// Wrap the Area to make this element behave as a control

func (r *guiRenderer) Destroy() {
	r.area.Destroy()
}
func (r *guiRenderer) LibuiControl() uintptr {
	return r.area.LibuiControl()
}
func (r *guiRenderer) Handle() uintptr {
	return r.area.Handle()
}
func (r *guiRenderer) Show() {
	r.area.Show()
}
func (r *guiRenderer) Hide() {
	r.area.Hide()
}
func (r *guiRenderer) Enable() {
	r.area.Enable()
}
func (r *guiRenderer) Disable() {
	r.area.Disable()
}
