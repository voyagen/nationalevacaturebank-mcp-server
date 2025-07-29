package main

import (
	"fmt"

	"github.com/voyagen/nationalevacaturebank-mcp-server/internal/api"
	"github.com/voyagen/nationalevacaturebank-mcp-server/internal/config"
	"github.com/voyagen/nationalevacaturebank-mcp-server/internal/handlers"
	"github.com/voyagen/nationalevacaturebank-mcp-server/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Configuration error: %v\n", err)
		return
	}

	// Create API client
	apiClient := api.NewClient(cfg.BaseURL, cfg.Timeout)

	// Create handlers
	mcpHandlers := handlers.New(apiClient)

	// Create and configure MCP server
	mcpServer := server.New(cfg.ServerName, cfg.ServerVersion, mcpHandlers)

	// Start the server
	if err := mcpServer.Serve(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}