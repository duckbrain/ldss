package ldsorg

import (
	"strings"

	"github.com/duckbrain/ldss/lib"
	"github.com/duckbrain/ldss/lib/dl"
)

var itemsByLangAndPath = make(map[ref]lib.Item)

var _ lib.Source = (*source)(nil)
var _ dl.Downloader = (*source)(nil)

type source struct {
	dl.Template
}
type ref struct {
	langCode string
	path     string
}

func init() {
	s := source{}
	s.Template.Src = getServerAction("languages.query")
	s.Template.Dest = languagesPath()
	lib.Register("lds.org", s)
}

func (s source) Name() string {
	return "LDS.org"
}

func (s source) Hash() string {
	return "ldsorg"
}

// Lookup finds an Item by it's path. Expects a fully qualified path. "/" will
// return the catalog. Will return an error if there is an error
// loading the item or it is not downloaded.
func (s source) Lookup(lang lib.Lang, path string) (lib.Item, error) {
	lcode := lang.Code()

	var item lib.Item

	var previousAttemptPath string
	for item == nil || item.Path() != path {
		// Traverse the path backwards to find an ancestor available
		ppath := path
		item = nil
		for {
			if i, ok := itemsByLangAndPath[ref{lcode, ppath}]; ok {
				item = i
				break
			}
			ppath = ppath[:strings.LastIndex(ppath, "/")]
			if ppath == "" {
				ppath = "/"
			}
		}

		if ppath == previousAttemptPath {
			return nil, lib.ErrNotFound
		}

		// Work back down the ancestory to the path
		err := lib.DlAndOpen(item)
		if err != nil {
			return nil, err
		}
		previousAttemptPath = ppath
	}

	// c, err := lang.Catalog()
	// if err != nil {
	// 	return nil, err
	// }
	// if r.Path == "/" {
	// 	return c, nil
	// }
	// if folder, ok := c.foldersByPath[r.Path]; ok {
	// 	return folder, nil
	// }
	// sections := strings.Split(r.Path, "/")
	// if sections[0] != "" {
	// 	return nil, fmt.Errorf("Invalid path \"%v\", must start with '/'", r.Path)
	// }
	// for i := 2; i <= len(sections); i++ {
	// 	temppath := strings.Join(sections[0:i], "/")
	// 	if book, ok := c.booksByPath[temppath]; ok {
	// 		if r.Path == book.Path() {
	// 			return book, nil
	// 		}
	// 		node, err := book.lookupPath(r.Path)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("Path %v not found - %v", r.Path, err.Error())
	// 		}
	// 		return node, err
	// 	}
	// }
	if item != nil {
		return item, nil
	} else {
		return nil, lib.ErrNotFound
	}
}
