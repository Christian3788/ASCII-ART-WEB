package tests

import (
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
}

func TestGenerateUnknownStyle(t *testing.T) {
    _, err := asciiart.Generate("Hello", "unknown")
    if err == nil {
        t.Fatal("expected error for unknown style")
    }
}
