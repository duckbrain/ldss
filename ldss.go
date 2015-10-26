package main

import (
	"fmt"
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
	loadFileOptions(op)
	args := loadParameterOptions(op)
	alen := len(args)
	config := LoadConfiguration(op)

	if len(args) == 0 {
		PrintInstructions()
	} else {
		switch args[0] {
		case "help":
			if alen == 1 {
				PrintInstructions()
			} else {
				for _, instr := range args[1:] {
					PrintCommandInstructions(instr)
				}
			}
		case "lookup":
			//l := NewLookupLoader(config.SelectedLanguage, config.OfflineContent)
			LookupPath(args[1])
		case "languages":
			if alen == 1 {
				for _, l := range config.Languages.GetAll() {
					fmt.Println(l.String())
				}
			} else {
				fmt.Println(config.Languages.GetByUnknown(args[1]).String())
			}
		case "catalog", "cat":
			catalog := NewCatalogLoader(config.SelectedLanguage, config.OfflineContent)

			if len(args) == 1 {
				for _, l := range config.Languages.GetAll() {
					fmt.Println(l.String())
				}
			} else {
				fmt.Println(catalog.GetCatalog().String())
			}
		case "download", "dl":
			switch args[1] {
			case "languages", "lang":
				config.Download.DownloadLanguages()
			case "all":
				panic("Not implemented")
			default:
				language := config.Languages.GetByUnknown(args[1])

				if language != nil {
					config.Download.DownloadCatalog(language)
				} else {
					catalog := NewCatalogLoader(config.SelectedLanguage, config.OfflineContent)
					book := catalog.GetBookByUnknown(args[1])
					if book != nil {
						config.Download.DownloadBook(book)
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
