package main

import (
	"fmt"
	"github.com/fatih/color"
)

var _ color.Color
var apps map[string]app

func init() {
	apps = make(map[string]app)
}

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
			fmt.Print(op)
		default:
			var ok bool
			var app app
			if app, ok = apps[args[0]]; !ok {
				app = &cmd{}
			}
			app.setInfo(args, config)
			app.run()
		}
	}
}
