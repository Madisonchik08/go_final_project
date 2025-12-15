package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const (
	defaultPort = 7540
	webDir      = "web"
)

func resolvePort() int {
	envPort := os.Getenv("TODO_PORT")
	if envPort == "" {
		return defaultPort
	}
	port, err := strconv.Atoi(envPort)
	if err != nil || port <= 0 {
		return defaultPort
	}
	return port
}

func resolveWebDir() string {
	dir, err := filepath.Abs(webDir)
	if err != nil {
		return webDir
	}
	return dir
}

func main() {
	dir := resolveWebDir()

	http.Handle("/", http.FileServer(http.Dir(dir)))

	addr := fmt.Sprintf(":%d", resolvePort())
	log.Printf("Starting file server for %s on http://localhost%s", dir, addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
