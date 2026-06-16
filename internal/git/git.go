package git

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
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
	LastCommitTime   time.Time
	UncommittedCount int
	IsRepo           bool
}

func GetStatus(dir string, runner Runner) (*Status, error) {
	if runner == nil {
		runner = &CLIRunner{}
	}

	raw, err := runner.Run(dir, "log", "-1", "--format=%cI")
	if err != nil {
		return &Status{IsRepo: false}, nil
	}

	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, fmt.Errorf("parse commit time %q: %w", raw, err)
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
		LastCommitTime:   t,
		UncommittedCount: count,
		IsRepo:           true,
	}, nil
}
