// +build debug

package main

import (
	"ldss/lib"
	"log"
	"os"
)

func (a *appinfo) setInfo(args []string, config *lib.Configuration) {
	a.args = args
	a.config = config
	a.fmt = log.New(os.Stdin, "", 0)
	a.efmt = log.New(os.Stderr, "", 0)
	a.debug = log.New(os.Stderr, "", 0)
}
