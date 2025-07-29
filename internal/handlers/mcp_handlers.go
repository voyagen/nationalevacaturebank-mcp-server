package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/voyagen/nationalevacaturebank-mcp-server/internal/api"
	"github.com/voyagen/nationalevacaturebank-mcp-server/pkg/nvb"
)

// Handlers contains all MCP tool handlers
type Handlers struct {
	apiClient *api.Client
}

// New creates a new Handlers instance
func New(client *api.Client) *Handlers {
	return &Handlers{
		apiClient: client,
	}
}

// SearchFunctionTitles handles the search_function_titles MCP tool
func (h *Handlers) SearchFunctionTitles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := request.RequireString("query")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid query parameter: %v", err)), nil
	}

	response, err := h.apiClient.SearchFunctionTitles(ctx, query)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("search failed: %v", err)), nil
	}

	resultJSON, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to serialize response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(resultJSON)), nil
}

// SearchCity handles the search_city MCP tool
func (h *Handlers) SearchCity(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	startsWith, err := request.RequireString("startsWith")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid startsWith parameter: %v", err)), nil
	}

	cities, err := h.apiClient.SearchCities(ctx, startsWith)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("city search failed: %v", err)), nil
	}

	resultJSON, err := json.Marshal(cities)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to serialize response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(resultJSON)), nil
}

// GeoLocation handles the get_geolocation MCP tool
func (h *Handlers) GeoLocation(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cityName, err := request.RequireString("cityName")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid cityName parameter: %v", err)), nil
	}

	response, err := h.apiClient.GetGeoLocation(ctx, cityName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("geolocation lookup failed: %v", err)), nil
	}

	resultJSON, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to serialize response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(resultJSON)), nil
}

// FindJobs handles the find_jobs MCP tool
func (h *Handlers) FindJobs(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters with defaults and validation
	page := request.GetInt("page", 1)
	limit := request.GetInt("limit", 10)
	sort := request.GetString("sort", "relevance")

	// Validate parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Build search parameters
	params := nvb.JobSearchParams{
		Page:      page,
		Limit:     limit,
		Sort:      sort,
		City:      request.GetString("city", ""),
		JobTitle:  request.GetString("jobTitle", ""),
		Latitude:  request.GetFloat("latitude", 0),
		Longitude: request.GetFloat("longitude", 0),
		Distance:  request.GetFloat("distance", 40.0),
	}

	response, err := h.apiClient.FindJobs(ctx, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("job search failed: %v", err)), nil
	}

	resultJSON, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to serialize response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(resultJSON)), nil
}