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
	// 1. Strict 404 Check for incorrect URL paths
	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	// 2. Strict Check: Ensure the method is strictly GET
	if r.Method != http.MethodGet {
		http.Error(w, "400 Bad Request: Method must be GET", http.StatusBadRequest)
		return
	}

	// 3. Dynamic Template Verification (Returns a 404 if index.html is missing)
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
	// Ensure method is strictly POST
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

	// Input Validation
	if text == "" || (banner != "standard" && banner != "shadow" && banner != "thinkertoy") {
		http.Error(w, "400 Bad Request: Missing text or invalid banner choice", http.StatusBadRequest)
		return
	}

	// 4. Verify banner existence on disk BEFORE loading (Returns a 404 if file is deleted)
	filePath := fmt.Sprintf("banners/%s.txt", banner)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "404 Not Found: Banner file missing", http.StatusNotFound)
		return
	}

	// Generate ASCII Art using your processing logic
	result, err := generateAscii(text, banner)
	if err != nil {
		http.Error(w, "500 Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Re-parse template safely for execution layout injection
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
	// 1. Sanitize web carriage returns from user input (\r\n -> \n)
	text = strings.ReplaceAll(text, "\r\n", "\n")

	// 2. Read the appropriate banner file
	filePath := fmt.Sprintf("banners/%s.txt", banner)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("could not read banner file for %s", banner)
	}

	// Standardize file line breaks to prevent indexing drifts across different OS formats
	fileStr := strings.ReplaceAll(string(content), "\r\n", "\n")
	lines := strings.Split(fileStr, "\n")

	// Safety check to ensure file loaded correctly and isn't empty
	if len(lines) < 2 {
		return "", fmt.Errorf("banner file %s appears to be empty or misformatted", banner)
	}

	// 3. Process the input chunks split by newline
	inputLines := strings.Split(text, "\n")
	var output strings.Builder

	for _, line := range inputLines {
		if line == "" {
			output.WriteString("\n")
			continue
		}

		// Each character block is 8 lines of art + 1 blank line separating characters = 9 lines total
		const blockHeight = 9

		// Build the 8 vertical slices for the row of characters simultaneously
		for i := 1; i <= 8; i++ {
			for _, runeVal := range line {
				if runeVal < 32 || runeVal > 126 {
					continue // Filter non-printable/unsupported ASCII characters
				}

				// Calculate exact mathematical starting index for standard template banners
				// Space (32) starts at line index 1 of the lines array
				startingLine := (int(runeVal)-32)*blockHeight + i

				// Bounds check to ensure we don't go past the file
				if startingLine >= len(lines) {
					continue
				}

				// Append the line from the banner file
				output.WriteString(lines[startingLine])
				output.WriteString(" ")
			}
			output.WriteString("\n")
		}
		output.WriteString("\n")
	}

	return output.String(), nil
}