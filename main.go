package main

import (
	"fmt"
	"ldss/lib"
)

func main() {

	Config().RegisterOption(ConfigOption{
		Name:     "Language",
		Default:  "eng",
		ShortArg: 'l',
		LongArg:  "lang",
	})

	if err := Config().Init(); err != nil {
		panic(err)
	}

	lib.SetReferenceParseReader(func(lang *lib.Language) ([]byte, error) {
		return Asset("data/reference/" + lang.GlCode)
	})

	args := Config().Args()

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
			fmt.Print(Config().String())
		default:
			var ok bool
			var app app
			if app, ok = apps[args[0]]; !ok {
				app = &cmd{}
			}
			app.setInfo(args)
			app.run()
		}
	}
}
