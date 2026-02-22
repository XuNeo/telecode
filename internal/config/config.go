package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// WorkspaceConfig represents a single workspace/bot configuration
type WorkspaceConfig struct {
	Name           string        `yaml:"name"`
	WorkingDir     string        `yaml:"working_dir"`
	BotToken       string        `yaml:"bot_token"`
	AllowedChats   []int64       `yaml:"allowed_chats,omitempty"`
	DefaultCLI     string        `yaml:"default_cli,omitempty"`
	CommandTimeout time.Duration `yaml:"command_timeout,omitempty"`
	Model          string        `yaml:"model,omitempty"`
}

// Config represents the complete telecode configuration
type Config struct {
	Workspaces []WorkspaceConfig `yaml:"workspaces"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults and validate
	for i := range cfg.Workspaces {
		if cfg.Workspaces[i].DefaultCLI == "" {
			cfg.Workspaces[i].DefaultCLI = "claude"
		}
		if cfg.Workspaces[i].CommandTimeout == 0 {
			cfg.Workspaces[i].CommandTimeout = 20 * time.Minute
		}
		if cfg.Workspaces[i].WorkingDir == "" {
			return nil, fmt.Errorf("workspace %d: working_dir is required", i)
		}
		if cfg.Workspaces[i].BotToken == "" {
			return nil, fmt.Errorf("workspace %d: bot_token is required", i)
		}
	}

	return &cfg, nil
}

// GetDefaultConfigPath returns the default configuration file path
func GetDefaultConfigPath() string {
	// Check for config in home directory
	home, err := os.UserHomeDir()
	if err == nil {
		configPath := home + "/.telecode/config.yml"
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	// Check for config in current directory
	if _, err := os.Stat("telecode.yml"); err == nil {
		return "telecode.yml"
	}

	// Check for config in /etc
	if _, err := os.Stat("/etc/telecode/config.yml"); err == nil {
		return "/etc/telecode/config.yml"
	}

	return ""
}

// CreateExampleConfig creates an example configuration file
func CreateExampleConfig(path string) error {
	example := `# Telecode Multi-Bot Configuration
# Each workspace represents a separate project with its own bot

workspaces:
  - name: project-a
    working_dir: /home/user/project-a
    bot_token: "YOUR_BOT_TOKEN_1"
    allowed_chats:
      - 123456789
    default_cli: opencode
    command_timeout: 20m
    # model: anthropic/opus-4.6  # Optional: OpenCode model (defaults to opus-4.6)

  - name: project-b
    working_dir: /home/user/project-b
    bot_token: "YOUR_BOT_TOKEN_2"
    allowed_chats:
      - 987654321
    default_cli: claude
    # command_timeout defaults to 20m if not specified
`
	return os.WriteFile(path, []byte(example), 0644)
}
