package executor

import (
	"os/exec"
	"regexp"
)

// OpenCodeExecutor implements Executor for OpenCode CLI
type OpenCodeExecutor struct{}

// BuildCommand builds the OpenCode command
func (e *OpenCodeExecutor) BuildCommand(prompt, sessionID, imagePath string, model string) []string {
	// Use default model if not specified
	if model == "" {
		model = "anthropic/opus-4.6"
	}
	cmd := []string{"opencode", "run", "--format", "json", "--model", model, prompt}

	if sessionID != "" {
		cmd = append(cmd, "--session", sessionID)
	}

	if imagePath != "" {
		cmd = append(cmd, "--file", imagePath)
	}

	return cmd
}

// ParseSessionID extracts session ID from OpenCode JSON output
func (e *OpenCodeExecutor) ParseSessionID(output string) string {
	// Parse sessionID from JSON output (e.g., {"type":"step_start","sessionID":"ses_xxx",...})
	// Look for "sessionID":"ses_..." pattern in the first JSON line
	re := regexp.MustCompile(`"sessionID":"(ses_[a-zA-Z0-9_-]+)"`)
	match := re.FindStringSubmatch(output)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// Name returns the Executor name
func (e *OpenCodeExecutor) Name() string {
	return "opencode"
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
