package git_test

import (
	"fmt"
	"testing"

	"github.com/mskasa/sumika/internal/git"
)

type mockRunner struct {
	responses map[string]string
	errors    map[string]error
}

func (m *mockRunner) Run(dir string, args ...string) (string, error) {
	key := fmt.Sprintf("%v", args)
	if err, ok := m.errors[key]; ok {
		return "", err
	}
	if v, ok := m.responses[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("unexpected git call: %v", args)
}

func TestGetStatus(t *testing.T) {
	tests := []struct {
		name      string
		runner    *mockRunner
		wantIsRepo bool
		wantCommit string
		wantCount  int
	}{
		{
			name: "clean repo",
			runner: &mockRunner{
				responses: map[string]string{
					"[log -1 --format=%ci]":   "2026-06-16 10:00:00 +0900",
					"[status --porcelain]": "",
				},
			},
			wantIsRepo: true,
			wantCommit: "2026-06-16 10:00:00 +0900",
			wantCount:  0,
		},
		{
			name: "repo with uncommitted changes",
			runner: &mockRunner{
				responses: map[string]string{
					"[log -1 --format=%ci]":   "2026-06-16 10:00:00 +0900",
					"[status --porcelain]": " M file1.go\n?? file2.go",
				},
			},
			wantIsRepo: true,
			wantCommit: "2026-06-16 10:00:00 +0900",
			wantCount:  2,
		},
		{
			name: "not a git repo",
			runner: &mockRunner{
				errors: map[string]error{
					"[log -1 --format=%ci]": fmt.Errorf("not a git repo"),
				},
			},
			wantIsRepo: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st, err := git.GetStatus("/some/dir", tt.runner)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if st.IsRepo != tt.wantIsRepo {
				t.Errorf("IsRepo: got %v, want %v", st.IsRepo, tt.wantIsRepo)
			}
			if !tt.wantIsRepo {
				return
			}
			if st.LastCommit != tt.wantCommit {
				t.Errorf("LastCommit: got %q, want %q", st.LastCommit, tt.wantCommit)
			}
			if st.UncommittedCount != tt.wantCount {
				t.Errorf("UncommittedCount: got %d, want %d", st.UncommittedCount, tt.wantCount)
			}
		})
	}
}
