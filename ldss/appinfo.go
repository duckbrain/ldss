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

var apps map[string]app

func addApp(name string, a app) {
	if apps == nil {
		apps = make(map[string]app)
	}
	apps[name] = a
}
