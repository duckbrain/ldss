package ldsorg

import (
	"fmt"

	"github.com/duckbrain/ldss/lib"
)

// Represents a node in a Book
type node struct {
	id          int64
	name        string
	path        string
	hasContent  bool
	childCount  int
	parentId    int
	parent      lib.Item
	book        *book
	children    []lib.Item
	subtitle    string
	sectionName string
	shortTitle  string
	content     lib.Content
	footnotes   []lib.Footnote
}

func (n *node) Open() error {
	// TODO: Make sure children, content, parent, etc. is populated
	path := bookPath(n.book)
	l, err := opendb(path)
	if err != nil {
		return err
	}
	defer l.Close()

	// Populate Children
	nodes, err := l.childrenByParentID(n.id, n, n.book)
	if err != nil {
		return err
	}
	n.children = nodes

	langCode := n.Lang().Code()
	for _, node := range nodes {
		itemsByLangAndPath[ref{langCode, node.Path()}] = node
	}

	// Populate content
	content, err := l.contentByNodeID(n.id)
	if err != nil {
		return err
	}
	n.content = lib.Content(content)

	// Populate Footnotes
	n.footnotes, err = l.footnotesByNode(n, nil)
	if err != nil {
		return err
	}

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
	if len(verses) == 0 {
		return n.footnotes, nil
	}

	path := bookPath(n.book)
	l, err := opendb(path)
	if err != nil {
		return nil, err
	}
	defer l.Close()

	return l.footnotesByNode(n, verses)
}

// Returns the content of the Node, to use as HTML or Parse
func (n *node) Content() (lib.Content, error) {
	return n.content, nil
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
