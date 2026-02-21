# Telecode - 텔레그램용 코딩 에이전트 봇

## 컨셉

텔레그램에서 Claude Code 또는 OpenCode CLI를 직접 호출하는 경량 봇. OpenClaw를 거치지 않아 토큰 비용 이중 과징 없음.

## 왜 만드는가

**문제**: OpenClaw를 통해 코딩 에이전트를 사용하면 호스팅 비용 + 토큰 비용이 이중으로 발생

**해결**: 텔레그램 봇이 직접 CLI를 실행하면 토큰 비용만 소모

**핵심 목표**: 텔레그램에서 **대화형으로** 코딩 에이전트를 사용하는 경험 제공. 터미널에서 직접 쓰는 것과 동일한 컨텍스트 유지

## 기술 스택

| 구성요소 | 기술 | 비고 |
|---------|------|------|
| 봇 서버 | **Go** | 단일 바이너리, 크로스 플랫폼 (macOS/Linux/Windows) |
| 메시징 | Telegram Bot API | Long Polling 방식 |
| 실행 엔진 | Claude Code / OpenCode | Executor 추상화로 교체 가능 |
| 배포 | 로컬 실행 | macOS: launchd, Linux: systemd, Windows: 서비스 |

## 설정

환경변수 사용. 설정 파일 관리 오버헤드 없음.

| 환경변수 | 설명 | 예시 |
|----------|------|------|
| `TELECODE_BOT_TOKEN` | Telegram Bot API 토큰 (BotFather 발급) | `123456:ABC-DEF...` |
| `TELECODE_ALLOWED_CHATS` | 허용된 chat_id (쉼표 구분) | `5788362055` 또는 `5788362055,123456789` |
| `TELECODE_DEFAULT_CLI` | 기본 CLI (claude 또는 opencode) | `claude` (기본값) |

**CLI API 키**는 각 CLI가 자체 관리하므로 봇에서 신경 쓸 필요 없음.

## 크로스 플랫폼

```bash
# macOS
GOOS=darwin GOARCH=amd64 go build -o telecode-darwin-amd64 ./cmd/telecode

# Linux
GOOS=linux GOARCH=amd64 go build -o telecode-linux-amd64 ./cmd/telecode

# Windows
GOOS=windows GOARCH=amd64 go build -o telecode-windows-amd64.exe ./cmd/telecode
```

**배포 방식**

| 플랫폼 | 데몬 등록 | 비고 |
|--------|----------|------|
| **macOS (메인)** | launchd | `~/Library/LaunchAgents/` |
| Linux | systemd | `/etc/systemd/system/` |
| Windows | 서비스 또는 작업 스케줄러 | |

## 동작 흐름

```
┌─────────────┐    Long Polling     ┌──────────────┐
│  텔레그램    │ ◄────────────────── │  봇 서버     │
│             │                     │              │
│  메시지 전송 │ ─────────────────► │  메시지 수신  │
│             │                     │      ↓       │
│             │                     │ CLI 선택     │
│             │                     │ (chatSettings)│
│             │                     │      ↓       │
│             │                     │ 세션 조회    │
│             │                     │ (chat_id →   │
│             │                     │  session_id) │
│             │                     │      ↓       │
│             │                     │ CLI 실행     │
│             │                     │ (--resume)   │
│             │                     │      ↓       │
│             │                     │ 세션 ID 저장  │
│             │                     │ (새 세션인 경우)│
│             │                     │      ↓       │
│  응답 수신   │ ◄────────────────── │  결과 전송   │
└─────────────┘                     └──────────────┘
```

### 상세 흐름

1. **메시지 수신**: 봇 서버가 `getUpdates` 폴링으로 새 메시지 감지
2. **명령어 처리**: `/new`면 세션 초기화, `/cli`면 CLI 전환
3. **CLI 실행**:
   - Claude Code: `claude -p "<프롬프트>" --resume <session_id>`
   - OpenCode: `opencode run "<프롬프트>" --session <session_id>`
4. **세션 ID 저장**: 새 세션이면 출력에서 session_id 추출 후 매핑 저장
5. **결과 전송**: 긴 메시지는 분할하여 `sendMessage` API로 전송

