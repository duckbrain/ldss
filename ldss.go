package main

import (
	"fmt"
	"os"
	"log"
	"github.com/fatih/color"
)

var _ color.Color

func main() {
	/*defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			color.Println("@rerror@{|}: " + err.Error())
		}
	}()*/

	op := loadDefaultOptions()
	op = loadFileOptions(op)
	op, args := loadParameterOptions(op)
	config := LoadConfiguration(op)
	efmt := log.New(os.Stderr, "", 0)

	if len(args) == 0 {
		PrintInstructions()
	} else {
		switch args[0] {
		case "help":
			if len(args) == 1 {
				PrintInstructions()
			} else {
				for _, instr := range args[1:] {
					PrintCommandInstructions(instr)
				}
			}
		case "config":
			fmt.Printf("Language:      %v\n", op.Language)
			fmt.Printf("ServerURL:     %v\n", op.ServerURL)
			fmt.Printf("DataDirectory: %v\n", op.DataDirectory)
			fmt.Printf("WebPort:       %v\n", op.WebPort)
		case "lookup":
			//l := NewLookupLoader(config.SelectedLanguage(), config.OfflineContent)
			//LookupPath(args[1])
		case "languages":
			if len(args) == 1 {
				for _, l := range config.Languages() {
					fmt.Println(l.String())
				}
			} else {
				fmt.Println(config.Language(args[1]).String())
			}
		case "catalog", "cat":
			catalog := config.Library.Catalog(config.SelectedLanguage())

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
				config.Download.Languages()
			case "all":
				panic("Not implemented")
			default:
				language := config.Language(args[1])

				if language != nil {
					config.Download.Catalog(language)
				} else {
					catalog := config.Library.Catalog(language)
					book := config.Library.Book(args[1], catalog)
					if book != nil {
						config.Download.Book(book)
					} else {
						panic("Unknown download \"" + args[1] + "\"")
					}
				}
			}
		default:
			fmt.Printf("Unknown command \"%s\"\n", args[0])
		}
	}
}
