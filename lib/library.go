package lib

var catalogsByLanguageId map[int]*Catalog
var booksByLangBookId map[langBookID]*Book

type langBookID struct {
	langID int
	bookID int
}

func init() {
	catalogsByLanguageId = make(map[int]*Catalog)
	booksByLangBookId = make(map[langBookID]*Book)
}

// Calls the open function and downloads any missing content if has not
// been downloaded yet. If it does download something, it will retry the
// open function repeatedly until it gets no download errors (or repeated
// ones). It provides information about these downloads and the final
// result through a series of messages from the returned channel.
func AutoDownload(open func() (interface{}, error)) <-chan Message {
	c := make(chan Message)
	go func() {
		item, err := open()
		defer close(c)
		var dlErr, preDlErr NotDownloadedErr
		dlErr, ok := err.(NotDownloadedErr)
		for ok {
			if dlErr == preDlErr {
				return
			}
			c <- MessageDownload{dlErr}
			err = dlErr.Download()
			if err != nil {
				c <- MessageError{err}
				return
			}
			item, err = open()
			preDlErr = dlErr
			dlErr, ok = err.(NotDownloadedErr)
		}

		if err == nil {
			c <- MessageDone{item}
		} else {
			c <- MessageError{err}
		}
	}()
	return c
}

// Gets the next or previous sibling of an item using it's
// interface functions. Used to implement Next and Previous
// on the Item interface. Only pass -1 or 1 to direction.
func genericNextPrevious(item Item, direction int) Item {
	parent := item.Parent()
	if parent == nil {
		return nil
	}
	siblings, err := parent.Children()
	if err != nil {
		panic(err)
	}

	getSideSibs := func() []Item {
		parentSibling := genericNextPrevious(parent, direction)
		if parentSibling == nil {
			return nil
		}
		if sideSibs, err := parentSibling.Children(); err == nil && len(sideSibs) > 0 {
			return sideSibs
		}
		return nil
	}

	for i, sibling := range siblings {
		if sibling.Path() == item.Path() {
			if i+direction < 0 {
				//Get last child of parent's sibling
				if sideSibs := getSideSibs(); sideSibs != nil {
					return sideSibs[len(sideSibs)-1]
				}
				return nil
			}
			if i+direction >= len(siblings) {
				//Get first child of parent's sibling
				if sideSibs := getSideSibs(); sideSibs != nil {
					return sideSibs[0]
				}
				return nil
			}
			return siblings[i+direction]
		}
	}
	panic("Item not found as child's parent.")
}

// Uses AutoDownload to lookup a string. Uses reference lookups or paths
// to find the item.
func Lookup(lang *Language, q string) <-chan Message {
	return AutoDownload(func() (interface{}, error) {
		ref, err := lang.ref()
		if err == nil {
			q, err = ref.lookup(q)
			if err != nil {
				return nil, err
			}
		}
		catalog, err := lang.Catalog()
		if err != nil {
			return nil, err
		}
		item, err := catalog.LookupPath(q)
		if err != nil {
			return nil, err
		}
		return item, nil
	})
}
