package main

import (
	"fmt"
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
		case "web":
			web{args, config}.run()
		case "gui":
			gui(args, config)
		case "curses":
			curses{args, config}.run()
		default:
			cmd(args, config)
		}
	}
}
