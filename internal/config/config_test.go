package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mskasa/sumika/internal/config"
)

func TestLoadFrom(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    *config.Config
		wantErr bool
	}{
		{
			name: "valid config",
			yaml: `version: 1
settings:
  port: 9000
  editor: code
  ai_tool: claude
projects:
  - name: my-api
    path: /projects/my-api
    description: REST API
    tags:
      - backend
`,
			want: &config.Config{
				Version: 1,
				Settings: config.Settings{
					Port:   9000,
					Editor: "code",
					AITool: "claude",
				},
				Projects: []config.Project{
					{
						Name:        "my-api",
						Path:        "/projects/my-api",
						Description: "REST API",
						Tags:        []string{"backend"},
					},
				},
			},
		},
		{
			name: "default port when zero",
			yaml: `version: 1
settings: {}
`,
			want: &config.Config{
				Version:  1,
				Settings: config.Settings{Port: 8964},
			},
		},
		{
			name:    "invalid yaml",
			yaml:    ":\tinvalid:",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "config.yaml")
			if err := os.WriteFile(path, []byte(tt.yaml), 0o644); err != nil {
				t.Fatal(err)
			}

			got, err := config.LoadFrom(path)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Version != tt.want.Version {
				t.Errorf("Version: got %d, want %d", got.Version, tt.want.Version)
			}
			if got.Settings.Port != tt.want.Settings.Port {
				t.Errorf("Settings.Port: got %d, want %d", got.Settings.Port, tt.want.Settings.Port)
			}
			if got.Settings.Editor != tt.want.Settings.Editor {
				t.Errorf("Settings.Editor: got %q, want %q", got.Settings.Editor, tt.want.Settings.Editor)
			}
			if len(got.Projects) != len(tt.want.Projects) {
				t.Errorf("Projects len: got %d, want %d", len(got.Projects), len(tt.want.Projects))
			}
		})
	}
}

func TestSaveTo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "config.yaml")

	cfg := &config.Config{
		Version: 1,
		Settings: config.Settings{
			Port:   8964,
			Editor: "nvim",
		},
		Projects: []config.Project{
			{Name: "foo", Path: "/foo"},
		},
	}

	if err := cfg.SaveTo(path); err != nil {
		t.Fatalf("SaveTo: %v", err)
	}

	loaded, err := config.LoadFrom(path)
	if err != nil {
		t.Fatalf("LoadFrom after save: %v", err)
	}
	if loaded.Settings.Editor != "nvim" {
		t.Errorf("Editor: got %q, want %q", loaded.Settings.Editor, "nvim")
	}
	if len(loaded.Projects) != 1 || loaded.Projects[0].Name != "foo" {
		t.Errorf("Projects: got %+v", loaded.Projects)
	}
}

func TestLoadFrom_FileNotFound(t *testing.T) {
	_, err := config.LoadFrom("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
