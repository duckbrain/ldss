package lib

import "fmt"

// The Generic interface for all items in the gospel library. This should be a
// *Catalog, *Folder, *Book, or *Node. For the most part, an interface should be
// possible without converting this interface to it's base type, except to convert
// to a Node to get it's content and media.
type Item interface {
	Name() string
	Children() ([]Item, error)
	Path() string
	Language() *Language
	Parent() Item
	Next() Item
	Previous() Item
	Search(<-chan Reference, Reference)
	fmt.Stringer
}
