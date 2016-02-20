package lib

import (
	"fmt"
)

// Represents a node in a Book
type Node struct {
	id         int
	name       string
	path       string
	Book       *Book
	hasContent bool
	childCount int
	parentId   int
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
func (n *Node) Language() *Language {
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
		items[i] = n
	}
	return items, nil
}

// Returns the content of the Node, to use as HTML or Parse
func (n *Node) Content() (Content, error) {
	rawContent, err := n.Book.nodeContent(n)
	return Content(rawContent), err
}

// Parent node or book
func (n *Node) Parent() Item {
	if n.parentId == 0 {
		return n.Book
	} else {
		node, _ := n.Book.lookupId(n.parentId)
		return node
	}
}

// Next sibling node
func (n *Node) Next() Item {
	return genericNextPrevious(n, 1)
}

// Preivous sibling node
func (n *Node) Previous() Item {
	return genericNextPrevious(n, -1)
}
