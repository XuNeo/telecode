# Telecode - Telegram Coding Agent Bot

A multi-bot server for remotely using AI coding assistants (Claude Code, OpenCode) via Telegram. Supports a configuration file-based structure for managing multiple projects simultaneously.

## Features

- 🚀 **Lightweight**: Single binary execution (statically linked)
- 💰 **Cost-effective**: Only token costs (no hosting fees)
- 🔒 **Secure**: Allowlist-based access control
- 💬 **Interactive Sessions**: Per-chat_id session persistence
- 🖼️ **Image Support**: Analyze Telegram images
- 🔄 **Multi-CLI**: Choose between Claude Code and OpenCode
- 🏗️ **Multi-Bot**: Manage multiple projects with separate bots
- 📁 **Project Isolation**: Each bot works in its own working directory
- ⏱️ **Configurable Timeout**: Set command execution timeout per workspace
- 📊 **Smart Output**: Automatic JSON parsing for OpenCode responses

## Installation

### Requirements

- Go 1.25.5 or higher
- Telegram Bot API token (from @BotFather)
- Claude Code or OpenCode CLI installed

### Quick Install (Recommended)

Install Telecode with a single command:

```bash
curl -sSL https://raw.githubusercontent.com/futureCreator/telecode/main/install.sh | bash
```

Or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/futureCreator/telecode/main/install.sh | bash
```

This will:
- Download the latest statically-linked binary from GitHub
- Install to `~/.local/bin`
- Create config at `~/.telecode/config.yml`
- Warn if `~/.local/bin` is not in your PATH

Then edit the config and run:

```bash
# Edit config
nano ~/.telecode/config.yml

# Run telecode
telecode
```

### Manual Installation

#### Build Binary (Statically Linked)

```bash
# macOS
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -extldflags '-static'" -o telecode-darwin-amd64 ./cmd/telecode

# Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -extldflags '-static'" -o telecode-linux-amd64 ./cmd/telecode

# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -extldflags '-static'" -o telecode-windows-amd64.exe ./cmd/telecode
```

#### Build Locally (using Makefile)

```bash
# Build for current platform (statically linked)
make build

# Or build manually
go build -o telecode ./cmd/telecode
```

#### Install with Local Script

After building, use the local install script:

```bash
# Build first
make build

# Install locally
./install.sh --local
```

#### Manual Binary Install

```bash
# Create directories
mkdir -p ~/.local/bin ~/.telecode

# Copy binary
cp telecode ~/.local/bin/

# Copy config
cp telecode.yml ~/.telecode/config.yml

# Make sure ~/.local/bin is in your PATH
export PATH="$HOME/.local/bin:$PATH"
```

## Configuration

### Configuration File

Generate an example configuration file:

```bash
./telecode -generate-config
```

Edit the generated `telecode.yml`:

```yaml
workspaces:
  - name: project-a
    working_dir: /home/user/project-a
    bot_token: "YOUR_BOT_TOKEN_1"
    allowed_chats:
      - 123456789
    default_cli: opencode

  - name: project-b
    working_dir: /home/user/project-b
    bot_token: "YOUR_BOT_TOKEN_2"
    allowed_chats:
      - 987654321
    default_cli: claude
```

| Configuration | Description | Required | Default |
|--------------|-------------|----------|---------|
| `name` | Workspace name | ✅ | - |
| `working_dir` | Directory where CLI executes | ✅ | - |
| `bot_token` | Telegram Bot API token | ✅ | - |
| `allowed_chats` | List of allowed chat_ids | ❌ | All blocked |
| `default_cli` | Default CLI (claude/opencode) | ❌ | `claude` |
| `model` | OpenCode model (provider/model format) | ❌ | `anthropic/opus-4.6` |
| `command_timeout` | Command execution timeout | ❌ | `20m` |

### CLI API Keys

Claude Code and OpenCode manage their own API keys, no additional configuration needed.

## Usage

### Command Line Flags

| Flag | Description | Example |
|------|-------------|---------|
| `-config` | Path to configuration file | `-config /path/to/config.yml` |
| `-generate-config` | Generate example configuration file | `-generate-config` |
| `-version` | Show version information | `-version` |

### Start the Server

```bash
# Using default config file (auto-detects ~/.telecode/config.yml)
telecode

# Specify custom config file
telecode -config /path/to/config.yml

