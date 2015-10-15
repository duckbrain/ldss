package main

import (
	"fmt"
	"os"
	"github.com/wsxiaoys/terminal/color"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			color.Println("@rerror@{|}: " + err.Error())
		}
	}()
	
	args := os.Args[1:]
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
		case "download", "dl":
			switch (args[1]) {
			case "languages", "lang":
				DownloadLanguages()
			default:
				panic("Unknown download");
			}
		default:
			fmt.Printf("Unknown command \"%s\"\n", args[0])
		}
	}
}
