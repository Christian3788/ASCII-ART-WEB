package asciiart

import (
	"errors"
	"fmt"
	"strings"

	"github.com/christianotieno/ascii-art-web/ascii-art/artstyles"
)

var ErrBannerNotFound = errors.New("banner not found")

func Render(text, banner string) (string, error) {
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
	switch strings.ToLower(name) {
	case "standard":
		return artstyles.Standard(), nil
	case "shadow":
		return artstyles.Shadow(), nil
	case "thinkertoy":
		return artstyles.Thinkertoy(), nil
	default:
		return nil, ErrBannerNotFound
	}
}
