package handlers

import (
    "html/template"
    "net/http"
    "path/filepath"

    "ascii-art-web/ascii-art"
)

type TemplateData struct {
    Art           string
    Error         string
    Text          string
    SelectedStyle string
    Styles        []string
}

func RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/generate", generateHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, TemplateData{Styles: asciiart.AvailableStyles})
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        renderTemplate(w, TemplateData{Error: "Unable to parse form.", Styles: asciiart.AvailableStyles})
        return
    }

    text := r.FormValue("text")
    style := r.FormValue("style")
    art, err := asciiart.Generate(text, style)
    data := TemplateData{
        Text:          text,
        SelectedStyle: style,
        Styles:        asciiart.AvailableStyles,
    }

    if err != nil {
        data.Error = err.Error()
    } else {
        data.Art = art
    }

    renderTemplate(w, data)
}

func renderTemplate(w http.ResponseWriter, data TemplateData) {
    layout := filepath.Join("templates", "layout.html")
    index := filepath.Join("templates", "index.html")
    tpl, err := template.ParseFiles(layout, index)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err := tpl.ExecuteTemplate(w, "layout", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
