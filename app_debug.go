// +build debug

package main

import (
	"log"
	"os"
)

func (a *appinfo) getDebug() *log.Logger {
	return log.New(os.Stderr, "", 0)
}
