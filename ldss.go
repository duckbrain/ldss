package main

import (
	"fmt"
	"os"
	//"github.com/wsxiaoys/terminal/color"
)

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
				lang, err := config.Library.Language(args[1])
				if err == os.ErrNotExist {
					fmt.Printf("Languages have not been downloaded. Please run ldss dl lang to download.")
					return
				}
				fmt.Println(lang.String())
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
