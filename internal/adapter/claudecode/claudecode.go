package claudecode

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Adapter struct{}

func New() *Adapter {
	return &Adapter{}
}

func (a *Adapter) Name() string {
	return "claude-code"
}

func (a *Adapter) IsAvailable() bool {
	_, err := exec.LookPath("claude")
	return err == nil
}

func (a *Adapter) Launch(projectPath string) error {
	cmd := exec.Command("claude")
	cmd.Dir = projectPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("launch claude: %w", err)
	}
	return nil
}

func (a *Adapter) GetContextFile(projectPath string) (string, error) {
	p := filepath.Join(projectPath, "CLAUDE.md")
	data, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("read CLAUDE.md: %w", err)
	}
	return string(data), nil
}
