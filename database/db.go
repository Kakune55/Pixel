package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Initdb() {
	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建表的SQL语句
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS mytable (
		link TEXT PRIMARY KEY,
		md5 TEXT NOT NULL,
		ext TEXT NOT NULL
	);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func NewFile(linkID string, md5 string, ext string) {
	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// SQL语句
	SQL := `
	INSERT INTO "main"."mytable" ("link", "md5" ,"ext") VALUES (?, ? , ?)
	`

	stmt, err := db.Prepare(SQL)
	if err != nil {
		log.Fatal(err)
	}

 	_, err = stmt.Exec(linkID,md5,ext) //插入记录
	if err != nil {
		log.Fatal(err)
	}
}

func GetFileName(linkID string) string {
	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// SQL语句
	SQL := `
	SELECT * FROM "mytable" WHERE link = ?
	`

	row := db.QueryRow(SQL,linkID)
	if err != nil {
		log.Fatal(err)
	}

	// 扫描查询结果
	var md5 string
	var linkIDx string
	var ext string
	err = row.Scan(&linkIDx,&md5,&ext)
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		} else {
			log.Fatal(err)
		}
		return ""
	}
	return md5+ext
}

func QueryId() ([]string, error) {
	var result []string

	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT link FROM mytable")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		result = append(result, name)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}