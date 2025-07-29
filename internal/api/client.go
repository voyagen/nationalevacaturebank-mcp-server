package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/voyagen/nationalevacaturebank-mcp-server/pkg/nvb"
)

// Client implements the NationalevacaturebankAPI interface
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new API client
func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL: baseURL,
	}
}

// makeRequest performs an HTTP GET request and handles common error cases
func (c *Client) makeRequest(ctx context.Context, url string, target interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return NewAPIError(0, "failed to create request", url, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return NewAPIError(0, "API request failed", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return NewAPIError(resp.StatusCode, string(body), url, nil)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewAPIError(resp.StatusCode, "failed to read response", url, err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return NewAPIError(resp.StatusCode, "failed to parse JSON", url, err)
	}

	return nil
}

// SearchFunctionTitles searches for job function titles
func (c *Client) SearchFunctionTitles(ctx context.Context, query string) (*nvb.FunctionTitlesResponse, error) {
	if query == "" {
		return nil, NewValidationError("query", query, "query cannot be empty")
	}

	apiURL := fmt.Sprintf("%s/api/jobs/v3/sites/nationalevacaturebank.nl/function-titles?query=%s",
		c.baseURL, url.QueryEscape(query))

	var response nvb.FunctionTitlesResponse
	if err := c.makeRequest(ctx, apiURL, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SearchCities searches for cities by prefix
func (c *Client) SearchCities(ctx context.Context, startsWith string) ([]string, error) {
	if startsWith == "" {
		return nil, NewValidationError("startsWith", startsWith, "startsWith cannot be empty")
	}

	apiURL := fmt.Sprintf("%s/api/v1/cities/nl?startsWith=%s",
		c.baseURL, url.QueryEscape(startsWith))

	var cities []string
	if err := c.makeRequest(ctx, apiURL, &cities); err != nil {
		return nil, err
	}

	return cities, nil
}

// GetGeoLocation gets geolocation data for a city
func (c *Client) GetGeoLocation(ctx context.Context, cityName string) (*nvb.GeoLocationResponse, error) {
	if cityName == "" {
		return nil, NewValidationError("cityName", cityName, "cityName cannot be empty")
	}

	apiURL := fmt.Sprintf("%s/api/v1/geolocations/nl/%s",
		c.baseURL, url.PathEscape(cityName))

	// First get the raw response to handle string coordinates
	var rawResponse struct {
		CityCenter struct {
			Latitude  string `json:"latitude"`
			Longitude string `json:"longitude"`
		} `json:"cityCenter"`
		CityName string `json:"cityName"`
	}

	if err := c.makeRequest(ctx, apiURL, &rawResponse); err != nil {
		return nil, err
	}

	// Convert string coordinates to float64
	lat, err := strconv.ParseFloat(rawResponse.CityCenter.Latitude, 64)
	if err != nil {
		return nil, NewAPIError(0, "invalid latitude format", apiURL, err)
	}

	lng, err := strconv.ParseFloat(rawResponse.CityCenter.Longitude, 64)
	if err != nil {
		return nil, NewAPIError(0, "invalid longitude format", apiURL, err)
	}

	response := &nvb.GeoLocationResponse{
		CityCenter: nvb.Coordinates{
			Latitude:  lat,
			Longitude: lng,
		},
		CityName: rawResponse.CityName,
	}

	return response, nil
}

// FindJobs searches for job listings with filters
func (c *Client) FindJobs(ctx context.Context, params nvb.JobSearchParams) (*nvb.JobsResponse, error) {
	// Validate parameters
	if params.Page < 1 {
		return nil, NewValidationError("page", params.Page, "page must be >= 1")
	}
	if params.Limit < 1 {
		return nil, NewValidationError("limit", params.Limit, "limit must be >= 1")
	}
	if params.Limit > 100 {
		return nil, NewValidationError("limit", params.Limit, "limit must be <= 100")
	}

	baseURL := fmt.Sprintf("%s/api/jobs/v3/sites/nationalevacaturebank.nl/jobs", c.baseURL)
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, NewAPIError(0, "failed to parse URL", baseURL, err)
	}

	// Build query parameters
	q := u.Query()
	q.Set("page", fmt.Sprintf("%d", params.Page))
	q.Set("limit", fmt.Sprintf("%d", params.Limit))
	q.Set("sort", params.Sort)

	// Build filters
	var filters []string

	// Add geolocation filters if both latitude and longitude are provided
	if params.Latitude != 0 && params.Longitude != 0 {
		filters = append(filters, fmt.Sprintf("latitude:%.6f", params.Latitude))
		filters = append(filters, fmt.Sprintf("longitude:%.6f", params.Longitude))
		filters = append(filters, fmt.Sprintf("distance:%.0f", params.Distance))
	}

	// Add city filter if provided
	if params.City != "" {
		filters = append(filters, fmt.Sprintf("city:%s", params.City))
	}

	// Add job title filter if provided
	if params.JobTitle != "" {
		filters = append(filters, fmt.Sprintf("dcoTitle:%s", params.JobTitle))
	}

	if len(filters) > 0 {
		q.Set("filters", strings.Join(filters, " "))
	}

	u.RawQuery = q.Encode()

	var response nvb.JobsResponse
	if err := c.makeRequest(ctx, u.String(), &response); err != nil {
		return nil, err
	}

	return &response, nil
}