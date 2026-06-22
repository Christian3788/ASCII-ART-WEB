package main

import (
    "errors"
    "fmt"
    "html/template"
    "net/http"
    "os"
    "path/filepath"
    "sort"
    "strings"
)

const (
    templateDir = "templates"
    bannerDir   = "banners"
)

var defaultBanners = []string{"standard", "shadow", "thinkertoy"}

type pageData struct {
    Text    string
    Banner  string
    Result  string
    Banners []string
    Error   string
}

func runServer() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", handleIndex)
    mux.HandleFunc("/ascii-art", handleAsciiArt)

    fmt.Println("Server running at http://localhost:8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        fmt.Fprintf(os.Stderr, "server error: %v\n", err)
        os.Exit(1)
    }
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    data := pageData{Banners: listBanners(), Banner: "standard"}
    if err := renderTemplate(w, "index.html", data); err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
    banners := listBanners()
    data := pageData{Text: text, Banner: banner, Banners: banners}

    if text == "" {
        data.Error = "Text is required"
        _ = renderTemplate(w, "index.html", data)
        return
    }
    if banner == "" {
        data.Error = "Banner selection is required"
        _ = renderTemplate(w, "index.html", data)
        return
    }

    result, err := Render(text, banner)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            data.Error = "Banner not found"
            _ = renderTemplate(w, "index.html", data)
            return
        }
        data.Error = err.Error()
        _ = renderTemplate(w, "index.html", data)
        return
    }

    data.Result = result
    if err := renderTemplate(w, "index.html", data); err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

func listBanners() []string {
    entries, err := os.ReadDir(bannerDir)
    if err != nil {
        return defaultBanners
    }

    var banners []string
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        name := strings.ToLower(entry.Name())
        if strings.HasSuffix(name, ".txt") {
            base := strings.TrimSuffix(name, filepath.Ext(name))
            banners = append(banners, base)
        }
    }

    if len(banners) == 0 {
        return defaultBanners
    }

    sort.Strings(banners)
    return banners
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) error {
    t, err := template.ParseFiles(filepath.Join(templateDir, name))
    if err != nil {
        return err
    }
    return t.Execute(w, data)
}

func Render(text, banner string) (string, error) {
    font, err := loadFont(banner)
    if err != nil {
        return "", err
    }

    text = strings.ReplaceAll(text, "\r\n", "\n")
    lines := strings.Split(text, "\n")
    var out strings.Builder

    for li, line := range lines {
        if line == "" {
            if li < len(lines)-1 {
                out.WriteString("\n")
            }
            continue
        }

        runes := []rune(strings.ToUpper(line))
        patterns := make([][]string, len(runes))
        widths := make([]int, len(runes))
        maxHeight := 0

        for i, ch := range runes {
            if ch == ' ' {
                widths[i] = 5
                continue
            }

            pattern, ok := font[ch]
            if !ok {
                return "", fmt.Errorf("unsupported character %q", ch)
            }
            patterns[i] = pattern
            for _, row := range pattern {
                if len(row) > widths[i] {
                    widths[i] = len(row)
                }
            }
            if len(pattern) > maxHeight {
                maxHeight = len(pattern)
            }
        }

        for row := 0; row < maxHeight; row++ {
            for i, ch := range runes {
                if i > 0 {
                    out.WriteByte(' ')
                }

                if ch == ' ' {
                    out.WriteString(strings.Repeat(" ", widths[i]))
                    continue
                }

                pattern := patterns[i]
                if row < len(pattern) {
                    out.WriteString(pattern[row])
                } else {
                    out.WriteString(strings.Repeat(" ", widths[i]))
                }
            }
            out.WriteString("\n")
        }

        if li < len(lines)-1 {
            out.WriteString("\n")
        }
    }

    return out.String(), nil
}

func loadFont(name string) (map[rune][]string, error) {
    path := filepath.Join(bannerDir, strings.ToLower(name)+".txt")
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    return parseBannerFile(string(data)), nil
}

func parseBannerFile(contents string) map[rune][]string {
    contents = strings.ReplaceAll(contents, "\r\n", "\n")
    contents = strings.ReplaceAll(contents, "\r", "\n")
    lines := strings.Split(contents, "\n")

    var blocks [][]string
    curr := []string{}
    for _, line := range lines {
        if strings.TrimSpace(line) == "" {
            blocks = append(blocks, curr)
            curr = []string{}
            continue
        }
        curr = append(curr, line)
    }
    blocks = append(blocks, curr)

    font := make(map[rune][]string)
    for i, block := range blocks {
        code := rune(32 + i)
        font[code] = block
    }
    return font
}
