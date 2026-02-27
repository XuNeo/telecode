# AGENTS.md — Telecode

> ⚠️ **CRITICAL RULE: NEVER push to remote without explicit user permission.**
> 
> All commits and changes must be reviewed by the user before pushing to any remote repository.
> The agent must always ask for permission before running `git push` or any command that uploads code to a remote.

## Project Overview

## Project Overview

Telecode is a multi-bot Telegram server for remotely using AI coding assistants (Claude Code, OpenCode).
It is a small, focused Go binary (~10 source files) with no test suite yet.

**Module**: `telecode` (see `go.mod`)
**Go version**: 1.25.5
**Key dependency**: `github.com/mymmrac/telego` (Telegram Bot API)

## Project Structure

```
cmd/telecode/main.go          # Entry point: flag parsing, config loading, signal handling
internal/
  config/config.go             # YAML config loading, defaults, validation
  bot/
    bot.go                     # Core bot logic: session mgmt, CLI selection, command building
    manager.go                 # Multi-bot orchestrator: creates WorkspaceBot per config entry
    handlers.go                # Telegram command/message handlers (/new, /cli, /status, /stats)
    utils.go                   # Command execution, ANSI stripping, OpenCode JSON parsing
  executor/
    executor.go                # Executor interface (BuildCommand, ParseSessionID, Name, Stats)
    claude.go                  # Claude Code CLI executor
    opencode.go                # OpenCode CLI executor
  session/
    manager.go                 # In-memory session ID store (map + sync.RWMutex)
```

### Package Dependency Flow

```
main → config, bot
bot  → config, executor, session
executor → (no internal deps)
session  → (no internal deps)
```

## Build / Run / Test Commands

```bash
# Build (statically linked, CGO_ENABLED=0)
make build

# Run
make run                # build + run
make run-config         # build + run with telecode.yml

# Test
make test               # go test -v ./...
go test -v ./internal/bot/...          # single package
go test -v -run TestFoo ./internal/... # single test by name

# Lint & Format
make fmt                # go fmt ./...
make lint               # golangci-lint run (if installed)
make tidy               # go mod tidy

# Cross-compile
make cross-build        # linux/darwin/windows amd64+arm64

# Clean
make clean
```

There is **no golangci-lint config** (`.golangci.yml`) — the lint target runs the default ruleset.
There are **no tests yet** (`*_test.go` files do not exist). When adding tests, follow stdlib conventions.
There is **no CI/CD config** (`.github/workflows/` does not exist).

## Code Style Guidelines

### Imports

Two groups separated by a blank line: stdlib, then everything else (third-party and internal mixed).

```go
import (
    "context"
    "fmt"
    "os"

    "github.com/mymmrac/telego"
    "telecode/internal/config"
)
```

Aliasing: used for telego utilities only → `tu "github.com/mymmrac/telego/telegoutil"`.

### Error Handling

- **Wrap with `fmt.Errorf` + `%w`** for propagated errors:
  ```go
  return nil, fmt.Errorf("failed to read config file: %w", err)
  ```
- **No custom error types** — all errors are wrapped stdlib errors.
- **No sentinel errors** — errors are checked with `err != nil` only.
- In `main()`, errors print with `fmt.Printf` and `os.Exit(1)` — no `log.Fatal`.
- Handler functions return `error` up the chain; the manager logs them with `fmt.Printf`.

### Naming Conventions

- **Packages**: short, lowercase, single word (`bot`, `config`, `executor`, `session`).
- **Exported types**: PascalCase nouns — `Manager`, `Bot`, `WorkspaceBot`, `Executor`.
- **Constructors**: `NewX()` pattern — `NewBot()`, `NewManager()`.
- **Interfaces**: verb-based — `Executor` (not `IExecutor` or `ExecutorInterface`).
- **Unexported fields**: camelCase — `sessionMgr`, `chatSettings`, `settingsMu`.
- **Acronyms**: mixed — `CLI` stays uppercase in exported names (`DefaultCLI`, `GetCLI`), `ID` uppercase (`chatID`, `sessionID` in params, `ParseSessionID` in methods).
- **Receiver names**: single letter matching type — `(m *Manager)`, `(b *Bot)`, `(e *ClaudeExecutor)`.

### Function Signatures

- `context.Context` is the first parameter in handler functions.
- No functional options pattern — config is passed via structs or direct args.
- Return `error` as the last return value.
- Multiple returns use named returns only for simple status functions:
  ```go
  func (b *Bot) GetStatus(chatID int64) (cli, sessionID string)
  ```

### Structs

- YAML tags on config structs: `` `yaml:"field_name"` `` with `omitempty` for optional fields.
- JSON tags on internal structs: `` `json:"field"` ``.
- No struct embedding — composition via explicit fields.
- Constructor functions return pointer types: `func NewBot(...) *Bot`.

### Concurrency

- `sync.RWMutex` protects shared maps (`sessions`, `chatSettings`).
- Pattern: `mu.RLock()` / `defer mu.RUnlock()` for reads, `mu.Lock()` / `defer mu.Unlock()` for writes.
- Goroutines launched with `go func(param Type) { ... }(capturedVar)` to avoid closure capture bugs.
- `context.WithCancel` / `context.WithTimeout` for lifecycle control.

### Logging

- **No logging library** — all output uses `fmt.Printf` / `fmt.Println` with emoji prefixes.
- Error: `fmt.Printf("❌ ...")`, success: `fmt.Printf("✅ ...")`, info: `fmt.Printf("🤖 ...")`.
- This is intentional for a lightweight CLI tool — do not introduce a logging framework.

### Comments

- Godoc-style on all exported types and functions: `// TypeName does X`.
- Brief inline comments for non-obvious logic.
- No package-level doc comments (files start with `package name`).

### Telegram Message Patterns

- Use `tu.Message(tu.ID(chatID), text)` for building messages.
- Markdown formatting: `.WithParseMode(telego.ModeMarkdown)`.
- Non-critical send errors ignored with `_, _ = ws.TgBot.SendMessage(...)`.
- Critical send errors returned: `_, err := ws.TgBot.SendMessage(...); return err`.

## Architecture Notes

- **No dependency injection framework** — manual wiring in `NewManager()`.
- **Config search order**: `-config` flag → `./telecode.yml` → `~/.telecode/config.yml` → `/etc/telecode/config.yml`.
- **Version**: injected via ldflags at build time (`-X main.version=$(VERSION)`).
- **Static linking**: all builds use `CGO_ENABLED=0` with `-extldflags '-static'`.
- **Long polling**: bots use `UpdatesViaLongPolling`, not webhooks.

## Adding New Features

### Adding a new CLI executor
1. Create `internal/executor/newcli.go` implementing the `Executor` interface.
2. Register it in `bot.NewBot()` executors map.
3. Update the CLI validation in `handlers.go` `handleCLI()`.

### Adding a new bot command
1. Add handler method on `Manager` in `handlers.go`.
2. Add case in `handleUpdate()` switch in `manager.go`.

### Adding tests
- No test infrastructure exists. Use stdlib `testing` package.
- Place `*_test.go` next to the file being tested.
- Prefer table-driven tests with `t.Run()` subtests.
- Run with `go test -v ./internal/...`.
