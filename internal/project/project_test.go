package project_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mskasa/sumika/internal/config"
	"github.com/mskasa/sumika/internal/project"
)

func newCfg(projects ...config.Project) *config.Config {
	return &config.Config{
		Version:  1,
		Projects: projects,
	}
}

func TestAdd(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name        string
		path        string
		projectName string
		description string
		wantName    string
		wantErr     bool
	}{
		{
			name:     "default name from dir",
			path:     dir,
			wantName: filepath.Base(dir),
		},
		{
			name:        "custom name",
			path:        dir,
			projectName: "my-project",
			wantName:    "my-project",
		},
		{
			name:     "non-existent path resolves anyway",
			path:     filepath.Join(dir, "sub"),
			wantName: "sub",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.MkdirAll(tt.path, 0o755); err != nil {
				t.Fatal(err)
			}
			cfg := newCfg()
			err := project.Add(cfg, tt.path, tt.projectName, tt.description)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(cfg.Projects) != 1 {
				t.Fatalf("expected 1 project, got %d", len(cfg.Projects))
			}
			if cfg.Projects[0].Name != tt.wantName {
				t.Errorf("Name: got %q, want %q", cfg.Projects[0].Name, tt.wantName)
			}
		})
	}
}

func TestAdd_DuplicateName(t *testing.T) {
	dir := t.TempDir()
	cfg := newCfg(config.Project{Name: "foo", Path: dir})
	if err := project.Add(cfg, dir, "foo", ""); err == nil {
		t.Error("expected error for duplicate name, got nil")
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name     string
		existing []config.Project
		remove   string
		wantLen  int
		wantErr  bool
	}{
		{
			name:     "remove existing",
			existing: []config.Project{{Name: "foo"}, {Name: "bar"}},
			remove:   "foo",
			wantLen:  1,
		},
		{
			name:     "remove not found",
			existing: []config.Project{{Name: "foo"}},
			remove:   "baz",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newCfg(tt.existing...)
			err := project.Remove(cfg, tt.remove)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(cfg.Projects) != tt.wantLen {
				t.Errorf("Projects len: got %d, want %d", len(cfg.Projects), tt.wantLen)
			}
		})
	}
}

func TestFind(t *testing.T) {
	cfg := newCfg(
		config.Project{Name: "alpha", Path: "/alpha"},
		config.Project{Name: "beta", Path: "/beta"},
	)

	p, err := project.Find(cfg, "alpha")
	if err != nil {
		t.Fatalf("Find: %v", err)
	}
	if p.Path != "/alpha" {
		t.Errorf("Path: got %q, want %q", p.Path, "/alpha")
	}

	_, err = project.Find(cfg, "gamma")
	if err == nil {
		t.Error("expected error for missing project, got nil")
	}
}
