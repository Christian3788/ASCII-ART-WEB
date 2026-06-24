package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var tmpl *template.Template

func main() {
	// Endpoints mapped to handlers
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiHandler)

	log.Println("Server starting on http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// GET / - Renders the main page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "400 Bad Request: Method must be GET", http.StatusBadRequest)
		return
	}

	var err error
	tmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "404 Not Found: Template file missing", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

// POST /ascii-art - Processes form and displays ASCII art
func asciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "400 Bad Request: Method must be POST", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "400 Bad Request: Failed to parse form", http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")

	if text == "" || (banner != "standard" && banner != "shadow" && banner != "thinkertoy") {
		http.Error(w, "400 Bad Request: Missing text or invalid banner choice", http.StatusBadRequest)
		return
	}

	filePath := fmt.Sprintf("banners/%s.txt", banner)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "404 Not Found: Banner file missing", http.StatusNotFound)
		return
	}

	result, err := generateAscii(text, banner)
	if err != nil {
		http.Error(w, "500 Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "404 Not Found: Template missing during rendering", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, result)
}

// Core Engine: Calculates and generates the ASCII lines cleanly
func generateAscii(text, banner string) (string, error) {
	text = strings.ReplaceAll(text, "\r\n", "\n")

	filePath := fmt.Sprintf("banners/%s.txt", banner)
	content, err := ioutil.ReadFile(filePath)
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
