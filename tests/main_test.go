package tests

import (
    "strings"
    "testing"

    "ascii-art-web/ascii-art"
)

func TestGenerateAsciiArt(t *testing.T) {
    output, err := asciiart.Generate("Hello", "standard")
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if output == "" {
        t.Fatal("expected generated art text")
    }
    if !strings.Contains(output, "H") && !strings.Contains(output, "|") {
        t.Fatalf("unexpected ascii art output: %q", output)
    }
}

func TestGenerateAsciiArtWithNewline(t *testing.T) {
    output, err := asciiart.Generate("Hi\\nYou", "shadow")
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if strings.Count(output, "\n") < 2 {
        t.Fatalf("expected multiline output, got %q", output)
    }
}

func TestGenerateUnknownStyle(t *testing.T) {
    _, err := asciiart.Generate("Hello", "unknown")
    if err == nil {
        t.Fatal("expected error for unknown style")
    }
}
