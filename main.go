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
	config := LoadConfiguration(op)

	if len(args) == 0 {
		PrintInstructions()
	} else {
		switch args[0] {
		case "help":
			if len(args) == 1 {
				PrintInstructions()
			} else {
				PrintCommandInstructions(args[1])
			}
		case "lookup":
			LookupPath(args[1])
		case "languages", "lang", "langs":
			if len(args) == 1 {
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
			default:
				if language := config.Languages.GetByUnknown(args[1]); language != nil {
					config.Download.DownloadCatalog(language)
				} else {
					panic("Unknown download \"" + args[1] + "\"")
				}
			}
		default:
			fmt.Printf("Unknown command \"%s\"\n", args[0])
		}
	}
}
