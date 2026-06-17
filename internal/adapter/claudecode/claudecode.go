package claudecode

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mskasa/sumika/internal/adapter"
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

// projectDirName converts a project path to the directory name used by Claude Code.
// Both '/' and '.' are replaced with '-'.
// e.g. /Users/foo.bar/baz → -Users-foo-bar-baz
func projectDirName(projectPath string) string {
	r := strings.ReplaceAll(projectPath, "/", "-")
	return strings.ReplaceAll(r, ".", "-")
}

type jsonlEntry struct {
	Message struct {
		Role    string `json:"role"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"message"`
	Timestamp string `json:"timestamp"`
}

func (a *Adapter) GetSessionSummary(projectPath string) (*adapter.SessionSummary, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home dir: %w", err)
	}

	sessionDir := filepath.Join(home, ".claude", "projects", projectDirName(projectPath))
	entries, err := os.ReadDir(sessionDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read session dir: %w", err)
	}

	// JSONLファイルを新しい順に並べる
	type sessionFile struct {
		path    string
		modTime time.Time
	}
	var files []sessionFile
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		files = append(files, sessionFile{
			path:    filepath.Join(sessionDir, e.Name()),
			modTime: info.ModTime(),
		})
	}
	if len(files) == 0 {
		return nil, nil
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].modTime.After(files[j].modTime)
	})

	// アシスタントメッセージが見つかるまで新しい順に遡る
	var summary string
	var lastActive time.Time
	for _, f := range files {
		s, t, err := parseLatestAssistantText(f.path)
		if err != nil {
			continue
		}
		if s != "" {
			summary = s
			lastActive = t
			break
		}
	}
	if summary == "" {
		return nil, nil
	}

	return &adapter.SessionSummary{
		ProjectName: filepath.Base(projectPath),
		LastActive:  lastActive,
		Summary:     truncate(summary, 200),
	}, nil
}

// parseLatestAssistantText reads the JSONL file and returns the last assistant text message.
func parseLatestAssistantText(path string) (string, time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("open jsonl: %w", err)
	}
	defer f.Close()

	var lastText string
	var lastTime time.Time

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB per line
	for scanner.Scan() {
		var e jsonlEntry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			continue
		}
		if e.Message.Role != "assistant" {
			continue
		}
		for _, c := range e.Message.Content {
			if c.Type == "text" && strings.TrimSpace(c.Text) != "" {
				t, _ := time.Parse(time.RFC3339, e.Timestamp)
				lastText = strings.TrimSpace(c.Text)
				lastTime = t
				break
			}
		}
	}
	return lastText, lastTime, scanner.Err()
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

func truncate(s string, maxRunes int) string {
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxRunes]) + "…"
}
