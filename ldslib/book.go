package ldslib

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type bookParser struct {
	source Source
	book   *Book
	db     *sql.DB
	stmtChildren, stmtUri *sql.Stmt
}

type nodeInfo struct {
	id int
	docVersion *string
	parentId int
	title string
	uri string
	content string
}

const sqlQueryNode = "SELECT id, doc_version, parent_id, title, uri, content FROM node "

func newBookParser(book *Book, source Source) *bookParser {
	l := bookParser{}
	l.source = source
	l.book = book
	return &l
}

func (l *bookParser) populate() error {
	if l.db == nil {
		db, err := sql.Open("sqlite3", l.source.BookPath(l.book))
		if err != nil {
			return err
		}
		l.db = db
		l.stmtChildren, err = db.Prepare("SELECT id, title, uri FROM node WHERE parent_id = ?")
		l.stmtUri, err = db.Prepare("SELECT id, title, uri FROM node WHERE uri = ?")
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *bookParser) Close() {
	l.db.Close()
}

func (l *bookParser) Index() ([]Node, error) {
	return l.Children(Node{})
}

func (l *bookParser) Children(node Node) ([]Node, error) {
	if err := l.populate(); err != nil {
		return nil, err
	}
	rows, err := l.stmtChildren.Query(node.ID)
	if err != nil {
		return nil, err
	}
	nodes := make([]Node, 0)
	for rows.Next() {
		node := Node{}
		if err := rows.Scan(&node.ID, &node.Name, &node.GlURI); err != nil {
			return nil, err
		}
		node.Book = l.book
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (l *bookParser) GlUri(uri string) (Node, error) {
	if err := l.populate(); err != nil {
		return Node{}, err
	}
	node := Node{}
	err := l.stmtUri.QueryRow(uri).Scan(&node.ID, &node.Name, &node.GlURI)
	node.Book = l.book
	if err != nil {
		return Node{}, err
	}
	return node, nil
}

func (l *bookParser) Content(node Node) (string, error) {
	return "", nil
}
