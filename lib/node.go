package lib

import (
	"fmt"
)

type Node struct {
	id         int
	name       string
	path       string
	Book       *Book
	hasContent bool
	childCount int
	parentId   int
}

func (n *Node) Name() string {
	return n.name
}

func (n *Node) String() string {
	return fmt.Sprintf("%v {%v}", n.name, n.path)
}

func (n *Node) Path() string {
	return n.path
}

func (n *Node) Language() *Language {
	return n.Book.Language()
}

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

func (n *Node) Content() (*Content, error) {
	rawContent, err := n.Book.nodeContent(n)
	return &Content{rawHTML: rawContent}, err
}

func (n *Node) Parent() Item {
	if n.parentId == 0 {
		return n.Book
	} else {
		node, _ := n.Book.lookupId(n.parentId)
		return node
	}
}

func (n *Node) Next() Item {
	return genericNextPrevious(n, 1)
}

func (n *Node) Previous() Item {
	return genericNextPrevious(n, -1)
}
