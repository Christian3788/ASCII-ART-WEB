package asciiart_test

import (
	"strings"
	"testing"

	"github.com/christianotieno/ascii-art-web/ascii-art"
)

func TestRenderStandard(t *testing.T) {
	out, err := asciiart.Render("A", "standard")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "#####") {
		t.Fatalf("expected standard banner output, got %q", out)
	}
}

func TestRenderShadow(t *testing.T) {
	out, err := asciiart.Render("A", "shadow")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "..") {
		t.Fatalf("expected shadow banner output, got %q", out)
	}
}

func TestUnknownBanner(t *testing.T) {
	_, err := asciiart.Render("A", "unknown")
	if err == nil {
		t.Fatal("expected error for unknown banner")
	}
	if !strings.Contains(err.Error(), "banner not found") {
		t.Fatalf("expected banner not found, got %v", err)
	}
}
