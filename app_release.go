// +build !debug

package main

import (
	"io/ioutil"
	"log"
)

func (a *appinfo) getDebug() *log.Logger {
	return log.New(ioutil.Discard, "", 0)
}
