package main

import (
	"fmt"
	"ldss/lib"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

type cmd struct {
	appinfo
	colors cmdcolors
}

type cmdcolors struct {
	title, subtitle, summary, verse, content, message *color.Color
}

func colors(enabled bool) *cmdcolors {
	c := cmdcolors{}
	c.content = color.New()
	c.title = color.New(color.Bold).Add(color.Underline).Add(color.FgWhite).Add(color.BgHiMagenta)
	c.subtitle = color.New(color.Bold).Add(color.FgGreen)
	c.summary = color.New(color.Italic).Add(color.BgBlue).Add(color.FgBlack)
	c.verse = color.New(color.Bold).Add(color.FgRed)
	c.message = color.New(color.FgHiYellow).Add(color.Italic)
	color.NoColor = !enabled
	return &c
}

func (app *cmd) item(c <-chan lib.Message) interface{} {
	for m := range c {
		switch m.(type) {
		case lib.MessageDone:
			return m.(lib.MessageDone).Item()
		case lib.MessageError:
			panic(m)
		default:
			if m == nil {
				return nil
			}
			fmt.Printf("%v\n", m)
		}
	}
	panic(fmt.Errorf("Channel completed prematurely\n"))
}

func (app *cmd) dl(open func() (interface{}, error)) interface{} {
	return app.item(lib.AutoDownload(open))
}

func (app *cmd) run() {
	args := app.args
	efmt := log.New(os.Stderr, "", 0)
	lang, err := lib.DefaultLanguage()

	if err != nil {
		panic(err)
	}

	switch args[0] {
	case "lookup":
		lookupString := strings.Join(args[1:], " ")
		item := app.item(lib.Lookup(lang, lookupString)).(lib.Item)

		if node, ok := item.(*lib.Node); ok {
			if content, err := node.Content(); err == nil {
				c := colors(true)
				page, err := content.Page()
				if err != nil {
					panic(err)
				}
				c.title.Printf("   %v   \n", page.Title)
				if len(page.Subtitle) > 0 {
					c.subtitle.Println(page.Subtitle)
				}
				if len(page.Summary) > 0 {
					c.summary.Println(page.Summary)
				}
				for _, verse := range page.Verses {
					c.verse.Printf("%v ", verse.Number)
					c.content.Println(verse.Text)
				}
				break
			}
		}

		children := app.dl(func() (interface{}, error) {
			return item.Children()
		}).([]lib.Item)
		fmt.Println(item)
		for _, child := range children {
			fmt.Printf("- %v\n", child)
		}
	case "languages":
		if len(args) == 1 {
			langs, err := lib.Languages()
			if err != nil {
				panic(err)
			}
			for _, l := range langs {
				fmt.Println(l.String())
			}
		} else {
			lang, err := lib.LookupLanguage(args[1])
			if err != nil {
				panic(err)
			}
			fmt.Println(lang.String())
		}
	case "catalog", "cat":
		if len(args) == 1 {
			langs, err := lib.Languages()
			if err != nil {
				panic(err)
			}
			for _, l := range langs {
				fmt.Println(l.String())
			}
		} else {
			catalog := app.item(lib.Lookup(lang, "/")).(*lib.Catalog)
			fmt.Println(catalog.String())
		}
	case "download", "dl":
		if len(args) == 1 {
			app.item(lib.DownloadAll(lang, false))
			return
		}
		switch args[1] {
		case "languages", "lang":
			efmt.Println("Downloading language list")
			if err := lib.DownloadLanguages(); err != nil {
				panic(err)
			}
		case "all":
			app.item(lib.DownloadAll(lang, true))
		case "missing":
			app.item(lib.DownloadAll(lang, false))
		case "cat", "catalog":
			efmt.Println("Downloading \"" + lang.Name + "\" language catalog")
			lib.DownloadCatalog(lang)
		default:
			item := app.item(lib.Lookup(lang, args[1]))
			if book, ok := item.(*lib.Book); ok {
				efmt.Printf("Downloading book \"%v\"\n", book.Name())
				lib.DownloadBook(book)
			} else if folder, ok := item.(*lib.Folder); ok {
				efmt.Printf("Downloading folder \"%v\"\n", folder.Name())
				lib.DownloadChildren(folder, false)
			}
		}
	default:
		app.args = append([]string{"lookup"}, app.args...)
		app.run()
	}
}
