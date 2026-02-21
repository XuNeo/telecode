package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"telecode/internal/bot"
	"telecode/internal/config"
)

func main() {
	// Command line flags
	configPath := flag.String("config", "", "Path to config file (default: auto-detect)")
	generateConfig := flag.Bool("generate-config", false, "Generate example config file")
	flag.Parse()

	// Generate example config if requested
	if *generateConfig {
		path := "telecode.yml"
		if err := config.CreateExampleConfig(path); err != nil {
			fmt.Printf("❌ Failed to create example config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Example config created: %s\n", path)
		fmt.Println("📝 Please edit the file and add your bot tokens")
		os.Exit(0)
	}

	// Determine config path
	if *configPath == "" {
		*configPath = config.GetDefaultConfigPath()
		if *configPath == "" {
			fmt.Println("❌ No config file found")
			fmt.Println("💡 Use -generate-config to create an example config")
			fmt.Println("💡 Or specify config path with -config flag")
			os.Exit(1)
		}
	}

	// Load configuration
	fmt.Printf("📄 Loading config from: %s\n", *configPath)
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("❌ Failed to load config: %v\n", err)
		os.Exit(1)
	}

	if len(cfg.Workspaces) == 0 {
		fmt.Println("❌ No workspaces defined in config")
		os.Exit(1)
	}

	fmt.Printf("📦 Found %d workspace(s)\n", len(cfg.Workspaces))

	// Create multi-bot manager
	manager, err := bot.NewManager(cfg)
	if err != nil {
		fmt.Printf("❌ Failed to create bot manager: %v\n", err)
		os.Exit(1)
	}

	// Setup context for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start all bots
	if err := manager.Start(ctx); err != nil {
		fmt.Printf("❌ Failed to start bots: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ All bots are running!")
	fmt.Println("👋 Press Ctrl+C to stop")

	// Wait for shutdown signal
	<-ctx.Done()
	fmt.Println("\n👋 Shutting down...")
}
