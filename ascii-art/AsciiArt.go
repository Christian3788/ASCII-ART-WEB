package asciiart

import "fmt"

var AvailableStyles = []string{"standard", "shadow", "thinkertoy"}

func Generate(text, style string) (string, error) {
    if text == "" {
        return "", fmt.Errorf("text must not be empty")
    }

    switch style {
    case "standard", "shadow", "thinkertoy":
        return fmt.Sprintf("[%s] %s", style, text), nil
    default:
        return "", fmt.Errorf("unknown style %q", style)
    }
}
