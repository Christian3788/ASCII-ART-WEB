package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Global template cache container
var templates *template.Template

// PageData passes structured content to templates/index.html
type PageData struct {
	Result string // Holds ASCII art
	Error  string // Holds error message to display
}

func main() {
	// Clean path resolution for stability
	tmplPath := filepath.Clean("templates/index.html")
	templates = template.Must(template.ParseFiles(tmplPath))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiHandler)

	log.Println("Server starting on http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// renderError injects errors safely into index.html layout with explicit HTTP status code
func renderError(w http.ResponseWriter, msg string, statusCode int) {
	w.WriteHeader(statusCode)
	data := PageData{Error: msg}
	if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("Failed to render error template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderError(w, "404 Not Found: Page does not exist", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		renderError(w, "405 Method Not Allowed: Method must be GET", http.StatusMethodNotAllowed)
		return
	}

	// Render home cleanly via cache container
	if err := templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Printf("Failed to render home template: %v", err)
		renderError(w, "500 Internal Server Error: Template error", http.StatusInternalServerError)
	}
}

func asciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderError(w, "405 Method Not Allowed: Method must be POST", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		renderError(w, "400 Bad Request: Failed to parse form data", http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")

	if text == "" || (banner != "standard" && banner != "shadow" && banner != "thinkertoy") {
		renderError(w, "400 Bad Request: Missing text input or invalid banner choice", http.StatusBadRequest)
		return
	}

	// Sanitize path using filepath.Clean to enforce directory safety boundaries
	filePath := filepath.Clean(fmt.Sprintf("banners/%s.txt", banner))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		renderError(w, "404 Not Found: Banner font style missing on server", http.StatusNotFound)
		return
	}

	result, err := generateAscii(text, banner)
	if err != nil {
		renderError(w, "500 Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render output text cleanly via template layout
	data := PageData{Result: result}
	if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("Failed to render response template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func generateAscii(text, banner string) (string, error) {
	text = strings.ReplaceAll(text, "\r\n", "\n")

	filePath := filepath.Clean(fmt.Sprintf("banners/%s.txt", banner))
	// Replaced deprecated ioutil.ReadFile with os.ReadFile
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("could not read banner file for %s", banner)
	}

	fileStr := strings.ReplaceAll(string(content), "\r\n", "\n")
	lines := strings.Split(fileStr, "\n")

	if len(lines) < 2 {
		return "", fmt.Errorf("banner file %s appears to be empty or misformatted", banner)
	}

	inputLines := strings.Split(text, "\n")
	var output strings.Builder

	for _, line := range inputLines {
		if line == "" {
			output.WriteString("\n")
			continue
		}

		const blockHeight = 9

		for i := 1; i <= 8; i++ {
			for _, runeVal := range line {
				if runeVal < 32 || runeVal > 126 {
					continue
				}

				startingLine := (int(runeVal)-32)*blockHeight + i

				if startingLine >= len(lines) {
					continue
				}

				output.WriteString(lines[startingLine])
			}
			output.WriteString("\n")
		}
	}

	return output.String(), nil
}
