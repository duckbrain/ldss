package main

import (
	"os"
	"log"
)

type app interface {
	run()
	setInfo(args []string, config Config)
}

type appinfo struct {
	args []string
	config Config
	fmt *log.Logger
	efmt *log.Logger
}

func (a *appinfo) setInfo(args []string, config Config) {
	a.args = args;
	a.config = config
	a.fmt = log.New(os.Stdin, "", 0)
	a.efmt = log.New(os.Stderr, "", 0)
}