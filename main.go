package main

import (
	"fmt"
	"github.com/duckbrain/ldss/lib"
)

func main() {

	config := newConfiguration()

	config.RegisterOption(ConfigOption{
		Name:     "Language",
		Default:  "eng",
		ShortArg: 'l',
		LongArg:  "lang",
	})

	if err := config.Init(); err != nil {
		panic(err)
	}

	lib.SetReferenceParseReader(func(lang *lib.Language) ([]byte, error) {
		return Asset("data/reference/" + lang.GlCode)
	})

	args := config.Args()

	if len(args) == 0 {
		PrintInstructions()
		var ok bool
		var app app
		if app, ok = apps["web"]; !ok {
			app = &cmd{}
		}
		app.register(config)
		app.setInfo(config, args)
		app.run()
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
			fmt.Print(config.String())
		default:
			var ok bool
			var app app
			if app, ok = apps[args[0]]; !ok {
				app = &cmd{}
			}
			app.register(config)
			app.setInfo(config, args)
			app.run()
		}
	}
}