## 핵심 기능

### MVP (필수)

- [x] 텔레그램 메시지 수신 (Long Polling)
- [x] CLI 실행 및 결과 캡처
- [x] 텔레그램으로 응답 전송
- [x] **대화형 세션 유지** — chat_id별 session_id 매핑
- [x] **이미지 입력 지원** — 텔레그램 이미지 → CLI `--file` 전달
- [x] **접근 제어 (Allowlist)** — 허용된 chat_id만 사용 가능
- [x] **멀티 CLI 지원** — Claude Code / OpenCode 선택 가능
- [x] 기본 에러 처리

### v1.0 (봇 명령어)

| 명령어 | 기능 | 비고 |
|--------|------|------|
| `/new` | 새 세션 시작 | 세션 ID 초기화 |
| `/cli` | 현재 CLI 조회 | - |
| `/cli claude` | Claude Code로 전환 | - |
| `/cli opencode` | OpenCode로 전환 | - |
| `/model [name]` | 모델 조회/변경 | `--model` 플래그 |
| `/models` | 사용 가능한 모델 목록 | CLI별 조회 |
| `/stats` | 토큰 사용량/비용 통계 | CLI별 조회 |
| `/status` | 현재 세션/CLI/모델 요약 | 내부 상태 조회 |

### v1.1

- [ ] 스트리밍 응답 (긴 출력 분할 전송)
- [ ] 타임아웃 처리
- [ ] 로깅
- [ ] 세션 만료 처리

### Future

- [ ] 멀티 유저 지원 (chat_id별 세션 격리)
- [ ] 워크스페이스 프리셋 (`/workspace`)
- [ ] Agent 지원 (`/agent`, `/agents`)

## 비용

| 항목 | 비용 |
|------|------|
| 텔레그램 봇 API | 무료 |
| 봇 서버 호스팅 | 무료 (로컬 실행) |
| Claude Code / OpenCode 토큰 | 사용량 기반 |

**총 비용**: CLI 토큰만

## 아키텍처

### 프로젝트 구조

```
telecode/
├── cmd/
│   └── telecode/
│       └── main.go
├── internal/
│   ├── executor/
│   │   ├── executor.go       # 인터페이스 정의
│   │   ├── claude.go         # Claude Code 구현
│   │   └── opencode.go       # OpenCode 구현
│   ├── bot/
│   │   └── bot.go            # 텔레그램 봇 로직
│   └── session/
│       └── manager.go        # 세션 관리
├── go.mod
├── go.sum
└── README.md
```

### Executor 인터페이스

```go
package executor

type Executor interface {
    // 명령어 빌드
    BuildCommand(prompt string, sessionID string, imagePath string, model string) []string
    
    // 세션 ID 파싱
    ParseSessionID(output string) string
    
    // CLI 이름
    Name() string
    
    // 모델 목록 조회
    ListModels() ([]string, error)
    
    // 통계 조회
    Stats() (string, error)
}
```

### Claude Code Executor

```go
type ClaudeExecutor struct{}

func (e *ClaudeExecutor) BuildCommand(prompt, sessionID, imagePath, model string) []string {
    cmd := []string{"claude", "-p", prompt}
    
    if sessionID != "" {
        cmd = append(cmd, "--resume", sessionID)
    }
    
    if model != "" {
        cmd = append(cmd, "--model", model)
    }
    
    if imagePath != "" {
        cmd = append(cmd, imagePath) // Claude Code는 인자 끝에 파일 경로
    }
    
    return cmd
}

func (e *ClaudeExecutor) ParseSessionID(output string) string {
    // Claude Code 세션 ID 형식 파싱
    re := regexp.MustCompile(`session[:\s]+([a-zA-Z0-9-]+)`)
    match := re.FindStringSubmatch(output)
    if len(match) > 1 {
        return match[1]
    }
    return ""
}
```

### OpenCode Executor

