package adapter

type AIAdapter interface {
	Name() string
	IsAvailable() bool
	Launch(projectPath string) error
	GetContextFile(projectPath string) (string, error)
}
