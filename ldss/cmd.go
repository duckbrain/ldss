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
	//c := colors(true)
	args := app.args
	efmt := log.New(os.Stderr, "", 0)
	catalog := app.item(lib.DefaultCatalog()).(*lib.Catalog)

	switch args[0] {
	case "lookup":
		lookupString := strings.Join(args[1:], " ")
		item, err := catalog.Lookup(lookupString)
		if err != nil {
			efmt.Printf("Path \"%v\" not found.", lookupString)
			panic(err)
		}

		if node, ok := item.(lib.Node); ok {
			if content, err := node.Content(); err == nil {
				app.fmt.Printf("%v", content.HTML())
				//TODO: Format
				/*
					c.title.Printf("   %v   \n", content.Title)
					if len(content.Subtitle) > 0 {
						c.subtitle.Println(content.Subtitle)
					}
					if len(content.Summary) > 0 {
						c.summary.Println(content.Summary)
					}
					for _, verse := range content.Verses {
						c.verse.Printf("%v ", verse.Number)
						c.content.Println(verse.Text)
					}*/
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
			fmt.Println(catalog.String())
		}
	case "download", "dl":
		if len(args) == 1 {
			efmt.Println("Must provide argment of what to download")
			efmt.Println("usage: ldss download|dl lang|<lang>|<book>")
			return
		}
		switch args[1] {
		case "languages", "lang":
			efmt.Println("Downloading language list")
			if err := lib.DownloadLanguages(); err != nil {
				panic(err)
			}
		case "all":
			app.item(lib.DownloadAll(catalog.Language(), true))
		case "missing":
			app.item(lib.DownloadAll(catalog.Language(), false))
		case "cat", "catalog":
			lang, err := lib.DefaultLanguage()
			if err != nil {
				panic(err)
			}
			efmt.Println("Downloading \"" + lang.Name + "\" language catalog")
			lib.DownloadCatalog(lang)
		default:
			book, err := catalog.LookupBook(args[1])
			if err != nil {
				panic("Unknown download \"" + args[1] + "\"")
			}
			efmt.Printf("Downloading book \"%v\" for the \"%v\" catalog\n", book.Name, catalog.Name)
			lib.DownloadBook(book)
		}
	default:
		app.args = append([]string{"lookup"}, app.args...)
		app.run()
		//fmt.Printf("Unknown command \"%s\"\n", args[0])
	}
}
