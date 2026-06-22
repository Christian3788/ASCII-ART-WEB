package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

const templateDir = "templates"

var errBannerNotFound = errors.New("banner not found")
var supportedBanners = []string{"standard", "shadow", "thinkertoy"}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/ascii-art", handleAsciiArt)

	fmt.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

type pageData struct {
	Text    string
	Banner  string
	Result  string
	Banners []string
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

	result, err := renderAsciiArt(text, banner)
	if err != nil {
		if errors.Is(err, errBannerNotFound) {
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
	t, err := template.ParseFiles(fmt.Sprintf("%s/%s", templateDir, name))
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

func renderAsciiArt(text, banner string) (string, error) {
	font, err := loadFont(banner)
	if err != nil {
		return "", err
	}

	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	var out strings.Builder

	for li, line := range lines {
		runes := []rune(strings.ToUpper(line))
		for row := 0; row < 6; row++ {
			for ci, ch := range runes {
				if ch == ' ' {
					out.WriteString("     ")
				} else {
					charPattern, ok := font[ch]
					if !ok {
						return "", fmt.Errorf("unsupported character %q", ch)
					}
					out.WriteString(charPattern[row])
				}
				if ci < len(runes)-1 {
					out.WriteByte(' ')
				}
			}
			out.WriteByte('\n')
		}
		if li < len(lines)-1 {
			out.WriteByte('\n')
		}
	}

	return out.String(), nil
}

func loadFont(name string) (map[rune][6]string, error) {
	if name == "standard" {
		return baseFont(), nil
	}
	if name == "shadow" {
		return shadowFont(baseFont()), nil
	}
	if name == "thinkertoy" {
		return thinkertoyFont(baseFont()), nil
	}
	return nil, errBannerNotFound
}

func shadowFont(base map[rune][6]string) map[rune][6]string {
	shadow := make(map[rune][6]string, len(base))
	for ch, pattern := range base {
		var charLines [6]string
		for i, line := range pattern {
			shadowLine := strings.Map(func(r rune) rune {
				if r == ' ' {
					return ' '
				}
				return '.'
			}, line)
			charLines[i] = fmt.Sprintf("%s  %s", line, shadowLine)
		}
		shadow[ch] = charLines
	}
	return shadow
}

func thinkertoyFont(base map[rune][6]string) map[rune][6]string {
	thinker := make(map[rune][6]string, len(base))
	for ch, pattern := range base {
		var charLines [6]string
		for i, line := range pattern {
			charLines[i] = strings.Map(func(r rune) rune {
				if r == ' ' {
					return ' '
				}
				return '@'
			}, line)
		}
		thinker[ch] = charLines
	}
	return thinker
}

func baseFont() map[rune][6]string {
	return map[rune][6]string{
		'A': {"  #  ", " # # ", "#####", "#   #", "#   #", "#   #"},
		'B': {"#### ", "#   #", "#### ", "#   #", "#   #", "#### "},
		'C': {" ####", "#    ", "#    ", "#    ", "#    ", " ####"},
		'D': {"#### ", "#   #", "#   #", "#   #", "#   #", "#### "},
		'E': {"#####", "#    ", "###  ", "#    ", "#    ", "#####"},
		'F': {"#####", "#    ", "###  ", "#    ", "#    ", "#    "},
		'G': {" ####", "#    ", "#  ##", "#   #", "#   #", " ####"},
		'H': {"#   #", "#   #", "#####", "#   #", "#   #", "#   #"},
		'I': {"#####", "  #  ", "  #  ", "  #  ", "  #  ", "#####"},
		'J': {"#####", "   # ", "   # ", "   # ", "#  # ", " ##  "},
		'K': {"#   #", "#  # ", "###  ", "#  # ", "#   #", "#   #"},
		'L': {"#    ", "#    ", "#    ", "#    ", "#    ", "#####"},
		'M': {"#   #", "## ##", "# # #", "#   #", "#   #", "#   #"},
		'N': {"#   #", "##  #", "# # #", "#  ##", "#   #", "#   #"},
		'O': {" ### ", "#   #", "#   #", "#   #", "#   #", " ### "},
		'P': {"#### ", "#   #", "#### ", "#    ", "#    ", "#    "},
		'Q': {" ### ", "#   #", "#   #", "#   #", "#  ##", " ####"},
		'R': {"#### ", "#   #", "#### ", "#  # ", "#   #", "#   #"},
		'S': {" ####", "#    ", " ### ", "    #", "    #", "#### "},
		'T': {"#####", "  #  ", "  #  ", "  #  ", "  #  ", "  #  "},
		'U': {"#   #", "#   #", "#   #", "#   #", "#   #", " ### "},
		'V': {"#   #", "#   #", "#   #", "#   #", " # # ", "  #  "},
		'W': {"#   #", "#   #", "#   #", "# # #", "## ##", "#   #"},
		'X': {"#   #", " # # ", "  #  ", "  #  ", " # # ", "#   #"},
		'Y': {"#   #", " # # ", "  #  ", "  #  ", "  #  ", "  #  "},
		'Z': {"#####", "   # ", "  #  ", " #   ", "#    ", "#####"},
		'0': {" ### ", "#   #", "#  ##", "# # #", "##  #", " ### "},
		'1': {"  #  ", " ##  ", "  #  ", "  #  ", "  #  ", " ### "},
		'2': {" ### ", "#   #", "   # ", "  #  ", " #   ", "#####"},
		'3': {" ### ", "#   #", "   # ", "   # ", "#   #", " ### "},
		'4': {"#   #", "#   #", "#####", "    #", "    #", "    #"},
		'5': {"#####", "#    ", "#### ", "    #", "    #", "#### "},
		'6': {" ### ", "#    ", "#### ", "#   #", "#   #", " ### "},
		'7': {"#####", "    #", "   # ", "  #  ", " #   ", "#    "},
		'8': {" ### ", "#   #", " ### ", "#   #", "#   #", " ### "},
		'9': {" ### ", "#   #", "#   #", " ####", "    #", " ### "},
	}
}
