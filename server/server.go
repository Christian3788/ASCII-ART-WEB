package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/christianotieno/ascii-art-web/handlers"
)

const (
	templateDir = "templates"
	staticDir   = "static"
)

func Run() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(".", staticDir)))))
	handlers.RegisterRoutes(mux)

	fmt.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
