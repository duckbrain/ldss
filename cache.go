package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type CacheConnection struct {
	db *sql.DB
}

func NewCacheConnection() *CacheConnection {
	c := new(CacheConnection)
	return c
}

func (c *CacheConnection) Open(filename string) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic(err)
	}
	c.db = db
}

func (c *CacheConnection) SaveLanguages(languages []Language) {
	c.db.Exec(`
		CREATE TABLE IF NOT EXISTS language 
		(
			id int, 
			name string, 
			eng_name string, 
			code string, 
			gl_code
		);
		DELETE FROM language;
		`)

	tx, err := c.db.Begin()
	if err != nil {
		panic(err)
	}
	stmt, err := tx.Prepare(`
		INSERT INTO language
		(
			id, name, eng_name, code, gl_code
		) 
		VALUES (?, ?, ?, ?, ?);
		`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	for _, lang := range languages {
		_, err = stmt.Exec(lang.ID, lang.Name, lang.EnglishName, lang.Code, lang.GlCode)
		if err != nil {
			panic(err)
		}
	}
	tx.Commit()
}

func (c *CacheConnection) GetByUnknown(id string) *Language {
	//var l Language
	////err := c.db.QueryRow("SELECT id, name, eng_name, code, gl_code FROM languages;").Scan(&l.ID, &l.Name, &l.EnglishName, &l.Code, &l.GlCode)
	return nil

}

func (c *CacheConnection) GetAll() []Language {
	languages := make([]Language, 0)
	rows, err := c.db.Query("SELECT id, name, eng_name, code, gl_code FROM language;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var l Language
		rows.Scan(&l.ID, &l.Name, &l.EnglishName, &l.Code, &l.GlCode)
		languages = append(languages, l)
	}
	return languages
}
