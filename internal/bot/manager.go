package bot

import (
	"context"
	"fmt"

	"github.com/mymmrac/telego"
	"telecode/internal/config"
)

// WorkspaceBot represents a single workspace with its bot instance
type WorkspaceBot struct {
	Config config.WorkspaceConfig
	Bot    *Bot
	TgBot  *telego.Bot
}

// Manager handles multiple workspace bots
type Manager struct {
	workspaces map[string]*WorkspaceBot
}

// NewManager creates a new multi-bot manager
func NewManager(cfg *config.Config) (*Manager, error) {
	mgr := &Manager{
		workspaces: make(map[string]*WorkspaceBot),
	}

	for _, wsConfig := range cfg.Workspaces {
		// Convert allowed chats to map
		allowedChats := make(map[int64]bool)
		for _, chatID := range wsConfig.AllowedChats {
			allowedChats[chatID] = true
		}

		// Create bot logic instance
		botLogic := NewBot(allowedChats, wsConfig.DefaultCLI)

		// Create Telegram bot
		var botOpts []telego.BotOption
		tgBot, err := telego.NewBot(wsConfig.BotToken, botOpts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create bot for workspace %s: %w", wsConfig.Name, err)
		}

		// Store workspace bot
		mgr.workspaces[wsConfig.Name] = &WorkspaceBot{
			Config: wsConfig,
			Bot:    botLogic,
			TgBot:  tgBot,
		}
	}

	return mgr, nil
}

// Start starts all workspace bots
func (m *Manager) Start(ctx context.Context) error {
	for name, ws := range m.workspaces {
		fmt.Printf("🤖 Starting bot for workspace: %s (dir: %s)\n", name, ws.Config.WorkingDir)

		// Start this workspace's bot in a goroutine
		go func(ws *WorkspaceBot) {
			if err := m.runWorkspaceBot(ctx, ws); err != nil {
				fmt.Printf("❌ Bot error for workspace %s: %v\n", ws.Config.Name, err)
			}
		}(ws)
	}

	return nil
}

// runWorkspaceBot runs a single workspace bot
func (m *Manager) runWorkspaceBot(ctx context.Context, ws *WorkspaceBot) error {
	// Get updates
	updates, err := ws.TgBot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start long polling: %w", err)
	}

	// Process updates
	for {
		select {
		case <-ctx.Done():
			return nil
		case update, ok := <-updates:
			if !ok {
				return nil
			}
			if err := m.handleUpdate(ctx, ws, update); err != nil {
				fmt.Printf("❌ Error handling update for %s: %v\n", ws.Config.Name, err)
			}
		}
	}
}

// handleUpdate handles a single update for a workspace bot
func (m *Manager) handleUpdate(ctx context.Context, ws *WorkspaceBot, update telego.Update) error {
	if update.Message == nil {
		return nil
	}

	chatID := update.Message.Chat.ID

	// Check if chat is allowed
	if !ws.Bot.IsAllowed(chatID) {
		return nil
	}

	// Check if message has photo
	if len(update.Message.Photo) > 0 {
		return m.handlePhotoMessage(ctx, ws, update.Message)
	}

	// Get command handler
	cmd := getCommandFromMessage(update.Message.Text)

	switch cmd {
	case "/new":
		return m.handleNewSession(ctx, ws, chatID)
	case "/status":
		return m.handleStatus(ctx, ws, chatID)
	case "/cli":
		return m.handleCLI(ctx, ws, chatID, update.Message.Text)
	case "/stats":
		return m.handleStats(ctx, ws, chatID)
	default:
		// Handle regular message
		return m.handleMessage(ctx, ws, chatID, update.Message.Text, "")
	}
}

func getCommandFromMessage(text string) string {
	if len(text) == 0 {
		return ""
	}
	if text[0] == '/' {
		// Extract command (up to first space or end)
		for i, c := range text {
			if c == ' ' {
				return text[:i]
			}
		}
		return text
	}
	return ""
}
