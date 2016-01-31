package main

import (
	"ldss/lib"
	"log"
)

type app interface {
	run()
	setInfo(args []string, config *lib.Configuration)
}

type appinfo struct {
	args   []string
	config *lib.Configuration
	fmt    *log.Logger
	efmt   *log.Logger
	debug  *log.Logger
}
