package executor

// Executor defines the interface for CLI executors
type Executor interface {
	// BuildCommand builds the CLI command
	BuildCommand(prompt string, sessionID string, imagePath string, model string) []string

	// ParseSessionID extracts session ID from output
	ParseSessionID(output string) string

	// Name returns the CLI name
	Name() string

	// ListModels returns the list of available models
	ListModels() ([]string, error)

	// Stats returns statistics information
	Stats() (string, error)
}
