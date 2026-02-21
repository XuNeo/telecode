package executor

import (
	"os/exec"
	"regexp"
)

// OpenCodeExecutor implements Executor for OpenCode CLI
type OpenCodeExecutor struct{}

// BuildCommand builds the OpenCode command
func (e *OpenCodeExecutor) BuildCommand(prompt, sessionID, imagePath, model string) []string {
	cmd := []string{"opencode", "run", prompt}

	if sessionID != "" {
		cmd = append(cmd, "--session", sessionID)
	}

	if model != "" {
		cmd = append(cmd, "--model", model)
	}

	if imagePath != "" {
		cmd = append(cmd, "--file", imagePath)
	}

	return cmd
}

// ParseSessionID extracts session ID from OpenCode output
func (e *OpenCodeExecutor) ParseSessionID(output string) string {
	// Parse OpenCode session ID format
	// Examples:
	//   "Continue  opencode -s ses_37f9659a6ffemnd5vvn1GC2Y5Q"
	//   "session: ses_abc123"
	patterns := []string{
		`-s\s+([a-zA-Z0-9_-]+)`,
		`session[:\s]+([a-zA-Z0-9_-]+)`,
		`session\s*id[:\s]+([a-zA-Z0-9_-]+)`,
		`--session\s+([a-zA-Z0-9_-]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(output)
		if len(match) > 1 {
			return match[1]
		}
	}

	return ""
}

// Name returns the Executor name
func (e *OpenCodeExecutor) Name() string {
	return "opencode"
}

// ListModels returns the list of available models
func (e *OpenCodeExecutor) ListModels() ([]string, error) {
	// Try to get models from OpenCode command
	cmd := exec.Command("opencode", "models")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Return default list on failure
		return []string{
			"claude-3-5-sonnet",
			"claude-3-5-haiku",
			"gpt-4",
			"gpt-4o",
		}, nil
	}

	// TODO: Parse output to extract model list
	_ = output
	return []string{"Model list query coming soon"}, nil
}

// Stats returns statistics information
func (e *OpenCodeExecutor) Stats() (string, error) {
	// Try to get OpenCode stats
	cmd := exec.Command("opencode", "stats")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "Could not retrieve statistics", nil
	}
	return string(output), nil
}
