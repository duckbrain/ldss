package ldslib

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type BookLoader struct {
	content *LocalContent
	book    *Book
	db      *sql.DB
}

const sqlQueryNode = "SELECT id, doc_version, parent_id, title, uri, content FROM node "

func (l *BookLoader) NewBookLoader(book *Book) {
	l.book = book
}

func (l *BookLoader) populateIfNeeded() {
	if l.db == nil {
		db, err := sql.Open("sqlite3", l.content.GetBookPath(l.book))
		if err != nil {
			panic(err)
		}
		l.db = db
	}
}

func (l *BookLoader) nodesFromRows(rows *sql.Rows) []Node {
	nodes := make([]Node, 0)
	for rows.Next() {
		node := Node{}

		nodes = append(nodes, node)
	}
	return nodes
}

func (l *BookLoader) Close() {
	l.db.Close()
}

func (l *BookLoader) GetIndex() []Node {
	l.populateIfNeeded()
	rows, err := l.db.Query(sqlQueryNode + "WHERE parent_id = 0")
	if err != nil {
		panic(err)
	}
	return l.nodesFromRows(rows)
}
