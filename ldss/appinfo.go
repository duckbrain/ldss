package main

import (
	"log"
)

type app interface {
	run()
	setInfo(args []string)
}

type appinfo struct {
	args  []string
	fmt   *log.Logger
	efmt  *log.Logger
	debug *log.Logger
}
