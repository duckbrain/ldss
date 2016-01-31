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
	title, subtitle, summary, verse, content *color.Color
}

func colors(enabled bool) *cmdcolors {
	c := cmdcolors{}
	c.content = color.New()
	if enabled {
		c.title = color.New(color.Bold).Add(color.Underline).Add(color.FgWhite).Add(color.BgHiMagenta)
		c.subtitle = color.New(color.Bold).Add(color.FgGreen)
		c.summary = color.New(color.Italic).Add(color.BgBlue).Add(color.FgBlack)
		c.verse = color.New(color.Bold).Add(color.FgRed)
	} else {
		c.title = c.content
		c.subtitle = c.content
		c.summary = c.content
		c.verse = c.content
	}
	return &c
}

func (app *cmd) run() {
	c := colors(true)
	args := app.args
	efmt := log.New(os.Stderr, "", 0)
	var catalog *lib.Catalog

	for m := range lib.DefaultCatalog() {
		switch m.(type) {
		case lib.MessageDone:
			catalog = m.(lib.MessageDone).Item().(*lib.Catalog)
		default:
			efmt.Println(m.String())
		}
	}

	switch args[0] {
	case "lookup":
		lookupString := strings.Join(args[1:], " ")
		item, err := lib.Lookup(lookupString, catalog)
		if err != nil {
			efmt.Printf("Path \"%v\" not found.", lookupString)
			panic(err)
		}

		if node, ok := item.(lib.Node); ok {
			if content, err := node.Content(); err == nil {
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

		children, err := item.Children()
		if err != nil {
			panic(err)
		}
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
			//lang := config.SelectedLanguage()
			//efmt.Println("Downloading all content for \"" + lang.Name + "\" catalog")
			//config.Download.Missing()
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
