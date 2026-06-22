package main

import (
	"html/template"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

// TestMain sets up the test environment before running any unit tests.
// This ensures that global variables like 'tmpl' are initialized properly.
func TestMain(m *testing.M) {
	var err error
	// Safely parse template relative to root project scope
	tmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("Test setup failure: failed to compile template file: %v", err)
	}

	// Run the tests and exit with the correct status code
	os.Exit(m.Run())
}

// Tests that GET / returns 200 OK
func TestHomeHandler_Success(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %v", rr.Code)
	}
}

// Tests that a random URL path returns 404 Not Found
func TestHomeHandler_NotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/some-random-broken-path", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404 Not Found, got %v", rr.Code)
	}
}

// Tests that POST /ascii-art works with valid data
func TestAsciiHandler_Success(t *testing.T) {
	form := url.Values{}
	form.Add("text", "hello")
	form.Add("banner", "standard")

	req, err := http.NewRequest("POST", "/ascii-art", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(asciiHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %v", rr.Code)
	}
}

// Tests that POST /ascii-art returns 400 Bad Request when form data is invalid
func TestAsciiHandler_BadRequest(t *testing.T) {
	form := url.Values{}
	form.Add("text", "") // Empty string should trigger 400 Bad Request
	form.Add("banner", "invalid-banner-name")

	req, err := http.NewRequest("POST", "/ascii-art", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(asciiHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %v", rr.Code)
	}
}