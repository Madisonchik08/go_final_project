package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// schema defines the database table structure.
// schema определяет структуру таблицы базы данных.
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

// DB keeps opened database connection.
// DB хранит открытое соединение с базой данных.
var DB *sql.DB

// Init opens SQLite database at dbFile and ensures schema exists.
// Init открывает базу данных SQLite по пути dbFile и создает схему, если она не существует.
func Init(dbFile string) error {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	DB = db

	_, err = DB.Exec(schema)
	return err
}