# Show version
telecode -version
```

### Configuration File Locations

Telecode searches for config files in this order:

1. Path specified by `-config` flag
2. `./telecode.yml` (current directory)
3. `~/.telecode/config.yml` (home directory) ← **Default when using install.sh**
4. `/etc/telecode/config.yml` (system-wide)

## Bot Commands

All bots support the following commands:

| Command | Function |
|---------|----------|
| `/new` | Start new session (reset context) |
| `/cli` | Show current CLI |
| `/cli claude` | Switch to Claude Code |
| `/cli opencode` | Switch to OpenCode |
| `/status` | Show current status (workspace, CLI, session) |
| `/stats` | Show token usage statistics |

### Regular Messages

Simply send a message to interact with the AI:

```
Review this code for bugs
Refactor this function
Explain how this works
```

### Image Analysis

Send a photo with a caption to analyze it:

```
[Photo with caption: "What's wrong with this error?"]
```

If no caption is provided, it defaults to "Analyze this image".

### OpenCode JSON Output

When using OpenCode CLI, responses are automatically parsed from JSON format, providing clean, readable output in Telegram.

## Multi-Project Workflow Example

### Project A (Web Frontend)

1. Create a bot for your web project via @BotFather
2. Add to `telecode.yml`:
   ```yaml
   - name: web-frontend
     working_dir: /home/user/projects/web-app
     bot_token: "123456:ABC..."
     allowed_chats: [YOUR_CHAT_ID]
     default_cli: opencode
   ```
3. Start telecode: `./telecode`
4. Chat with your web project bot in Telegram:
   ```
   /new
   Fix the React rendering issue in this component
   ```

### Project B (API Backend)

1. Create another bot for your backend project
2. Add to the same `telecode.yml`:
   ```yaml
   - name: api-backend
     working_dir: /home/user/projects/api-server
     bot_token: "789012:XYZ..."
     allowed_chats: [YOUR_CHAT_ID]
     default_cli: claude
   ```
3. Restart telecode (or it will auto-reload)
4. Chat with your backend bot:
   ```
   /new
   Add error handling to this API endpoint
   ```

Each bot operates independently in its own working directory!

## Deployment

### macOS (launchd)

Create `~/Library/LaunchAgents/com.telecode.bot.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.telecode.bot</string>
    <key>ProgramArguments</key>
    <array>
        <string>~/.local/bin/telecode</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>WorkingDirectory</key>
    <string>/home/user</string>
    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/usr/local/bin:/usr/bin:/bin:~/.local/bin</string>
    </dict>
</dict>
</plist>
```

```bash
launchctl load ~/Library/LaunchAgents/com.telecode.bot.plist
```

Note: When using the install script, telecode will automatically find the config at `~/.telecode/config.yml`.

### Linux (systemd) - User Service

Create `~/.config/systemd/user/telecode.service`:

```ini
[Unit]
Description=Telecode Multi-Bot Server
After=network.target

[Service]
Type=simple
ExecStart=%h/.local/bin/telecode
Restart=always

[Install]
WantedBy=default.target
```

```bash
# Enable and start user service
systemctl --user enable telecode
systemctl --user start telecode

# Or for system-wide service (root required)
sudo tee /etc/systemd/system/telecode.service << EOF
[Unit]
Description=Telecode Multi-Bot Server
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/telecode -config /etc/telecode/config.yml
Restart=always
User=telecode
Group=telecode

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl enable telecode
sudo systemctl start telecode
```

## Development

### Makefile Commands

The project includes a comprehensive Makefile for development tasks:

```bash
# Build commands
make build              # Build for current platform (statically linked)
make cross-build        # Build for all platforms (Linux, macOS, Windows)
make build-race         # Build with race detector (for debugging)

# Installation
make install            # Install to GOPATH/bin
make install-system     # Install to /usr/local/bin (requires sudo)

# Development
make run                # Build and run
make run-config         # Run with telecode.yml config
make dev                # Run with hot reload (requires air)

# Code quality
make test               # Run tests
make test-coverage      # Run tests with coverage report
make fmt                # Format code with go fmt
make lint               # Run golangci-lint
make tidy               # Tidy go modules

# Utilities
make generate-config    # Generate example configuration
make verify-static      # Verify binary is statically linked
make clean              # Remove build artifacts
make help               # Show all available commands
```

### Static Linking

All binaries are statically linked for maximum portability:

```bash
# Build static binary
make build

# Verify static linking
make verify-static
```

## Project Structure

```
telecode/
├── cmd/telecode/
│   └── main.go              # Entry point with multi-bot support
├── internal/
│   ├── executor/
│   │   ├── executor.go      # Executor interface
│   │   ├── claude.go        # Claude Code implementation
│   │   └── opencode.go      # OpenCode implementation
│   ├── bot/
│   │   ├── bot.go           # Single bot logic
│   │   ├── manager.go       # Multi-bot manager
│   │   ├── handlers.go      # Telegram message handlers
│   │   └── utils.go         # Utility functions
│   ├── session/
│   │   └── manager.go       # Session management
│   └── config/
│       └── config.go        # Configuration file handling
├── install.sh               # Installation script
├── go.mod
├── go.sum
├── README.md
└── telecode.yml             # Example configuration
```

## Security

- Bot tokens and allowlists are managed via configuration files
- Messages from unauthorized chat_ids are silently ignored
- CLI executables are verified before execution
- Working directories are isolated per workspace
- Configuration files should have restricted permissions (chmod 600)

## Troubleshooting

### Check Version

```bash
telecode -version
```

### Bot not responding

1. Check bot token is correct
2. Verify your chat_id is in `allowed_chats`
3. Ensure you've started the bot with `/start` in Telegram
4. Check logs for errors

### CLI not found

```bash
# Verify CLI is installed and in PATH
which opencode
which claude

# Install if missing
npm install -g opencode
# or
npm install -g @anthropic-ai/claude-code
```

### Configuration file not found

```bash
# Generate example config
./telecode -generate-config

# Specify config path explicitly
./telecode -config /path/to/your/config.yml
```

### Permission denied on working directory

Ensure the user running telecode has read/write access:

```bash
chown -R $(whoami) /path/to/project
chmod 755 /path/to/project
```

### Command timeout issues

If commands are timing out, increase the timeout in your config:

```yaml
workspaces:
  - name: my-project
    working_dir: /path/to/project
    bot_token: "YOUR_TOKEN"
    command_timeout: 30m  # Increase from default 20m
```

## License

MIT
