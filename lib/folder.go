package lib

import (
	"fmt"
)

/*
 * Folder
 */

type Folder struct {
	base    jsonFolder
	parent  Item
	catalog *Catalog
}

func (f *Folder) Children() ([]Item, error) {
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

func (f *Folder) Folders() []*Folder {
	return f.base.Folders
}

func (f *Folder) Books() []*Book {
	return f.base.Books
}

func (f *Folder) ID() int {
	return f.base.ID
}

func (f *Folder) String() string {
	return fmt.Sprintf("%v {folders[%v] books[%v]}", f.Name(), len(f.Folders()), len(f.Books()))
}

func (f *Folder) Name() string {
	return f.base.Name
}

func (f *Folder) Path() string {
	return fmt.Sprintf("/%v", f.ID)
}

func (f *Folder) Language() *Language {
	return f.catalog.language
}

func (f *Folder) Parent() Item {
	return f.parent
}
