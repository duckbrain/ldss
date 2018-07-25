package ldsorg

import (
	"fmt"
	"strings"

	"github.com/duckbrain/ldss/lib"
)

type Item = lib.Item

// Represents a folder in the catalog. Could contain subfolders and books.
type folder struct {
	jsonFolder
	parent  Item
	catalog *catalog
}

// Full path of this folder. It will attempt to get a path from the references
// file or create a path based on the names of it's children. As a last resort,
// it will prepend it's ID with a forward slash.
func (f *Folder) Path() string {
	//Calculate path based on commonality with children
	var childFound = false
	var path []string
	var search func(folder *Folder)

	if p, err := languageQueryParser(f.Language()); err == nil {
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
func (f *Folder) Lang() Lang {
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
func (f *Folder) Prev() Item {
	return genericNextPrevious(f, -1)
}
