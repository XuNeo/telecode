package bot

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

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
		return fmt.Sprintf("Error: %v\n%s", err, string(output))
	}

	return string(output)
}
