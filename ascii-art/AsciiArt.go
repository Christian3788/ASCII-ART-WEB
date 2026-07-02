package asciiart

import (
    "embed"
    "fmt"
    "strings"
    "sync"
)

//go:embed artstyles/*.txt
var artStyleFS embed.FS

var (
    AvailableStyles = []string{"standard", "shadow", "thinkertoy"}
    fontRegistry   = map[string]font{}
    loadOnce       sync.Once
    loadErr        error
)

type font struct {
    Height int
    Chars  map[rune][]string
}

func Generate(text, style string) (string, error) {
    loadOnce.Do(initFonts)
    if loadErr != nil {
        return "", loadErr
    }

    fontData, ok := fontRegistry[style]
    if !ok {
        return "", fmt.Errorf("unknown style %q", style)
    }

    normalized := strings.ReplaceAll(text, "\\n", "\n")
    if strings.TrimSpace(normalized) == "" {
        return "", fmt.Errorf("text must not be empty")
    }

    var lines []string
    for _, inputLine := range strings.Split(normalized, "\n") {
        lines = append(lines, renderLine(inputLine, fontData))
    }

    return strings.Join(lines, "\n"), nil
}

func initFonts() {
    for _, style := range AvailableStyles {
        raw, err := artStyleFS.ReadFile("artstyles/" + style + ".txt")
        if err != nil {
            loadErr = fmt.Errorf("failed to load style %q: %w", style, err)
            return
        }

        parsed, err := parseFont(string(raw))
        if err != nil {
            loadErr = fmt.Errorf("failed to parse style %q: %w", style, err)
            return
        }

        fontRegistry[style] = parsed
    }
}

func parseFont(raw string) (font, error) {
    normalized := strings.ReplaceAll(raw, "\r\n", "\n")
    lines := strings.Split(normalized, "\n")

    for len(lines)%9 != 0 {
        lines = append(lines, "")
    }

    const startRune = 32
    charCount := len(lines) / 9
    if charCount < 1 {
        return font{}, fmt.Errorf("font file contains no glyphs")
    }

    chars := make(map[rune][]string, charCount)
    for i := 0; i < charCount; i++ {
        block := lines[i*9 : i*9+8]
        if len(block) < 8 {
            padded := make([]string, 8)
            copy(padded, block)
            block = padded
        }
        chars[rune(startRune+i)] = block
    }

    return font{Height: 8, Chars: chars}, nil
}

func renderLine(line string, f font) string {
    if line == "" {
        return ""
    }

    runes := []rune(line)
    rows := make([]strings.Builder, f.Height)

    for idx, r := range runes {
        if r < 32 || r > 126 {
            r = '?'
        }

        glyph, ok := f.Chars[r]
        if !ok {
            glyph, ok = f.Chars['?']
            if !ok {
                glyph = make([]string, f.Height)
            }
        }

        for i := 0; i < f.Height; i++ {
            rows[i].WriteString(glyph[i])
            if idx < len(runes)-1 {
                rows[i].WriteByte(' ')
            }
        }
    }

    out := make([]string, f.Height)
    for i := 0; i < f.Height; i++ {
        out[i] = strings.TrimRight(rows[i].String(), " ")
    }

    return strings.Join(out, "\n")
}
