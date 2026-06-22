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
	// Dynamically handle templates during requests to safely return proper HTTP status codes
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

	// 2. Strict 400/405 Check: Ensure the method is strictly GET
	if r.Method != http.MethodGet {
		http.Error(w, "400 Bad Request: Method must be GET", http.StatusBadRequest)
		return
	}

	// 3. Dynamic Template Verification (Fixes the missing template 404 requirement)
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

	// Check if banner file exists on disk BEFORE generating (Enforces 404 requirement for missing banners)
	filePath := fmt.Sprintf("banners/%s.txt", banner)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "404 Not Found: Banner file missing", http.StatusNotFound)
		return
	}

	// Generate ASCII Art using processing logic
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

// Core Engine stays exactly the same as your logic...
func generateAscii(text, banner string) (string, error) {
    // Keep your exact processing logic here
    return "", nil 
}
