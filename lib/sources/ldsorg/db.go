package ldsorg

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"html/template"
	"os"
	"sort"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/duckbrain/ldss/lib"
)

var MaxDatabaseQueueLength = 50

var conns = make(map[string]*sqlconn)
var connsLock = sync.Mutex{}
var connsQueue = []*sqlconn{}

type sqlconn struct {
	path          string
	db            *sql.DB
	stmtChildren  *sql.Stmt
	stmtUri       *sql.Stmt
	stmtId        *sql.Stmt
	stmtContent   *sql.Stmt
	stmtFootnotes *sql.Stmt
}

func opendb(path string) (*sqlconn, error) {
	connsLock.Lock()
	defer connsLock.Unlock()
	l, ok := conns[path]
	if ok {
		return l, nil
	}

	l = &sqlconn{}
	l.path = path
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.Wrap(err, "Failed to find the database file")
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM node;").Scan(&count)
	if err != nil {
		return nil, err
	}
	l.db = db
	const sqlQueryNode = `
		SELECT
			node.id,
			node.title,
			node.uri,
			node.parent_id,
			node.subtitle,
			node.section_name,
			node.short_title,
			CASE WHEN node.content IS NULL THEN 0 ELSE 1 END,
			(SELECT COUNT(*) FROM node subnode WHERE subnode.id = node.id) node_count
		FROM node`
	l.stmtChildren, err = db.Prepare(sqlQueryNode + " WHERE parent_id = ?")
	if err != nil {
		return nil, err
	}
	l.stmtUri, err = db.Prepare(sqlQueryNode + " WHERE uri = ?")
	if err != nil {
		return nil, err
	}
	l.stmtId, err = db.Prepare(sqlQueryNode + " WHERE id = ?")
	if err != nil {
		return nil, err
	}
	l.stmtContent, err = db.Prepare("SELECT content FROM node WHERE id = ?")
	if err != nil {
		return nil, err
	}

	const sqlQueryRef = `
		SELECT
			ref.ref_name,
			ref.link_name,
			ref.ref
		FROM ref
		WHERE
			ref.node_id = ?`
	l.stmtFootnotes, err = db.Prepare(sqlQueryRef)
	if err != nil {
		return nil, err
	}

	conns[path] = l
	connsQueue = append(connsQueue, l)
	if len(connsQueue) > MaxDatabaseQueueLength {
		x := connsQueue[0]
		connsQueue = connsQueue[1:]
		delete(conns, x.path)
		err = x.Close()
	}

	return l, err
}

func (l *sqlconn) Close() error {
	// Indicate that the connection is done.
	// Brainstoriming for connection pool: Create a map[*sqlconn]int with counts of opens without closes. A goroutine gets started after the max is hit and garbage collects the zeros, until it goes down to the allowed number.
	return nil
}

func (l *sqlconn) childrenByParentID(id int64, parent lib.Item, book *book) ([]lib.Item, error) {
	rows, err := l.stmtChildren.Query(id)
	if err != nil {
		return nil, err
	}
	nodes := make([]lib.Item, 0)
	for rows.Next() {
		n := &node{
			parent: parent,
			book:   book,
		}
		err := n.scan(rows)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}
	return nodes, nil

}

func (l *sqlconn) nodeByID(id int64, parent lib.Item, book *book) (*node, error) {
	row := l.stmtId.QueryRow(id)
	n := &node{
		parent: parent,
		book:   book,
	}
	err := n.scan(row)
	return n, err
}

func (l *sqlconn) nodeByGlURI(uri string, parent lib.Item, book *book) (*node, error) {
	row := l.stmtUri.QueryRow(uri)
	n := &node{
		parent: parent,
		book:   book,
	}
	err := n.scan(row)
	return n, err
}

func (l *sqlconn) contentByNodeID(id int64) (string, error) {
	var content sql.NullString
	err := l.stmtContent.QueryRow(id).Scan(&content)
	return content.String, err
}

func (l *sqlconn) footnotesByNode(n *node, verses []int) ([]lib.Footnote, error) {
	rows, err := l.stmtFootnotes.Query(n.id)
	if err != nil {
		return nil, err
	}
	refs := make([]lib.Footnote, 0)
	for rows.Next() {
		ref := lib.Footnote{}
		var content string
		err = rows.Scan(&ref.Name, &ref.LinkName, &content)
		ref.Content = template.HTML(content)
		ref.Item = n
		if err != nil {
			return nil, err
		}
		if len(verses) > 0 {
			verseNumString := ref.Name
			if char, length := utf8.DecodeLastRuneInString(verseNumString); unicode.IsDigit(char) {
				verseNumString = verseNumString[:len(verseNumString)-length-1]
			}
			if verseNum, err := strconv.Atoi(verseNumString); err != nil {
				if sort.SearchInts(verses, verseNum) > -1 {
					refs = append(refs, ref)
				}
			}
		} else {
			refs = append(refs, ref)
		}
	}
	return refs, nil
}

type sqlScanner interface {
	Scan(dest ...interface{}) error
}

func (n *node) scan(s sqlScanner) error {
	var subtitle, sectionName, shortTitle sql.NullString
	err := s.Scan(&n.id, &n.name, &n.path, &n.parentId,
		&subtitle, &sectionName, &shortTitle,
		&n.hasContent, &n.childCount)
	n.subtitle = subtitle.String
	n.sectionName = sectionName.String
	n.shortTitle = shortTitle.String
	return err
}

var BookConnectionLimit = 20
var connectionOpen = make(chan *sqlconn)

// func init() {
// 	go func() {
// 	WAIT:
// 		for conn := range connectionOpen {
// 			if a == nil {
// 				a = make([]*book, 0, BookConnectionLimit)
// 			}
//
// 			for i, x := range a {
// 				if x == book {
// 					for j := i; j > 0; j-- {
// 						a[j] = a[j-1]
// 					}
// 					a[0] = x
// 					continue WAIT
// 				}
// 			}
//
// 			if len(a) < cap(a) {
// 				a = append(a, book)
// 			} else {
// 				x := a[len(a)-1]
// 				for i := len(a) - 1; i > 0; i-- {
// 					a[i] = a[i-1]
// 				}
// 				a[0] = book
// 				x.dbCache.Close()
// 			}
// 		}
// 	}()
// }
