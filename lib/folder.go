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

// Represents a folder in the catalog. Could contain subfolders and books.
type Folder struct {
	folderBase
	parent  Item
	catalog *Catalog
}

// Combined Folders and Books
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

// Folders as direct children of this item
func (f *folderBase) Folders() []*Folder {
	return f.folders
}

// Books as direct children of this item
func (f *folderBase) Books() []*Book {
	return f.books
}

// Name of this item
func (f *folderBase) Name() string {
	return f.base.Name
}

// An ID that is unique to this Folder within it's language
func (f *Folder) ID() int {
	return f.base.ID
}

// A short human-readable representation of the folder, mostly useful for debugging.
func (f *Folder) String() string {
	return fmt.Sprintf("%v {%v folders[%v] books[%v]}", f.Name(), f.Path(), len(f.Folders()), len(f.Books()))
}

// Full path of this folder. It will attempt to get a path from the references
// file or create a path based on the names of it's children. As a last resort,
// it will prepend it's ID with a forward slash.
func (f *Folder) Path() string {
	//Calculate path based on commonality with children
	var childFound = false
	var path []string
	var search func(folder *Folder)

	if p, err := referenceParser(f.Language()); err == nil {
		if path, ok := p.matchFolder[f.ID()]; ok {
			return path
		}
	}

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
		if found, ok := f.catalog.foldersByPath[p]; !ok || found == f {
			return p
		}
	}

	return fmt.Sprintf("/%v", f.ID())
}

// Language of this folder
func (f *Folder) Language() *Language {
	return f.catalog.language
}

// Parent of this folder. Either a catalog or another folder
func (f *Folder) Parent() Item {
	return f.parent
}

// Next sibling of this folder
func (f *Folder) Next() Item {
	return genericNextPrevious(f, 1)
}

// Previous sibling of this folder
func (f *Folder) Previous() Item {
	return genericNextPrevious(f, -1)
}
