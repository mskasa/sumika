package launcher

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

// Open is for CLI use: runs commands and editor in background, then launches the AI tool
// in the foreground with inherited stdio (blocking until the user exits).
func Open(dir, editor, aiTool string, commands []string) error {
	for _, c := range commands {
		if err := launchBackground(dir, c); err != nil {
			slog.Warn("failed to run command", "cmd", c, "err", err)
		}
	}
	if editor != "" {
		if err := launchEditor(dir, editor); err != nil {
			slog.Warn("failed to launch editor", "editor", editor, "err", err)
		}
	}
	if aiTool != "" {
		return launchForeground(dir, aiTool)
	}
	return nil
}

// OpenBackground is for web dashboard use: launches editor in background.
// Commands are handled separately via RunCommands (Start button).
// The AI tool is intentionally omitted — no terminal is available from a web context.
func OpenBackground(dir, editor string) error {
	if editor != "" {
		if err := launchEditor(dir, editor); err != nil {
			slog.Warn("failed to launch editor", "editor", editor, "err", err)
		}
	}
	return nil
}

// RunCommands runs each command in the background (for web dashboard Start button).
func RunCommands(dir string, commands []string) {
	for _, c := range commands {
		if err := launchBackground(dir, c); err != nil {
			slog.Warn("failed to run command", "cmd", c, "err", err)
		}
	}
}

// launchEditor starts an editor with the project directory as an argument (e.g. `code /path/to/project`).
func launchEditor(dir, editor string) error {
	parts := strings.Fields(editor)
	if len(parts) == 0 {
		return fmt.Errorf("empty editor command")
	}
	args := append(parts[1:], dir)
	cmd := exec.Command(parts[0], args...)
	cmd.Dir = dir
	return cmd.Start()
}

// launchBackground starts a process detached from the terminal (no stdio).
func launchBackground(dir, command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = dir
	return cmd.Start()
}

// launchForeground runs a process in the foreground with inherited stdin/stdout/stderr.
// It blocks until the process exits (e.g. the user quits the AI tool).
func launchForeground(dir, command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
