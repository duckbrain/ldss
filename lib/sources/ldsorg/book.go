package ldsorg

import (
	"database/sql"
	"fmt"
	"html/template"
	"os"
	"sort"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/duckbrain/ldss/lib"
	_ "github.com/mattn/go-sqlite3"
)

type Footnote = lib.Footnote

// Represents a book in the catalog or one of it's folders and
// provides a way to access the nodes in it's database if
// downloaded.
// Used for parsing books in the catalog's JSON file
type book struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	GlURI   string `json:"gl_uri"`
	catalog *catalog
	parent  Item
}

type bookDBConnection struct {
	db                                                        *sql.DB
	stmtChildren, stmtUri, stmtId, stmtContent, stmtFootnotes *sql.Stmt
}

func (c bookDBConnection) Close() error {
	return c.db.Close()
}

func init() {
	var a []*book
	bookOpened = make(chan *book)
	go func() {
	WAIT:
		for book := range bookOpened {
			if a == nil {
				a = make([]*book, 0, BookConnectionLimit)
			}

			for i, x := range a {
				if x == book {
					for j := i; j > 0; j-- {
						a[j] = a[j-1]
					}
					a[0] = x
					continue WAIT
				}
			}

			if len(a) < cap(a) {
				a = append(a, book)
			} else {
				x := a[len(a)-1]
				for i := len(a) - 1; i > 0; i-- {
					a[i] = a[i-1]
				}
				a[0] = book
				x.dbCache.Close()
			}
		}
	}()
}

var BookConnectionLimit = 20
var bookOpened chan *book

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
	FROM node
`

const sqlQueryRef = `
	SELECT
		ref.ref_name,
		ref.link_name,
		ref.ref
	FROM ref
	WHERE
		ref.node_id = ?
`

type sqlScanner interface {
	Scan(dest ...interface{}) error
}

func (n *Node) scan(s sqlScanner) error {
	return s.Scan(&n.id, &n.name, &n.path, &n.parentId,
		&n.Subtitle, &n.SectionName, &n.ShortTitle,
		&n.hasContent, &n.childCount)
}

func newBook(base *jsonBook, catalog *catalog, parent Item) *Book {
	b := &Book{}
	b.base = base
	b.catalog = catalog
	b.parent = parent

	b.dbCache.construct = func() (interface{}, error) {
		var l bookDBConnection
		path := bookPath(b)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, err
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
		l.stmtFootnotes, err = db.Prepare(sqlQueryRef)
		if err != nil {
			return nil, err
		}
		return &l, nil
	}
	return b
}

// A short human-readable representation of the book, mostly useful for debugging.
func (b *book) String() string {
	return fmt.Sprintf("%v {%v}", b.base.Name, b.base.GlURI)
}

// An ID for the book. This is unique to this book within it's language.
func (b *book) ID() int {
	return b.base.ID
}

// The name of this book.
func (b *book) Name() string {
	return b.base.Name
}

// The URL the database of this book is located at online.
func (b *book) URL() string {
	return b.base.URL
}

// The Gospel Library Path of this book, unique within it's language
func (b *book) Path() string {
	return b.base.GlURI
}

// The language this book is in
func (b *book) Lang() Lang {
	return b.catalog.language
}

// Children in this book. This is identical to the Index function, but returns
// the index as a []Item
func (b *book) Children() ([]Item, error) {
	nodes, err := b.Index()
	if err != nil {
		return nil, err
	}
	items := make([]Item, len(nodes))
	for i, n := range nodes {
		if subnodes, err := n.Children(); err == nil && len(subnodes) == 1 {
			items[i] = subnodes[0]
		} else {
			items[i] = n
		}
	}
	return items, nil
}

// Parent Folder or Catalog of this book
func (b *book) Parent() Item {
	return b.parent
}

// Next book in the Folder
func (b *book) Next() Item {
	return genericNextPrevious(b, 1)
}

// Previous book in the Folder
func (b *book) Previous() Item {
	return genericNextPrevious(b, -1)
}

// The SQLite database connector with some prepared statements. Cached for subsequent uses.
func (b *book) db() (*bookDBConnection, error) {
	db, err := b.dbCache.get()
	if err != nil {
		return nil, err
	}
	bookOpened <- b
	return db.(*bookDBConnection), nil
}

// Returns the Nodes at the root of this book.
func (l *book) Index() ([]*Node, error) {
	return l.nodeChildren(nil)
}

// Returns the Nodes that are children of the passed node. If the passed Node
// is nil, it will return the Index
func (b *book) nodeChildren(parent *Node) ([]*Node, error) {
	l, err := b.db()
	parentId := 0
	if err != nil {
		return nil, err
	}
	if parent != nil {
		parentId = parent.id
	}
	rows, err := l.stmtChildren.Query(parentId)
	if err != nil {
		return nil, err
	}
	nodes := make([]*Node, 0)
	for rows.Next() {
		node := &Node{Book: b}
		err := node.scan(rows)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (b *book) lookupId(id int) (*Node, error) {
	node := &Node{book: b}
	l, err := b.db()
	if err != nil {
		return node, err
	}
	err = node.scan(l.stmtId.QueryRow(id))
	return node, err

}

func (b *book) lookupPath(uri string) (*Node, error) {
	node := &Node{book: b}
	l, err := b.db()
	if err != nil {
		return node, err
	}
	err = node.scan(l.stmtUri.QueryRow(uri))
	if err != nil {
		return nil, fmt.Errorf("Path %v not found", uri)
	}
	return node, err
}

func (b *book) nodeContent(node *Node) (string, error) {
	l, err := b.db()
	if err != nil {
		return "", err
	}
	var content string
	err = l.stmtContent.QueryRow(node.id).Scan(&content)
	return content, err
}

func (b *book) nodeFootnotes(node *Node, verses []int) ([]Footnote, error) {
	l, err := b.db()
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, fmt.Errorf("Cannot use nil Node in nodeFootnotes()")
	}
	rows, err := l.stmtFootnotes.Query(node.id)
	if err != nil {
		return nil, err
	}
	refs := make([]Footnote, 0)
	for rows.Next() {
		ref := Footnote{}
		var content string
		err = rows.Scan(&ref.Name, &ref.LinkName, &content)
		ref.Content = template.HTML(content)
		ref.Item = node
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
