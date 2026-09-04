package gdrive

import "testing"

func TestBuildListQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts ListOptions
		want string
	}{
		{
			name: "parent only",
			opts: ListOptions{Parent: "folder123"},
			want: "'folder123' in parents and trashed = false",
		},
		{
			name: "parent and query",
			opts: ListOptions{Parent: "folder123", Query: "report"},
			want: "'folder123' in parents and trashed = false and name contains 'report'",
		},
		{
			name: "single quote in query is escaped",
			opts: ListOptions{Parent: "folder123", Query: "it's"},
			want: `'folder123' in parents and trashed = false and name contains 'it\'s'`,
		},
		{
			name: "backslash in query is escaped",
			opts: ListOptions{Parent: "folder123", Query: `a\b`},
			want: `'folder123' in parents and trashed = false and name contains 'a\\b'`,
		},
		{
			name: "single quote in parent is escaped",
			opts: ListOptions{Parent: "x'y"},
			want: `'x\'y' in parents and trashed = false`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := BuildListQuery(tt.opts); got != tt.want {
				t.Errorf("BuildListQuery() = %q, want %q", got, tt.want)
			}
		})
	}
}
