package database

import (
	"os"
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(128) NOT NULL,
	comment TEXT,
	repeat VARCHAR(64) NOT NULL
);

CREATE INDEX IF NOT EXISTS date_index ON scheduler (date)`

func Close() error {

	if DB != nil {
		return DB.Close()
	}
	return nil
}

func Init() error {

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	DB, err = sql.Open("sqlite", "./" + dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}
