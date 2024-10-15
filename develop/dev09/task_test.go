package main

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "wget_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	filepath := filepath.Join(tempDir, "testfile.txt")
	err = os.WriteFile(filepath, []byte("This is a test file for wget utility.\n"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestExtractLinks(t *testing.T) {
	htmlContent := `
	<html>
		<body>
			<a href="/test">link</a>
			<img src="/image.jpg">
			<script src="/script.js"></script>
		</body>
	</html>`

	baseURL, _ := url.Parse("http://example.com")

	links, err := extractResourceLinks(baseURL, strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedLinks := []string{
		"http://example.com/test",
		"http://example.com/image.jpg",
		"http://example.com/script.js",
	}

	for i, link := range links {
		if link != expectedLinks[i] {
			t.Fatalf("expected %q, got %q", expectedLinks[i], link)
		}
	}
}

func TestDownloadPage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "wget_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	filepath := filepath.Join(tempDir, "index.html")
	err = os.WriteFile(filepath, []byte("<html><body>Test</body></html>"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
