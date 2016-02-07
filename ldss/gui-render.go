package main

import (
	"fmt"
	"ldss/lib"

	"github.com/andlabs/ui"
)

var _ fmt.Formatter

type guiRenderer struct {
	area, measureArea                   *ui.Area
	box                                 *ui.Box
	item                                lib.Item
	page                                *lib.Page
	elements                            []guiRenderElement
	titleFont, subtitleFont             *ui.Font
	summaryFont, verseFont              *ui.Font
	contentFont                         *ui.Font
	width, height, scrollY              float64
	marginTop, marginBottom, marginSide float64
	onItemChange                        func(lib.Item, *guiRenderer)
}

type guiRenderElement struct {
	layout *ui.TextLayout
	x, y   float64
	inline bool
}

func newGuiRenderer() *guiRenderer {
	r := &guiRenderer{
		marginTop:    20,
		marginBottom: 20,
		marginSide:   20,
	}
	// Scrolling Area
	r.area = ui.NewScrollingArea(r, 400, 400)
	r.measureArea = ui.NewArea(&guiRenderMeasure{r})
	r.box = ui.NewVerticalBox()
	toolbarContainer := ui.NewHorizontalBox()
	toolbarHeightRetainer := ui.NewHorizontalSeparator()
	toolbarContainer.Append(toolbarHeightRetainer, false)
	toolbarContainer.Append(r.measureArea, true)
	r.box.Append(r.area, true)
	r.box.Append(toolbarContainer, false)

	r.contentFont = ui.LoadClosestFont(&ui.FontDescriptor{
		Family: "Deja Vu",
		Size:   12,
	})
	r.titleFont = ui.LoadClosestFont(&ui.FontDescriptor{
		Family: "Deja Vu",
		Size:   12,
		Weight: ui.TextWeightHeavy,
	})
	r.titleFont = r.contentFont
	r.subtitleFont = r.contentFont
	r.summaryFont = ui.LoadClosestFont(&ui.FontDescriptor{
		Family: "Deja Vu",
		Size:   12,
		Italic: ui.TextItalicItalic,
	})
	r.verseFont = r.titleFont
	return r
}

func (r *guiRenderer) SetItem(item lib.Item) error {
	r.item = item
	if r.onItemChange != nil {
		r.onItemChange(item, r)
	}
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
	} else {
		r.page = nil
		children, err := item.Children()

	}

	// Add the elements
	if r.elements != nil {
		for _, ele := range r.elements {
			ele.layout.Free()
		}
	}
	r.elements = make([]guiRenderElement, 3+len(r.page.Verses)*2)

	if r.page == nil {
		return nil
	}

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
		r.elements[i*2+3] = guiRenderElement{
			layout: ui.NewTextLayout(fmt.Sprintf("%v ", v.Number), r.verseFont, r.width),
			inline: true,
		}
		r.elements[i*2+4] = guiRenderElement{
			layout: ui.NewTextLayout(v.Text, r.contentFont, r.width),
		}
	}

	//r.width = 0
	r.layout(r.width)
	r.measureArea.QueueRedrawAll()
	return nil
}

func (r *guiRenderer) layout(width float64) {
	x := r.marginSide
	y := r.marginTop
	for i, ele := range r.elements {
		ele.y = y
		ele.x = x
		ele.layout.SetWidth(width - r.marginSide - x)
		w, h := ele.layout.Extents()
		if ele.inline {
			x += w
		} else {
			y += h
			x = r.marginSide
			fmt.Printf("Layout: %v, %v, %v, %v\n", x, y, w, h)
		}

		r.elements[i] = ele
	}

	r.width = width
	r.height = y + r.marginBottom
	r.area.QueueRedrawAll()
	r.area.SetSize(int(r.width), int(r.height))
}

func (r *guiRenderer) Draw(a *ui.Area, dp *ui.AreaDrawParams) {
	//fmt.Printf("Area Size: %v, %v Clip Box: %v, %v, %v, %v\n", dp.AreaHeight, dp.AreaWidth, dp.ClipX, dp.ClipY, dp.ClipWidth, dp.ClipHeight)
	//r.layout(dp.AreaWidth)
	if r.elements == nil {
		return
	}
	for _, e := range r.elements {
		dp.Context.Text(e.x, e.y-r.scrollY, e.layout)
	}
}

func (r *guiRenderer) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
	//fmt.Printf("Up/Down:%v/%v Count:%v Modifiers:%v Held:%v \n", me.X, me.Y, me.AreaWidth, me.AreaHeight, me.Up, me.Down, me.Count, me.Modifiers, me.Held)
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

// Wrap the box to make this element behave as a control

func (r *guiRenderer) Destroy() {
	r.box.Destroy()
}
func (r *guiRenderer) LibuiControl() uintptr {
	return r.box.LibuiControl()
}
func (r *guiRenderer) Handle() uintptr {
	return r.box.Handle()
}
func (r *guiRenderer) Show() {
	r.box.Show()
}
func (r *guiRenderer) Hide() {
	r.box.Hide()
}
func (r *guiRenderer) Enable() {
	r.box.Enable()
}
func (r *guiRenderer) Disable() {
	r.box.Disable()
}

type guiRenderMeasure struct {
	parent *guiRenderer
}

func (r *guiRenderMeasure) Draw(a *ui.Area, dp *ui.AreaDrawParams) {
	if r.parent != nil && r.parent.width != dp.AreaWidth {
		fmt.Printf("Size change: %v\n", dp.AreaWidth)
		r.parent.layout(dp.AreaWidth)
	}

	//Fill background
	p := ui.NewPath(ui.Winding)
	p.AddRectangle(0, 0, dp.AreaWidth, dp.AreaHeight)
	p.End()
	dp.Context.Fill(p, &ui.Brush{
		A: 1,
	})
}

func (r *guiRenderMeasure) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {

}

func (r *guiRenderMeasure) MouseCrossed(a *ui.Area, left bool) {

}

func (r *guiRenderMeasure) DragBroken(a *ui.Area) {

}

func (r *guiRenderMeasure) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) bool {
	return false
}
