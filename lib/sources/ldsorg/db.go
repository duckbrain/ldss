package ldsorg

import (
	"context"
	"database/sql"

	"io"
	"io/ioutil"
	"os"

	"github.com/duckbrain/ldss/lib"
	_ "github.com/mattn/go-sqlite3"
)

type ZBook struct {
	filename      string
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
	if _, err := io.Copy(file, r); err != nil {
		return nil, err
	}
	if err := file.Close(); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", file.Name())
	if err != nil {
		return nil, err
	}

	z := &ZBook{db: db, filename: file.Name()}

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
	err := z.db.Close()
	if err != nil {
		return err
	}
	return os.Remove(z.filename)
}

func (z *ZBook) Children(ctx context.Context, id int64) ([]Node, error) {
	rows, err := z.stmtChildren.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}
	nodes := make([]Node, 0)
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
			return nil, err
		}
		nodes = append(nodes, n)
	}
	return nodes, nil
}

func (z *ZBook) Footnotes(ctx context.Context, id int64) ([]lib.Footnote, error) {
	rows, err := z.stmtFootnotes.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}
	footnotes := make([]lib.Footnote, 0)
	for rows.Next() {
		f := lib.Footnote{}
		err = rows.Scan(
			&f.Name,
			&f.LinkName,
			&f.Content,
		)
		if err != nil {
			return nil, err
		}
		footnotes = append(footnotes, f)
	}
	return footnotes, nil
}
