package lib

import (
	"fmt"
	"strings"
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
	return fmt.Sprintf("%v {%v folders[%v] books[%v]}", f.Name(), f.Path(), len(f.Folders()), len(f.Books()))
}

func (f *folderBase) Name() string {
	return f.base.Name
}

func (f *Folder) Path() string {
	//Calculate path based on commonality with children
	var childFound = false
	var path []string
	var search func(folder *Folder)

	search = func(folder *Folder) {
		for _, book := range folder.books {
			p := strings.Split(book.Path(), "/")
			if !childFound {
				path = p
				childFound = true
				continue
			}
			for i := 0; i < len(p) && i < len(path); i++ {
				if p[i] != path[i] {
					path = path[0:i]
				}
			}
		}
		for _, subFolder := range folder.folders {
			search(subFolder)
		}
	}

	search(f)

	if childFound && len(path) > 1 {
		p := strings.Join(path, "/")
		if _, ok := f.catalog.foldersByPath[p]; !ok {
			return p
		}
	}

	return fmt.Sprintf("/%v", f.ID())
}

func (f *Folder) Language() *Language {
	return f.catalog.language
}

func (f *Folder) Parent() Item {
	return f.parent
}
