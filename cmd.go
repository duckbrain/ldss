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

func (app *cmd) register(config *Configuration) {}

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

func (app *cmd) run() {
	args := app.args
	efmt := log.New(os.Stderr, "", 0)
	lang := app.lang

	switch args[0] {
	case "lookup":
		lookupString := strings.Join(args[1:], " ")
		refs := lib.Parse(lang, lookupString)
		if len(refs) != 1 {
			panic(fmt.Errorf("Multiple references not implemented"))
		}
		item, err := refs[0].Lookup()
		if err != nil {
			panic(err)
		}

		if node, ok := item.(*lib.Node); ok {
			if content, err := node.Content(); err == nil {
				c := colors(true)
				page := content.Page()
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
			catalog, err := lang.Catalog()
			if err != nil {
				panic(err)
			}
			fmt.Println(catalog.String())
		}
	case "download", "dl":
		if len(args) == 1 {
			lib.DownloadAll(lang, false)
			return
		}
		switch args[1] {
		case "languages", "lang":
			efmt.Println("Downloading language list")
			if err := lib.DownloadLanguages(); err != nil {
				panic(err)
			}
		case "all":
			lib.DownloadAll(lang, true)
		case "missing":
			lib.DownloadAll(lang, false)
		case "cat", "catalog":
			efmt.Println("Downloading \"" + lang.Name + "\" language catalog")
			lib.DownloadCatalog(lang)
		default:
			item, err := lib.AutoDownload(func() (lib.Item, error) {
				refs := lib.Parse(lang, args[1])
				if len(refs) != 1 {
					return nil, fmt.Errorf("Cannot yet handle multiple references")
				} else {
					return refs[0].Lookup()
				}
			})
			if err != nil {
				panic(err)
			}
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
