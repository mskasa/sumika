package task

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Item struct {
	LineIndex int
	Text      string
	Done      bool
}

// ResolvePath resolves tasks_file to an absolute path.
// Supports absolute paths, ~/... paths, and paths relative to the project root.
func ResolvePath(projectPath, tasksFile string) string {
	if filepath.IsAbs(tasksFile) {
		return tasksFile
	}
	if strings.HasPrefix(tasksFile, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return tasksFile
		}
		return filepath.Join(home, tasksFile[2:])
	}
	return filepath.Join(projectPath, tasksFile)
}

// Load reads a Markdown checkbox file and returns all task items.
func Load(path string) ([]Item, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open tasks file: %w", err)
	}
	defer f.Close()

	var items []Item
	scanner := bufio.NewScanner(f)
	lineIdx := 0
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "- [ ] "):
			items = append(items, Item{LineIndex: lineIdx, Text: line[6:], Done: false})
		case strings.HasPrefix(line, "- [x] "), strings.HasPrefix(line, "- [X] "):
			items = append(items, Item{LineIndex: lineIdx, Text: line[6:], Done: true})
		}
		lineIdx++
	}
	return items, scanner.Err()
}

// Check marks the task at lineIndex as done by rewriting the file.
func Check(path string, lineIndex int) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read tasks file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	if lineIndex < 0 || lineIndex >= len(lines) {
		return fmt.Errorf("line index out of range: %d", lineIndex)
	}
	if !strings.HasPrefix(lines[lineIndex], "- [ ] ") {
		return fmt.Errorf("line %d is not an unchecked task", lineIndex)
	}

	lines[lineIndex] = "- [x] " + lines[lineIndex][6:]
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644)
}
