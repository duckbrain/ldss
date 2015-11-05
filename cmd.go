package main

import (
	"fmt"
	"log"
	"ldslib"
	"strings"
	"os"
)

func cmd(args []string, config Config) {
	efmt := log.New(os.Stderr, "", 0)
	switch args[0] {
		case "lookup":
			item, err := config.Library.Lookup(strings.Join(args[1:], " "), config.SelectedCatalog())
			if err != nil {
				panic(err)
			}
			switch item.(type) {
				case ldslib.Node:
					fmt.Println(config.Library.Content(item.(ldslib.Node)))
				default:
					fmt.Println(item)
			}
		case "languages":
			if len(args) == 1 {
				for _, l := range config.Languages() {
					fmt.Println(l.String())
				}
			} else {
				fmt.Println(config.Language(args[1]).String())
			}
		case "catalog", "cat":
			catalog := config.SelectedCatalog()

			if len(args) == 1 {
				langs, err := config.Library.Languages()
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
				config.Download.Languages()
			case "all":
				lang := config.SelectedLanguage()
				efmt.Println("Downloading all content for \"" + lang.Name + "\" catalog")
				config.Download.Missing()
			case "cat", "catalog":
				lang := config.SelectedLanguage()
				efmt.Println("Downloading \"" + lang.Name + "\" language catalog")
				config.Download.Catalog(lang)
			default:
				catalog := config.SelectedCatalog()
				book, err := config.Library.Book(args[1], catalog)
				if err != nil {
					panic("Unknown download \"" + args[1] + "\"")
				}
				efmt.Printf("Downloading book \"%v\" for the \"%v\" catalog\n", book.Name, catalog.Name)
				config.Download.Book(book)
			}
		default:
			fmt.Printf("Unknown command \"%s\"\n", args[0])
	}
}