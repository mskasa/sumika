package launcher

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

func Open(dir, editor, aiTool string, commands []string) error {
	if editor != "" {
		if err := launch(dir, editor); err != nil {
			slog.Warn("failed to launch editor", "editor", editor, "err", err)
		}
	}
	if aiTool != "" {
		if err := launch(dir, aiTool); err != nil {
			slog.Warn("failed to launch ai tool", "ai_tool", aiTool, "err", err)
		}
	}
	for _, c := range commands {
		if err := runCommand(dir, c); err != nil {
			slog.Warn("failed to run command", "cmd", c, "err", err)
		}
	}
	return nil
}

func launch(dir, command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = dir
	return cmd.Start()
}

func runCommand(dir, command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = dir
	return cmd.Start()
}
