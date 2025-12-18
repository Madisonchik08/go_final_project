package main

import (
	"log"
	"os"

	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

// resolveDBFile returns database file path from environment variable or default.
// resolveDBFile возвращает путь к файлу базы данных из переменной окружения или значение по умолчанию.
func resolveDBFile() string {
	if dbFile := os.Getenv("TODO_DBFILE"); dbFile != "" {
		return dbFile
	}
	return "scheduler.db"
}

// main initializes database and starts HTTP server.
// main инициализирует базу данных и запускает HTTP сервер.
func main() {
	if err := db.Init(resolveDBFile()); err != nil {
		log.Fatal(err)
	}
	defer db.DB.Close()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
