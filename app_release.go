// +build release

package main

import (
	"io/ioutil"
	"log"
	"os"
)

func (a *appinfo) getDebug() *log.Logger {
	return log.New(ioutil.Discard, "", 0)
}
