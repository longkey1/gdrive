package gdrive

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
)

// OutputFormat represents the output format type
type OutputFormat string

const (
	OutputFormatText OutputFormat = "text"
	OutputFormatJSON OutputFormat = "json"
)

// ParseOutputFormat validates a user-supplied format string
func ParseOutputFormat(s string) (OutputFormat, error) {
	switch OutputFormat(s) {
	case OutputFormatText, OutputFormatJSON:
		return OutputFormat(s), nil
	default:
		return "", fmt.Errorf("unsupported format %q (must be text or json)", s)
	}
}

// FormatFileList outputs a file list in the specified format. JSON output is
// always an array, "[]" when empty.
func FormatFileList(w io.Writer, files []FileInfo, format OutputFormat) error {
	if format == OutputFormatJSON {
		if files == nil {
			files = []FileInfo{}
		}
		return writeJSON(w, files)
	}

	table := tablewriter.NewWriter(w)
	table.Header("ID", "NAME", "TYPE", "MODIFIED")

	for _, f := range files {
		_ = table.Append(f.ID, f.Name, f.MimeType, f.ModifiedTime)
	}

	_ = table.Render()
	return nil
}

// FormatFile outputs a single file's metadata in the specified format
func FormatFile(w io.Writer, file *FileInfo, format OutputFormat) error {
	if format == OutputFormatJSON {
		return writeJSON(w, file)
	}

	table := tablewriter.NewWriter(w)
	table.Header("FIELD", "VALUE")
	_ = table.Append("ID", file.ID)
	_ = table.Append("NAME", file.Name)
	_ = table.Append("TYPE", file.MimeType)
	_ = table.Append("SIZE", fmt.Sprintf("%d", file.Size))
	_ = table.Append("MODIFIED", file.ModifiedTime)
	_ = table.Append("LINK", file.WebViewLink)
	for _, p := range file.Parents {
		_ = table.Append("PARENT", p)
	}
	for _, o := range file.Owners {
		_ = table.Append("OWNER", fmt.Sprintf("%s <%s>", o.DisplayName, o.EmailAddress))
	}
	_ = table.Render()
	return nil
}

func writeJSON(w io.Writer, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal JSON: %w", err)
	}
	_, _ = fmt.Fprintln(w, string(data))
	return nil
}
