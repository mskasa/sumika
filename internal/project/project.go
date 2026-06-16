package project

import (
	"fmt"
	"path/filepath"

	"github.com/mskasa/sumika/internal/config"
)

func Add(cfg *config.Config, path, name, description string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}
	if name == "" {
		name = filepath.Base(abs)
	}
	for _, p := range cfg.Projects {
		if p.Name == name {
			return fmt.Errorf("project %q already exists", name)
		}
	}
	cfg.Projects = append(cfg.Projects, config.Project{
		Name:        name,
		Path:        abs,
		Description: description,
		Launch: config.Launch{
			Editor: true,
			AI:     true,
		},
	})
	return nil
}

func Remove(cfg *config.Config, name string) error {
	for i, p := range cfg.Projects {
		if p.Name == name {
			cfg.Projects = append(cfg.Projects[:i], cfg.Projects[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("project %q not found", name)
}

func Find(cfg *config.Config, name string) (*config.Project, error) {
	for i := range cfg.Projects {
		if cfg.Projects[i].Name == name {
			return &cfg.Projects[i], nil
		}
	}
	return nil, fmt.Errorf("project %q not found", name)
}
