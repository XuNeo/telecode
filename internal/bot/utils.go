package bot

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"time"
)

// stripAnsiCodes removes ANSI escape sequences and OSC sequences from text
func stripAnsiCodes(text string) string {
	// Remove ANSI escape sequences (color codes, cursor movements, etc.)
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	// Remove OSC sequences (terminal title changes, etc.)
	oscRegex := regexp.MustCompile(`\x1b\][0-9;]*.*?\x07`)
	// Remove any remaining escape sequences
	escapeRegex := regexp.MustCompile(`\x1b\[[\?0-9]*[hl]`)

	text = ansiRegex.ReplaceAllString(text, "")
	text = oscRegex.ReplaceAllString(text, "")
	text = escapeRegex.ReplaceAllString(text, "")
	return text
}

// runCommandWithDir executes a CLI command in a specific working directory
func runCommandWithDir(cmd []string, workingDir string) string {
	if len(cmd) == 0 {
		return "Error: Command is empty"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	command := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	command.Dir = workingDir // Set working directory
	output, err := command.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return "Error: Command execution timeout (5 minutes)"
	}

	if err != nil {
		return stripAnsiCodes(fmt.Sprintf("Error: %v\n%s", err, string(output)))
	}

	return stripAnsiCodes(string(output))
}
