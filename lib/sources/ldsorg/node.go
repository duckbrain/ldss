package ldsorg

import (
	"fmt"

	"github.com/duckbrain/ldss/lib"
)

// Represents a node in a Book
type Node struct {
	id          int
	name        string
	path        string
	Book        *book
	hasContent  bool
	childCount  int
	parentId    int
	Subtitle    string
	SectionName *string
	ShortTitle  *string
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
func (n *Node) Language() Lang {
	return n.Book.Language()
}

// The children of the node, will all be Nodes
func (n *Node) Children() ([]Item, error) {
	nodes, err := n.Book.nodeChildren(n)
	if err != nil {
		return nil, err
	}
	items := make([]Item, len(nodes))
	for i, n := range nodes {
		if subnodes, err := n.Children(); err == nil && len(subnodes) == 1 {
			items[i] = subnodes[0]
		} else {
			items[i] = n
		}
	}
	return items, nil
}

func (n *Node) Footnotes(verses []int) ([]Footnote, error) {
	return n.Book.nodeFootnotes(n, verses)
}

// Returns the content of the Node, to use as HTML or Parse
func (n *Node) Content() (lib.Content, error) {
	rawContent, err := n.Book.nodeContent(n)
	return lib.Content(rawContent), err
}

// Parent node or book
func (n *Node) Parent() (parent Item) {
	if n.parentId == 0 {
		parent = n.Book
	} else {
		node, _ := n.Book.lookupId(n.parentId)
		parent = node
	}
	if siblings, err := parent.Children(); err == nil && len(siblings) == 1 {
		parent = parent.Parent()
	}
	return
}

// Next sibling node
func (n *Node) Next() Item {
	return genericNextPrevious(n, 1)
}

// Preivous sibling node
func (n *Node) Previous() Item {
	return genericNextPrevious(n, -1)
}
