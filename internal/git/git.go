package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type Runner interface {
	Run(dir string, args ...string) (string, error)
}

type CLIRunner struct{}

func (r *CLIRunner) Run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out)), nil
}

type Status struct {
	LastCommit     string
	UncommittedCount int
	IsRepo         bool
}

func GetStatus(dir string, runner Runner) (*Status, error) {
	if runner == nil {
		runner = &CLIRunner{}
	}

	lastCommit, err := runner.Run(dir, "log", "-1", "--format=%ci")
	if err != nil {
		return &Status{IsRepo: false}, nil
	}

	porcelain, err := runner.Run(dir, "status", "--porcelain")
	if err != nil {
		return nil, fmt.Errorf("git status: %w", err)
	}

	count := 0
	if porcelain != "" {
		count = len(strings.Split(porcelain, "\n"))
	}

	return &Status{
		LastCommit:     lastCommit,
		UncommittedCount: count,
		IsRepo:         true,
	}, nil
}
