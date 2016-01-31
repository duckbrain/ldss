// +build !release

package main

import (
	"log"
	"os"
)

func (a *appinfo) setInfo(args []string) {
	a.args = args
	a.fmt = log.New(os.Stdin, "", 0)
	a.efmt = log.New(os.Stderr, "", 0)
	a.debug = log.New(os.Stderr, "", 0)
}
