package ldslib

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type bookParser struct {
	source Source
	book   *Book
	db     *sql.DB
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
	}
	return nil
}

func (l *bookParser) nodes(query string) ([]Node, error) {
	if err := l.populate(); err != nil {
		return nil, err
	}
	rows, err := l.db.Query(sqlQueryNode + query)
	if err != nil {
		return nil, err
	}
	nodes := make([]Node, 0)
	for rows.Next() {
		node := Node{}
		if err := rows.Scan(node); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (l *bookParser) Close() {
	l.db.Close()
}

func (l *bookParser) Index() ([]Node, error) {
	return l.nodes("WHERE parent_id = 0")
}

func (l *bookParser) Children(node Node) ([]Node, error) {
	return l.nodes(fmt.Sprintf("WHERE parent_id = %v", node.ID))
}
