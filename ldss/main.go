package main

import (
	"fmt"
	"ldss/lib"

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
			color.Println("@rfatal error@{|}: " + err.Error())
		}
	}()*/

	if err := lib.Config().Init(); err != nil {
		panic(err)
	}

	args := lib.Config().Args()

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
			fmt.Print(lib.Config().String())
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
