package executor

import (
	"os/exec"
	"regexp"
)

// ClaudeExecutor implements Executor for Claude Code CLI
type ClaudeExecutor struct{}

// BuildCommand builds the Claude Code command
func (e *ClaudeExecutor) BuildCommand(prompt, sessionID, imagePath, model string) []string {
	cmd := []string{"claude", "-p", prompt}

	if sessionID != "" {
		cmd = append(cmd, "--resume", sessionID)
	}

	if model != "" {
		cmd = append(cmd, "--model", model)
	}

	if imagePath != "" {
		// Claude Code appends file path at the end of arguments
		cmd = append(cmd, imagePath)
	}

	return cmd
}

// ParseSessionID extracts session ID from Claude Code output
func (e *ClaudeExecutor) ParseSessionID(output string) string {
	// Try various session ID formats
	patterns := []string{
		`session[:\s]+([a-zA-Z0-9-]+)`,
		`session\s*id[:\s]+([a-zA-Z0-9-]+)`,
		`resuming\s+session[:\s]+([a-zA-Z0-9-]+)`,
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
func (e *ClaudeExecutor) Name() string {
	return "claude"
}

// ListModels returns the list of available models
func (e *ClaudeExecutor) ListModels() ([]string, error) {
	// Return hardcoded list for Claude Code
	// Could parse from 'claude --help' in the future
	return []string{
		"claude-3-5-sonnet-latest",
		"claude-3-5-haiku-latest",
		"claude-3-opus-latest",
	}, nil
}

// Stats returns statistics information
func (e *ClaudeExecutor) Stats() (string, error) {
	// Claude Code has no stats API, just check installation
	_, err := exec.LookPath("claude")
	if err != nil {
		return "Claude Code is not installed", nil
	}
	return "Claude Code is installed", nil
}
