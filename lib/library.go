package lib

import (
	"github.com/duckbrain/ldss/lib/dl"
)

// Calls the open function and downloads any missing content if has not
// been downloaded yet. If it does download something, it will retry the
// open function repeatedly until it gets no download errors (or repeated
// ones). It provides information about these downloads and the final
// result through a series of messages from the returned channel.
func AutoDownload(open func() (Item, error)) (Item, error) {
	item, err := open()
	var dlr, preDlr dl.Downloader
	dlr, ok := dl.IsNotDownloaded(err)
	for ok {
		if dlr == preDlr {
			return nil, err
		}
		err = dl.EnqueueAndWait(dlr)
		if err != nil {
			return nil, err
		}
		item, err = open()
		preDlr = dlr
		dlr, ok = dl.IsNotDownloaded(err)
	}
	return item, nil
}

// Gets the next or previous sibling of an item using it's
// interface functions. Used to implement Next and Previous
// on the Item interface. Only pass -1 or 1 to direction.
func GenericNextPrevious(item Item, direction int) Item {
	parent := item.Parent()
	if parent == nil {
		return nil
	}
	siblings := parent.Children()

	getSideSibs := func() []Item {
		parentSibling := GenericNextPrevious(parent, direction)
		if parentSibling == nil {
			return nil
		}
		sideSibs := parentSibling.Children()
		if len(sideSibs) > 0 {
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
