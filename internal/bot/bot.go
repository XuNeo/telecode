package bot

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"telecode/internal/executor"
	"telecode/internal/session"
)

// ChatSettings stores per-chat configuration
type ChatSettings struct {
	CLI   string `json:"cli"`
	Model string `json:"model"`
}

// Bot handles the core logic of the Telegram bot
type Bot struct {
	sessionMgr   *session.Manager
	chatSettings map[int64]ChatSettings
	settingsMu   sync.RWMutex
	allowedChats map[int64]bool
	executors    map[string]executor.Executor
	defaultCLI   string
}

// NewBot creates a new bot instance
func NewBot(allowedChats map[int64]bool, defaultCLI string) *Bot {
	return &Bot{
		sessionMgr:   session.NewManager(),
		chatSettings: make(map[int64]ChatSettings),
		allowedChats: allowedChats,
		executors: map[string]executor.Executor{
			"claude":   &executor.ClaudeExecutor{},
			"opencode": &executor.OpenCodeExecutor{},
		},
		defaultCLI: defaultCLI,
	}
}

// IsAllowed checks if the chat_id is in the allowlist
func (b *Bot) IsAllowed(chatID int64) bool {
	return b.allowedChats[chatID]
}

// GetCLI returns the CLI setting for a chat
func (b *Bot) GetCLI(chatID int64) string {
	b.settingsMu.RLock()
	defer b.settingsMu.RUnlock()
	if cli := b.chatSettings[chatID].CLI; cli != "" {
		return cli
	}
	return b.defaultCLI
}

// SetCLI sets the CLI for a chat
func (b *Bot) SetCLI(chatID int64, cli string) error {
	// Check if CLI exists
	if _, err := exec.LookPath(cli); err != nil {
		return fmt.Errorf("CLI '%s' is not installed", cli)
	}

	b.settingsMu.Lock()
	defer b.settingsMu.Unlock()
	settings := b.chatSettings[chatID]
	settings.CLI = cli
	b.chatSettings[chatID] = settings

	// Reset session when CLI changes
	b.sessionMgr.Delete(chatID)

	return nil
}

// GetModel returns the model setting for a chat
func (b *Bot) GetModel(chatID int64) string {
	b.settingsMu.RLock()
	defer b.settingsMu.RUnlock()
	return b.chatSettings[chatID].Model
}

// SetModel sets the model for a chat
func (b *Bot) SetModel(chatID int64, model string) {
	b.settingsMu.Lock()
	defer b.settingsMu.Unlock()
	settings := b.chatSettings[chatID]
	settings.Model = model
	b.chatSettings[chatID] = settings
}

// GetSessionID returns the session ID for a chat
func (b *Bot) GetSessionID(chatID int64) string {
	return b.sessionMgr.Get(chatID)
}

// NewSession starts a new session
func (b *Bot) NewSession(chatID int64) {
	b.sessionMgr.Delete(chatID)
}

// UpdateSessionFromOutput extracts and saves session ID from output
func (b *Bot) UpdateSessionFromOutput(chatID int64, cli, output string) {
	if b.sessionMgr.Exists(chatID) {
		return // Don't update if session already exists
	}

	exec := b.executors[cli]
	if exec == nil {
		return
	}

	if sessionID := exec.ParseSessionID(output); sessionID != "" {
		b.sessionMgr.Set(chatID, sessionID)
	}
}

// GetExecutor returns the Executor for a CLI name
func (b *Bot) GetExecutor(cli string) executor.Executor {
	return b.executors[cli]
}

// BuildCommand builds the CLI command
func (b *Bot) BuildCommand(chatID int64, prompt, imagePath string) []string {
	cli := b.GetCLI(chatID)
	sessionID := b.GetSessionID(chatID)
	model := b.GetModel(chatID)

	exec := b.executors[cli]
	if exec == nil {
		return nil
	}

	return exec.BuildCommand(prompt, sessionID, imagePath, model)
}

// ListModels returns the list of available models for current CLI
func (b *Bot) ListModels(chatID int64) ([]string, error) {
	cli := b.GetCLI(chatID)
	exec := b.executors[cli]
	if exec == nil {
		return nil, fmt.Errorf("unsupported CLI: %s", cli)
	}
	return exec.ListModels()
}

// GetStats returns statistics for current CLI
func (b *Bot) GetStats(chatID int64) (string, error) {
	cli := b.GetCLI(chatID)
	exec := b.executors[cli]
	if exec == nil {
		return "", fmt.Errorf("unsupported CLI: %s", cli)
	}
	return exec.Stats()
}

// GetStatus returns the current status
func (b *Bot) GetStatus(chatID int64) (cli, sessionID, model string) {
	cli = b.GetCLI(chatID)
	sessionID = b.GetSessionID(chatID)
	model = b.GetModel(chatID)

	if sessionID == "" {
		sessionID = "none"
	}
	if model == "" {
		model = "default"
	}

	return
}

// ParseAllowedChats parses the allowed chat_id list from environment variable
func ParseAllowedChats(env string) map[int64]bool {
	result := make(map[int64]bool)
	if env == "" {
		return result
	}

	for _, idStr := range strings.Split(env, ",") {
		idStr = strings.TrimSpace(idStr)
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			result[id] = true
		}
	}
	return result
}
