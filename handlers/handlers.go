package handlers

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	asciiart "github.com/christianotieno/ascii-art-web/ascii-art"
)

const templateDir = "templates"

var (
	supportedBanners = []string{"standard", "shadow", "thinkertoy"}
)

type pageData struct {
	Text    string
	Banner  string
	Result  string
	Banners []string
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/ascii-art", handleAsciiArt)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := renderTemplate(w, "index.html", pageData{Banners: supportedBanners}); err != nil {
		handleTemplateError(w, err)
	}
}

func handleAsciiArt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	text := strings.TrimSpace(r.FormValue("text"))
	banner := strings.ToLower(strings.TrimSpace(r.FormValue("banner")))

	if text == "" || banner == "" {
		http.Error(w, "Bad Request: text and banner are required", http.StatusBadRequest)
		return
	}

	result, err := asciiart.Render(text, banner)
	if err != nil {
		if errors.Is(err, asciiart.ErrBannerNotFound) {
			http.Error(w, "Banner Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := pageData{
		Text:    text,
		Banner:  banner,
		Result:  result,
		Banners: supportedBanners,
	}

	if err := renderTemplate(w, "ascii-art.html", data); err != nil {
		handleTemplateError(w, err)
	}
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	t, err := template.ParseFiles(filepath.Join(templateDir, name))
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}

func handleTemplateError(w http.ResponseWriter, err error) {
	if os.IsNotExist(err) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
