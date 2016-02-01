package lib

type Node struct {
	id         int
	name       string
	glURI      string
	Book       *Book
	hasContent bool
	childCount int
	parent     Item
}

func (n Node) Name() string {
	return n.name
}

func (n Node) String() string {
	return n.name
}

func (n Node) Path() string {
	return n.glURI
}

func (n Node) Language() *Language {
	return n.Book.Language()
}

func (n Node) Children() ([]Item, error) {
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

func (n Node) Content() (*Content, error) {
	rawContent, err := n.Book.nodeContent(n)
	return &Content{rawHTML: rawContent}, err
}

func (n Node) Parent() Item {
	return n.parent
}
