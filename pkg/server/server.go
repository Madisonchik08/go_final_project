package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"go_final_project/pkg/api"
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

// Start launches HTTP server that serves files from web directory.
func Start() error {
	dir := resolveWebDir()

	mux := http.NewServeMux()
	api.Init(mux)
	mux.Handle("/", http.FileServer(http.Dir(dir)))

	addr := fmt.Sprintf(":%d", resolvePort())
	log.Printf("Starting file server for %s on http://localhost%s", dir, addr)

	return http.ListenAndServe(addr, mux)
}
