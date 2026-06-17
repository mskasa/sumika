package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version  int       `yaml:"version"`
	Settings Settings  `yaml:"settings"`
	Projects []Project `yaml:"projects"`
}

type Settings struct {
	Port   int    `yaml:"port"`
	Editor string `yaml:"editor"`
	AITool string `yaml:"ai_tool"`
}

type Project struct {
	Name        string   `yaml:"name"`
	Path        string   `yaml:"path"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
	TasksFile   string   `yaml:"tasks_file,omitempty"`
	Launch      Launch   `yaml:"launch"`
	Links       []Link   `yaml:"links"`
}

type Launch struct {
	Editor   bool     `yaml:"editor"`
	AI       bool     `yaml:"ai"`
	Commands []string `yaml:"commands"`
	Ports    []int    `yaml:"ports,omitempty"`
}

type Link struct {
	Label string `yaml:"label"`
	URL   string `yaml:"url"`
}

func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	return filepath.Join(home, ".config", "sumika", "config.yaml"), nil
}

func Load() (*Config, error) {
	path, err := DefaultPath()
	if err != nil {
		return nil, err
	}
	return LoadFrom(path)
}

func LoadFrom(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if cfg.Settings.Port == 0 {
		cfg.Settings.Port = 8964
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	path, err := DefaultPath()
	if err != nil {
		return err
	}
	return c.SaveTo(path)
}

func (c *Config) SaveTo(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}
