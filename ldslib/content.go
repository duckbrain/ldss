package ldslib

type ContentParser struct {
	node    Node
	contentHtml *string
}

func (p *ContentParser) OriginalHTML() *string {
	return p.contentHtml
}


