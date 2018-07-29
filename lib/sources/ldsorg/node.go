package ldsorg

import (
	"fmt"

	"github.com/duckbrain/ldss/lib"
)

// Represents a node in a Book
type Node struct {
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

func (n *Node) Open() error {
	// TODO: Make sure children, content, parent, etc. is populated
	return nil
}

// Name of the node
func (n *Node) Name() string {
	return n.name
}

// A short human-readable representation of the node, mostly useful for debugging.
func (n *Node) String() string {
	return fmt.Sprintf("%v {%v}", n.name, n.path)
}

// The full Gospel Library path of the node
func (n *Node) Path() string {
	return n.path
}

// The language the node is in.
func (n *Node) Lang() Lang {
	return n.parent.Lang()
}

// The children of the node, will all be Nodes
func (n *Node) Children() []Item {
	return n.children
}

func (n *Node) Footnotes(verses []int) ([]lib.Footnote, error) {
	return n.conn.footnotesByNode(n, verses)
}

// Returns the content of the Node, to use as HTML or Parse
func (n *Node) Content() (lib.Content, error) {
	rawContent, err := n.conn.contentByNodeID(n.id)
	return lib.Content(rawContent), err
}

// Parent node or book
func (n *Node) Parent() lib.Item {
	return n.parent
}

// Next sibling node
func (n *Node) Next() Item {
	return lib.GenericNextPrevious(n, 1)
}

// Preivous sibling node
func (n *Node) Prev() Item {
	return lib.GenericNextPrevious(n, -1)
}
