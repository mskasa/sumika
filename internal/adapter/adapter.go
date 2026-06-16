package adapter

import "time"

type SessionSummary struct {
	ProjectName string
	LastActive  time.Time
	Summary     string
	RawLog      string
}

type AIAdapter interface {
	Name() string
	IsAvailable() bool
	Launch(projectPath string) error
	GetSessionSummary(projectPath string) (*SessionSummary, error)
	GetContextFile(projectPath string) (string, error)
}
