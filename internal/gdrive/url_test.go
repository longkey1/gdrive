package gdrive

import "testing"

func TestParseFileInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "bare ID",
			input: "1mnXvCwniZFrx4W-divBlbO5rVMQ5ubii",
			want:  "1mnXvCwniZFrx4W-divBlbO5rVMQ5ubii",
		},
		{
			name:  "bare ID with surrounding whitespace",
			input: "  abc_123-XYZ \n",
			want:  "abc_123-XYZ",
		},
		{
			name:  "folder URL",
			input: "https://drive.google.com/drive/folders/1mnXvCwniZFrx4W-divBlbO5rVMQ5ubii",
			want:  "1mnXvCwniZFrx4W-divBlbO5rVMQ5ubii",
		},
		{
			name:  "folder URL with account index",
			input: "https://drive.google.com/drive/u/0/folders/1mnXvCwniZFrx4W-divBlbO5rVMQ5ubii",
			want:  "1mnXvCwniZFrx4W-divBlbO5rVMQ5ubii",
		},
		{
			name:  "folder URL with query and fragment",
			input: "https://drive.google.com/drive/u/2/folders/1mnXvCwniZFrx4W-divBlbO5rVMQ5ubii?usp=sharing#frag",
			want:  "1mnXvCwniZFrx4W-divBlbO5rVMQ5ubii",
		},
		{
			name:  "file URL",
			input: "https://drive.google.com/file/d/1AbCdEfGhIjKlMnOpQrStUvWxYz/view?usp=drive_link",
			want:  "1AbCdEfGhIjKlMnOpQrStUvWxYz",
		},
		{
			name:  "file URL with account index",
			input: "https://drive.google.com/file/u/1/d/1AbCdEfGhIjKlMnOpQrStUvWxYz/view",
			want:  "1AbCdEfGhIjKlMnOpQrStUvWxYz",
		},
		{
			name:  "docs URL",
			input: "https://docs.google.com/document/d/1DocId_-x/edit?tab=t.0",
			want:  "1DocId_-x",
		},
		{
			name:  "slides URL",
			input: "https://docs.google.com/presentation/d/1SlideId/edit#slide=id.p",
			want:  "1SlideId",
		},
		{
			name:  "sheets URL",
			input: "https://docs.google.com/spreadsheets/d/1SheetId/edit#gid=0",
			want:  "1SheetId",
		},
		{
			name:  "sheets URL with account index",
			input: "https://docs.google.com/spreadsheets/u/0/d/1SheetId/edit",
			want:  "1SheetId",
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: true,
		},
		{
			name:    "unrelated URL",
			input:   "https://example.com/drive/folders/abc",
			wantErr: false, // host is not checked; path layout is what matters
			want:    "abc",
		},
		{
			name:    "drive URL without ID segment",
			input:   "https://drive.google.com/drive/my-drive",
			wantErr: true,
		},
		{
			name:    "drive root URL",
			input:   "https://drive.google.com/",
			wantErr: true,
		},
		{
			name:    "malformed URL",
			input:   "https://drive.google.com/drive/folders/%zz",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseFileInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseFileInput(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseFileInput(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
