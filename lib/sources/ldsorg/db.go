package ldsorg

import (
	"database/sql"

	"io"
	"io/ioutil"
	"os"

	"github.com/duckbrain/ldss/lib"
	_ "github.com/mattn/go-sqlite3"
)

type ZBook struct {
	db            *sql.DB
	stmtChildren  *sql.Stmt
	stmtFootnotes *sql.Stmt
}

func NewZBook(r io.Reader) (*ZBook, error) {
	// Download and extract the book
	file, err := ioutil.TempFile(os.TempDir(), "ldss-*.sqlite3")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())
	if err := io.Copy(file, r); err != nil {
		return nil, err
	}
	if err := file.Close(); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", file.Name())
	if err != nil {
		return nil, err
	}

	z = &ZBook{db: db}

	const sqlQueryNode = `
		SELECT
			node.id,
			node.parent_id,
			node.uri,
			node.title,
			node.subtitle,
			node.section_name,
			node.short_title,
			node.content
		FROM node
		WHERE parent_id = ?`
	z.stmtChildren, err = db.Prepare(sqlQueryNode)
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
	z.stmtFootnotes, err = db.Prepare(sqlQueryRef)
	if err != nil {
		return nil, err
	}

	return z, nil
}

func (z *ZBook) Close() error {
	return z.db.Close()
}

func (z *ZBook) Children(id int64, c chan<- Node) error {
	defer close(c)
	rows, err := z.stmtChildren.Query(id)
	if err != nil {
		return err
	}
	for rows.Next() {
		n := Node{}
		err = rows.Scan(
			&n.ID,
			&n.ParentID,
			&n.Path,
			&n.Name,
			&n.Subtitle,
			&n.SectionName,
			&n.ShortTitle,
		)
		if err != nil {
			return err
		}
		c <- n
	}
	return nil
}

func (z *ZBook) Footnotes(id int64, c chan<- lib.Footnote) error {
	defer close(c)
	rows, err := z.stmtFootnotes.Query(id)
	if err != nil {
		return err
	}
	for rows.Next() {
		f := lib.Footnote{}
		err = rows.Scan(
			&f.Name,
			&f.LinkName,
			&f.Content,
		)
		if err != nil {
			return err
		}
		c <- f
	}
	return nil
}
