package claudecode

import "testing"

func TestProjectDirName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/Users/foo/bar", "-Users-foo-bar"},
		{"/Users/foo.bar/baz", "-Users-foo-bar-baz"},
		{"/Users/masahiro.kasatani/Documents/github/mskasa/sumika",
			"-Users-masahiro-kasatani-Documents-github-mskasa-sumika"},
	}
	for _, tt := range tests {
		got := projectDirName(tt.path)
		if got != tt.want {
			t.Errorf("projectDirName(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		max      int
		wantTail string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "…"},
		{"あいうえお", 3, "…"},
	}
	for _, tt := range tests {
		got := truncate(tt.input, tt.max)
		if len([]rune(tt.input)) <= tt.max {
			if got != tt.input {
				t.Errorf("truncate(%q, %d) = %q, want unchanged", tt.input, tt.max, got)
			}
		} else {
			if got[len(got)-len("…"):] != "…" {
				t.Errorf("truncate(%q, %d) = %q, want trailing ellipsis", tt.input, tt.max, got)
			}
		}
	}
}
