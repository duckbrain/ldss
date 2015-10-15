
package main

import (
	"fmt"
	"errors"
	"github.com/wsxiaoys/terminal/color"
)

type Reference struct {
	bookName string
	glPath string
	node int
	chapter int
	verseSelected int
	versesHighlighted []int
}

func (ref Reference) String() string {
	return string(ref.bookName) + " " + string(ref.chapter) + ":" + string(ref.verseSelected)
}

func ParsePath(path string) (ref *Reference, err error) {
	return nil, errors.New("Not Implemented")
}

func LookupPath(path string) {
	ref, err := ParsePath(path)

	if err != nil {
		color.Println("@rerror@{|}: " + err.Error())
	} else {
		fmt.Println(ref.String())
	}
}
