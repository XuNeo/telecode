---
layout: default
title: Telecode
---

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

## Quick Start

### Installation

Install Telecode with a single command:

```bash
curl -sSL https://raw.githubusercontent.com/futureCreator/telecode/main/install.sh | bash
```

Or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/futureCreator/telecode/main/install.sh | bash
```

### Configuration

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

### Start the Server

```bash
# Using default config file (auto-detects ~/.telecode/config.yml)
telecode

# Specify custom config file
telecode -config /path/to/config.yml

# Show version
telecode -version
```

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

## Multi-Project Workflow

### Project A (Web Frontend)

1. Create a bot for your web project via @BotFather
2. Add to `telecode.yml`
3. Start telecode: `./telecode`
4. Chat with your web project bot in Telegram

### Project B (API Backend)

1. Create another bot for your backend project
2. Add to the same `telecode.yml`
3. Restart telecode
4. Chat with your backend bot

Each bot operates independently in its own working directory!

## Learn More

For detailed documentation, visit our [GitHub Repository](https://github.com/futureCreator/telecode).

## License

MIT License
