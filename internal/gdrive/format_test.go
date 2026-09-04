package gdrive

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestParseOutputFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input   string
		want    OutputFormat
		wantErr bool
	}{
		{input: "text", want: OutputFormatText},
		{input: "json", want: OutputFormatJSON},
		{input: "yaml", wantErr: true},
		{input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			got, err := ParseOutputFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseOutputFormat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseOutputFormat(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatFileListJSON(t *testing.T) {
	t.Parallel()

	t.Run("nil slice renders empty array", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		if err := FormatFileList(&buf, nil, OutputFormatJSON); err != nil {
			t.Fatal(err)
		}
		if got := strings.TrimSpace(buf.String()); got != "[]" {
			t.Errorf("got %q, want %q", got, "[]")
		}
	})

	t.Run("empty slice renders empty array", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		if err := FormatFileList(&buf, []FileInfo{}, OutputFormatJSON); err != nil {
			t.Fatal(err)
		}
		if got := strings.TrimSpace(buf.String()); got != "[]" {
			t.Errorf("got %q, want %q", got, "[]")
		}
	})

	t.Run("required fields are always present as strings", func(t *testing.T) {
		t.Parallel()

		files := []FileInfo{{ID: "id1"}}
		var buf bytes.Buffer
		if err := FormatFileList(&buf, files, OutputFormatJSON); err != nil {
			t.Fatal(err)
		}

		var decoded []map[string]any
		if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
			t.Fatalf("output is not a JSON array: %v\n%s", err, buf.String())
		}
		if len(decoded) != 1 {
			t.Fatalf("got %d elements, want 1", len(decoded))
		}
		for _, key := range []string{"id", "name", "mimeType", "webViewLink", "modifiedTime"} {
			v, ok := decoded[0][key]
			if !ok {
				t.Errorf("key %q missing from output", key)
				continue
			}
			if _, isString := v.(string); !isString {
				t.Errorf("key %q = %v (%T), want string", key, v, v)
			}
		}
		if _, ok := decoded[0]["size"]; !ok {
			t.Errorf("key %q missing from output", "size")
		}
	})
}

func TestFormatFileListText(t *testing.T) {
	t.Parallel()

	files := []FileInfo{
		{ID: "id1", Name: "recording.mp4", MimeType: "video/mp4", ModifiedTime: "2026-09-01T00:00:00Z"},
		{ID: "id2", Name: "slides", MimeType: "application/vnd.google-apps.presentation", ModifiedTime: "2026-09-02T00:00:00Z"},
	}
	var buf bytes.Buffer
	if err := FormatFileList(&buf, files, OutputFormatText); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	for _, want := range []string{"ID", "NAME", "TYPE", "MODIFIED", "id1", "recording.mp4", "video/mp4", "id2", "slides"} {
		if !strings.Contains(out, want) {
			t.Errorf("text output missing %q:\n%s", want, out)
		}
	}
}

func TestFormatFileJSON(t *testing.T) {
	t.Parallel()

	file := &FileInfo{
		ID:       "id1",
		Name:     "doc",
		MimeType: "application/vnd.google-apps.document",
		Parents:  []string{"folder1"},
		Owners:   []OwnerInfo{{EmailAddress: "a@example.com", DisplayName: "A"}},
	}
	var buf bytes.Buffer
	if err := FormatFile(&buf, file, OutputFormatJSON); err != nil {
		t.Fatal(err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not a JSON object: %v\n%s", err, buf.String())
	}
	if decoded["id"] != "id1" {
		t.Errorf("id = %v, want id1", decoded["id"])
	}
	parents, ok := decoded["parents"].([]any)
	if !ok || len(parents) != 1 || parents[0] != "folder1" {
		t.Errorf("parents = %v, want [folder1]", decoded["parents"])
	}
	owners, ok := decoded["owners"].([]any)
	if !ok || len(owners) != 1 {
		t.Fatalf("owners = %v, want one owner", decoded["owners"])
	}
	owner, ok := owners[0].(map[string]any)
	if !ok || owner["emailAddress"] != "a@example.com" {
		t.Errorf("owners[0] = %v, want emailAddress a@example.com", owners[0])
	}
}
