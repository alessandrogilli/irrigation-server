package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite", "file:data/schedules.db?_foreign_keys=on")
	if err != nil {
		return err
	}

	createTables := `
	CREATE TABLE IF NOT EXISTS lines (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		is_on INTEGER DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS schedules (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		line_id INTEGER,
		description TEXT,
		on_time TEXT,
		off_time TEXT,
		scheduled INTEGER DEFAULT 0,
		in_progress INTEGER DEFAULT 0,
		FOREIGN KEY(line_id) REFERENCES lines(id) ON DELETE CASCADE
	);
	`

	_, err = DB.Exec(createTables)
	return err
}
