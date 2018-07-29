package ldsorg

import (
	"fmt"

	"github.com/duckbrain/ldss/lib"
)

// Represents a node in a Book
type node struct {
	id          int64
	conn        *sqlconn
	name        string
	path        string
	hasContent  bool
	childCount  int
	parentId    int
	parent      lib.Item
	children    []lib.Item
	subtitle    string
	sectionName *string
	shortTitle  *string
}

func (n *node) Open() error {
	// TODO: Make sure children, content, parent, etc. is populated
	return nil
}

// Name of the node
func (n *node) Name() string {
	return n.name
}

// A short human-readable representation of the node, mostly useful for debugging.
func (n *node) String() string {
	return fmt.Sprintf("%v {%v}", n.name, n.path)
}

// The full Gospel Library path of the node
func (n *node) Path() string {
	return n.path
}

// The language the node is in.
func (n *node) Lang() lib.Lang {
	return n.parent.Lang()
}

// The children of the node, will all be Nodes
func (n *node) Children() []lib.Item {
	return n.children
}

func (n *node) Footnotes(verses []int) ([]lib.Footnote, error) {
	return n.conn.footnotesByNode(n, verses)
}

// Returns the content of the Node, to use as HTML or Parse
func (n *node) Content() (lib.Content, error) {
	rawContent, err := n.conn.contentByNodeID(n.id)
	return lib.Content(rawContent), err
}

// Parent node or book
func (n *node) Parent() lib.Item {
	return n.parent
}

// Next sibling node
func (n *node) Next() lib.Item {
	return lib.GenericNextPrevious(n, 1)
}

// Preivous sibling node
func (n *node) Prev() lib.Item {
	return lib.GenericNextPrevious(n, -1)
}