```go
type OpenCodeExecutor struct{}

func (e *OpenCodeExecutor) BuildCommand(prompt, sessionID, imagePath, model string) []string {
    cmd := []string{"opencode", "run", prompt}
    
    if sessionID != "" {
        cmd = append(cmd, "--session", sessionID)
    }
    
    if model != "" {
        cmd = append(cmd, "--model", model)
    }
    
    if imagePath != "" {
        cmd = append(cmd, "--file", imagePath)
    }
    
    return cmd
}
```

## 구현 예시

### main.go

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "regexp"
    "strconv"
    "strings"
    "sync"
    "time"
    
    tb "gopkg.in/telebot.v3"
    "telecode/internal/executor"
)

type ChatSettings struct {
    CLI   string `json:"cli"`
    Model string `json:"model"`
}

var (
    sessions     = make(map[int64]string)
    chatSettings = make(map[int64]ChatSettings)
    mu           sync.RWMutex
    allowedChats map[int64]bool
    
    executors = map[string]executor.Executor{
        "claude":   &executor.ClaudeExecutor{},
        "opencode": &executor.OpenCodeExecutor{},
    }
)

func isAllowed(chatID int64) bool {
    return allowedChats[chatID]
}

func parseAllowedChats() map[int64]bool {
    result := make(map[int64]bool)
    env := os.Getenv("TELECODE_ALLOWED_CHATS")
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

func getDefaultCLI() string {
    cli := os.Getenv("TELECODE_DEFAULT_CLI")
    if cli == "" {
        return "claude"
    }
    return cli
}

func runCommand(cmd []string) string {
    out, _ := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
    return string(out)
}

func chunkString(s string, size int) []string {
    var chunks []string
    for len(s) > 0 {
        if len(s) < size {
            size = len(s)
        }
        chunks = append(chunks, s[:size])
        s = s[size:]
    }
    return chunks
}

func main() {
    allowedChats = parseAllowedChats()
    
    b, _ := tb.NewBot(tb.Settings{
        Token:  os.Getenv("TELECODE_BOT_TOKEN"),
        Poller: &tb.LongPoller{Timeout: 10 * time.Second},
    })
    
    // 접근 제어 미들웨어
    b.Use(func(next tb.HandlerFunc) tb.HandlerFunc {
        return func(c tb.Context) error {
            if !isAllowed(c.Chat().ID) {
                return nil // 무시
            }
            return next(c)
        }
    })
    
    // /new - 새 세션
    b.Handle("/new", func(c tb.Context) error {
        mu.Lock()
        delete(sessions, c.Chat().ID)
        mu.Unlock()
        return c.Send("🆕 새 세션을 시작합니다.")
    })
    
    // /cli - CLI 조회/변경
    b.Handle("/cli", func(c tb.Context) error {
        args := strings.Fields(c.Message().Text)
        chatID := c.Chat().ID
        
        if len(args) == 1 {
            // 조회
            mu.RLock()
            cli := chatSettings[chatID].CLI
            mu.RUnlock()
            if cli == "" {
                cli = getDefaultCLI()
            }
            return c.Send(fmt.Sprintf("📋 현재 CLI: `%s`", cli), tb.ModeMarkdown)
        }
        
        // 변경
        newCLI := args[1]
        if newCLI != "claude" && newCLI != "opencode" {
            return c.Send("❌ 지원하지 않는 CLI입니다. (claude | opencode)")
        }
        
        // CLI 존재 확인
        if _, err := exec.LookPath(newCLI); err != nil {
            return c.Send(fmt.Sprintf("❌ `%s` CLI가 설치되어 있지 않습니다.", newCLI), tb.ModeMarkdown)
        }
        
        mu.Lock()
        settings := chatSettings[chatID]
        settings.CLI = newCLI
        chatSettings[chatID] = settings
        // 세션도 초기화 (CLI 변경 시)
        delete(sessions, chatID)
        mu.Unlock()
        
        return c.Send(fmt.Sprintf("✅ CLI 변경: `%s` (세션 초기화)", newCLI), tb.ModeMarkdown)
    })
    
    // /status - 현재 상태
    b.Handle("/status", func(c tb.Context) error {
        chatID := c.Chat().ID
        mu.RLock()
        sessionID := sessions[chatID]
        settings := chatSettings[chatID]
        mu.RUnlock()
        
        if settings.CLI == "" {
            settings.CLI = getDefaultCLI()
        }
        if sessionID == "" {
            sessionID = "없음"
        }
        if settings.Model == "" {
            settings.Model = "기본"
        }
        
        return c.Send(fmt.Sprintf(`📊 **현재 상태**
- CLI: `%s`
- 세션: `%s`
- 모델: `%s``, settings.CLI, sessionID, settings.Model), tb.ModeMarkdown)
    })
    
    // 이미지 처리
    b.Handle(tb.OnPhoto, func(c tb.Context) error {
        chatID := c.Chat().ID
        
        mu.RLock()
        cli := chatSettings[chatID].CLI
        sessionID := sessions[chatID]
        model := chatSettings[chatID].Model
        mu.RUnlock()
        
        if cli == "" {
            cli = getDefaultCLI()
        }
        
        exec := executors[cli]
        
        // 이미지 다운로드
        photo := c.Message().Photo
        fileID := photo[len(photo)-1].FileID
        file, _ := b.FileByID(fileID)
        
        tempPath := fmt.Sprintf("/tmp/telecode_img_%d.jpg", chatID)
        b.Download(&file, tempPath)
        
        prompt := c.Message().Caption
        if prompt == "" {
            prompt = "이 이미지를 분석해줘"
        }
        
        cmd := exec.BuildCommand(prompt, sessionID, tempPath, model)
        out := runCommand(cmd)
        
        // 세션 ID 저장
        mu.Lock()
        if _, ok := sessions[chatID]; !ok {
            if newSessionID := exec.ParseSessionID(out); newSessionID != "" {
                sessions[chatID] = newSessionID
            }
        }
        mu.Unlock()
        
        for _, chunk := range chunkString(out, 4000) {
            c.Send(chunk)
        }
        return nil
    })
    
    // 일반 메시지 처리
    b.Handle(tb.OnText, func(c tb.Context) error {
        chatID := c.Chat().ID
        prompt := c.Message().Text
        
        mu.RLock()
        cli := chatSettings[chatID].CLI
        sessionID := sessions[chatID]
        model := chatSettings[chatID].Model
        mu.RUnlock()
        
        if cli == "" {
            cli = getDefaultCLI()
        }
        
        exec := executors[cli]
        cmd := exec.BuildCommand(prompt, sessionID, "", model)
        out := runCommand(cmd)
        
        // 세션 ID 저장
        mu.Lock()
        if _, ok := sessions[chatID]; !ok {
            if newSessionID := exec.ParseSessionID(out); newSessionID != "" {
                sessions[chatID] = newSessionID
            }
        }
        mu.Unlock()
        
        for _, chunk := range chunkString(out, 4000) {
            c.Send(chunk)
        }
        return nil
    })
    
    b.Start()
}
```

## 배포

### 로컬 실행 (iMac)

```bash
# 빌드
go build -o telecode ./cmd/telecode

# 실행
TELECODE_BOT_TOKEN=xxx TELECODE_ALLOWED_CHATS=5788362055 ./telecode
```

### launchd (macOS)

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.telecode.bot</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/telecode</string>
    </array>
    <key>EnvironmentVariables</key>
    <dict>
        <key>TELECODE_BOT_TOKEN</key>
        <string>your-bot-token</string>
        <key>TELECODE_ALLOWED_CHATS</key>
        <string>5788362055</string>
        <key>TELECODE_DEFAULT_CLI</key>
        <string>claude</string>
    </dict>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
```

## 보안 고려사항

- 봇 토큰, Allowlist 모두 환경변수로 관리 (파일로 남지 않음)
- 미허용 chat_id에서 메시지 수신 시 무시 (응답 없음)
- CLI 존재 여부 확인 후 실행

## 상태

- 기획 완료
- 구현 미착수

## 타임라인

- 2026-02-21: 기획서 초안 작성 (opencode-telegram-bot)
- 2026-02-21: Claude Code 지원 추가, 프로젝트명 Telecode로 변경
- 2026-02-21: Executor 추상화 설계, 멀티 CLI 아키텍처 확정

