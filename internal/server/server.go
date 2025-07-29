package server

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/voyagen/nationalevacaturebank-mcp-server/internal/handlers"
)

// MCPServer wraps the MCP server with our handlers
type MCPServer struct {
	server   *server.MCPServer
	handlers *handlers.Handlers
}

// New creates a new MCP server instance
func New(name, version string, handlers *handlers.Handlers) *MCPServer {
	s := server.NewMCPServer(
		name,
		version,
		server.WithToolCapabilities(true),
	)

	mcpServer := &MCPServer{
		server:   s,
		handlers: handlers,
	}

	mcpServer.registerTools()
	return mcpServer
}

// registerTools registers all MCP tools with their handlers
func (s *MCPServer) registerTools() {
	// Search function titles tool
	searchTitlesTool := mcp.NewTool("search_function_titles",
		mcp.WithDescription("Search job function titles from National Evacature Bank API"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query for function titles"),
		),
	)
	s.server.AddTool(searchTitlesTool, s.handlers.SearchFunctionTitles)

	// Search city tool
	searchCityTool := mcp.NewTool("search_city",
		mcp.WithDescription("Search cities in Netherlands"),
		mcp.WithString("startsWith",
			mcp.Required(),
			mcp.Description("City name prefix to search for"),
		),
	)
	s.server.AddTool(searchCityTool, s.handlers.SearchCity)

	// Geolocation tool
	geoLocationTool := mcp.NewTool("get_geolocation",
		mcp.WithDescription("Get geolocation data for a city in Netherlands"),
		mcp.WithString("cityName",
			mcp.Required(),
			mcp.Description("City name to get geolocation for"),
		),
	)
	s.server.AddTool(geoLocationTool, s.handlers.GeoLocation)

	// Find jobs tool
	findJobsTool := mcp.NewTool("find_jobs",
		mcp.WithDescription("Search for job listings with filters"),
		mcp.WithString("city", mcp.Description("City name")),
		mcp.WithString("jobTitle", mcp.Description("Job title or function")),
		mcp.WithNumber("latitude", mcp.Description("Latitude coordinate")),
		mcp.WithNumber("longitude", mcp.Description("Longitude coordinate")),
		mcp.WithNumber("distance", mcp.Description("Search radius in kilometers (default: 40)")),
		mcp.WithNumber("page", mcp.Description("Page number (default: 1)")),
		mcp.WithNumber("limit", mcp.Description("Results per page (default: 10)")),
		mcp.WithString("sort", mcp.Description("Sort by: relevance, date, distance, random (default: relevance)")),
	)
	s.server.AddTool(findJobsTool, s.handlers.FindJobs)
}

// Serve starts the MCP server
func (s *MCPServer) Serve() error {
	return server.ServeStdio(s.server)
}