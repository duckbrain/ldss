package ldslib

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type bookParser struct {
	source                             Source
	book                               *Book
	db                                 *sql.DB
	stmtChildren, stmtUri, stmtContent *sql.Stmt
}

const sqlQueryNode = `
	SELECT
		node.id,
		node.title,
		node.uri,
		CASE WHEN node.content IS NULL THEN 0 ELSE 1 END,
		(SELECT COUNT(*) FROM node subnode WHERE subnode.id = node.id) node_count
	FROM node
`

func newBookParser(book *Book, source Source) *bookParser {
	return &bookParser{source: source, book: book}
}

func (l *bookParser) populate() error {
	if l.db == nil {
		path := l.source.BookPath(l.book)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Println("File not found")
			return &NotDownloadedBookErr{err, l.book}
		}
		db, err := sql.Open("sqlite3", path)
		if err != nil {
			return err
		}
		l.db = db
		l.stmtChildren, err = db.Prepare(sqlQueryNode + " WHERE parent_id = ?")
		if err != nil {
			return err
		}
		l.stmtUri, err = db.Prepare(sqlQueryNode + " WHERE uri = ?")
		if err != nil {
			return err
		}
		l.stmtContent, err = db.Prepare("SELECT content FROM node WHERE id = ?")
		if err != nil {
			return err
		}
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM node;").Scan(&count)
		if err != nil {
			return &NotDownloadedBookErr{err, l.book}
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

func (l *bookParser) Children(parent Node) ([]Node, error) {
	if err := l.populate(); err != nil {
		return nil, err
	}
	rows, err := l.stmtChildren.Query(parent.ID)
	if err != nil {
		return nil, err
	}
	nodes := make([]Node, 0)
	for rows.Next() {
		node := Node{Book: l.book}
		if node.ID > 0 {
			node.parent = node
		}
		err := rows.Scan(&node.ID, &node.Name, &node.GlURI, &node.HasContent, &node.ChildCount)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (l *bookParser) GlUri(uri string) (Node, error) {
	node := Node{Book: l.book}
	if err := l.populate(); err != nil {
		return node, err
	}
	err := l.stmtUri.QueryRow(uri).Scan(&node.ID, &node.Name, &node.GlURI, &node.HasContent, &node.ChildCount)
	return node, err
}

func (l *bookParser) Content(node Node) (string, error) {
	if err := l.populate(); err != nil {
		return "", err
	}
	var content string
	err := l.stmtContent.QueryRow(node.ID).Scan(&content)
	return content, err
}
