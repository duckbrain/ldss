package main

import (
	"log"
)

type app interface {
	run()
	setInfo(args []string, config Config)
}

type appinfo struct {
	args   []string
	config Config
	fmt    *log.Logger
	efmt   *log.Logger
	debug  *log.Logger
}
