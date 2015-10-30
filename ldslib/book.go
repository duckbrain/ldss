package ldslib

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type bookParser struct {
	source Source
	book   *Book
	db     *sql.DB
}

const sqlQueryNode = "SELECT id, doc_version, parent_id, title, uri, content FROM node "

func (l *bookParser) NewBookLoader(book *Book) {
	l.book = book
}

func (l *bookParser) populateIfNeeded() {
	if l.db == nil {
		db, err := sql.Open("sqlite3", l.source.BookPath(l.book))
		if err != nil {
			panic(err)
		}
		l.db = db
	}
}

func (l *bookParser) nodesFromRows(rows *sql.Rows) []Node {
	nodes := make([]Node, 0)
	for rows.Next() {
		node := Node{}

		nodes = append(nodes, node)
	}
	return nodes
}

func (l *bookParser) Close() {
	l.db.Close()
}

func (l *bookParser) GetIndex() []Node {
	l.populateIfNeeded()
	rows, err := l.db.Query(sqlQueryNode + "WHERE parent_id = 0")
	if err != nil {
		panic(err)
	}
	return l.nodesFromRows(rows)
}
