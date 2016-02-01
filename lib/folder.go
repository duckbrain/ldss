package lib

import (
	"fmt"
)

type folderBase struct {
	base    *jsonFolder
	folders []*Folder
	books   []*Book
}

type Folder struct {
	folderBase
	parent  Item
	catalog *Catalog
}

func (f *folderBase) Children() ([]Item, error) {
	folderLen := len(f.base.Folders)
	items := make([]Item, folderLen+len(f.base.Books))
	for i, f := range f.Folders() {
		items[i] = f
	}
	for i, f := range f.Books() {
		items[folderLen+i] = f
	}
	return items, nil
}

func (f *folderBase) Folders() []*Folder {
	return f.folders
}

func (f *folderBase) Books() []*Book {
	return f.books
}

func (f *Folder) ID() int {
	return f.base.ID
}

func (f *Folder) String() string {
	return fmt.Sprintf("%v {folders[%v] books[%v]}", f.Name(), len(f.Folders()), len(f.Books()))
}

func (f *folderBase) Name() string {
	return f.base.Name
}

func (f *Folder) Path() string {
	//TODO: Calculate path based on commonality with children
	return fmt.Sprintf("/%v", f.ID)
}

func (f *Folder) Language() *Language {
	return f.catalog.language
}

func (f *Folder) Parent() Item {
	return f.parent
}
