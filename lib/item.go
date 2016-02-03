package lib

import "fmt"

type Item interface {
	Name() string
	Children() ([]Item, error)
	Path() string
	Language() *Language
	Parent() Item
	Next() Item
	Previous() Item
	fmt.Stringer
}
