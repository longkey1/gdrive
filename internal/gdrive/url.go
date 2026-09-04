package gdrive

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// filePathPatterns match the path segment that carries the file or folder ID
// in the Google Drive / Docs / Sheets / Slides URL layouts.
var filePathPatterns = []*regexp.Regexp{
	// https://drive.google.com/drive/folders/<id>
	// https://drive.google.com/drive/u/0/folders/<id>
	regexp.MustCompile(`^/drive/(?:u/\d+/)?folders/([a-zA-Z0-9_-]+)`),
	// https://drive.google.com/file/d/<id>/view
	regexp.MustCompile(`^/file/(?:u/\d+/)?d/([a-zA-Z0-9_-]+)`),
	// https://docs.google.com/document/d/<id>/edit
	// https://docs.google.com/presentation/d/<id>/edit
	// https://docs.google.com/spreadsheets/d/<id>/edit
	regexp.MustCompile(`^/(?:document|presentation|spreadsheets)/(?:u/\d+/)?d/([a-zA-Z0-9_-]+)`),
}

// ParseFileInput accepts either a bare Drive file/folder ID or a Google
// Drive, Docs, Sheets, or Slides URL and returns the ID. Bare IDs are
// returned as-is; query strings and fragments on URLs are ignored.
func ParseFileInput(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("file ID or URL is required")
	}

	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		return input, nil
	}

	u, err := url.Parse(input)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	for _, p := range filePathPatterns {
		if m := p.FindStringSubmatch(u.Path); len(m) == 2 {
			return m[1], nil
		}
	}

	return "", fmt.Errorf("could not find file ID in URL: %s", input)
}
