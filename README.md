# Nationale Vacaturebank MCP Server

A Model Context Protocol (MCP) server that provides access to the Nationale Vacaturebank (National Job Bank) API for job searching and location services in the Netherlands.

[![Go Version](https://img.shields.io/badge/Go-1.24.5+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/voyagen/nationalevacaturebank-mcp-server)

## Features

- **Job Search**: Search for job listings with advanced filtering options
- **Function Titles**: Search and autocomplete job function titles
- **City Search**: Find Dutch cities with prefix matching
- **Geolocation**: Get coordinates for Dutch cities
- **MCP Integration**: Seamless integration with Claude and other MCP clients

## Available Tools

| Tool | Description | Parameters |
|------|-------------|------------|
| `search_function_titles` | Search job function titles | `query` (required) |
| `search_city` | Search Dutch cities by prefix | `startsWith` (required) |
| `get_geolocation` | Get city coordinates | `cityName` (required) |
| `find_jobs` | Search job listings | `city`, `jobTitle`, `latitude`, `longitude`, `distance`, `page`, `limit`, `sort` |

## Installation

### Prerequisites

- Go 1.24.5 or later
- Access to the internet (for API calls)

### Clone and Build

```bash
git clone https://github.com/voyagen/nationalevacaturebank-mcp-server.git
cd nationalevacaturebank-mcp-server
go mod download
go build -o bin/server ./cmd/server
```

### Quick Start

```bash
# Run the MCP server
./bin/server

# Or run directly with Go
go run ./cmd/server
```

## Configuration

The server can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `NVB_BASE_URL` | Nationale Vacaturebank API base URL | `https://api.nationalevacaturebank.nl` |
| `NVB_TIMEOUT` | HTTP request timeout | `30s` |
| `LOG_LEVEL` | Logging level (info, debug, error) | `info` |
| `MAX_RETRIES` | Maximum API retry attempts | `3` |
| `SERVER_NAME` | MCP server name | `Nationale Vacature Bank` |
| `SERVER_VERSION` | MCP server version | `1.0.0` |

### Example Configuration

```bash
export NVB_TIMEOUT=45s
export LOG_LEVEL=debug
export MAX_RETRIES=5
./bin/server
```

## Usage Examples

### With Claude Desktop

Add to your Claude Desktop configuration (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "nationalevacaturebank": {
      "command": "/path/to/nationalevacaturebank-mcp-server/bin/server",
      "env": {
        "NVB_TIMEOUT": "30s",
        "LOG_LEVEL": "info"
      }
    }
  }
}
```

### Search for Frontend Developer Jobs

```typescript
// Example MCP tool usage
const jobs = await callTool("find_jobs", {
  jobTitle: "frontend developer",
  city: "Amsterdam",
  limit: 10,
  sort: "relevance"
});
```

### Get City Coordinates

```typescript
// Get Amsterdam coordinates
const location = await callTool("get_geolocation", {
  cityName: "Amsterdam"
});
// Returns: { cityCenter: { latitude: 52.3676, longitude: 4.9041 }, cityName: "Amsterdam" }
```

### Search Job Function Titles

```typescript
// Search for developer-related functions
const functions = await callTool("search_function_titles", {
  query: "developer"
});
// Returns: { suggestions: ["Frontend Developer", "Backend Developer", "Full Stack Developer", ...] }
```

## Architecture

The project follows Go best practices with a clean architecture:

```
├── cmd/server/           # Application entry point
├── internal/
│   ├── api/             # HTTP client and error handling
│   ├── config/          # Configuration management
│   ├── handlers/        # MCP tool handlers
│   └── server/          # MCP server setup
├── pkg/nvb/             # Public types and interfaces
└── bin/                 # Compiled binaries
```

### Key Components

- **API Client**: Handles HTTP communication with Nationale Vacaturebank API
- **Error Handling**: Comprehensive error types with context
- **Type Safety**: Proper type definitions with validation
- **Configuration**: Environment-based configuration management
- **MCP Integration**: Clean separation between MCP logic and business logic

## API Reference

### Job Search Parameters

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `city` | string | City name filter | - |
| `jobTitle` | string | Job title filter | - |
| `latitude` | number | Geographic latitude | 0 |
| `longitude` | number | Geographic longitude | 0 |
| `distance` | number | Search radius in km | 40 |
| `page` | number | Page number (1-based) | 1 |
| `limit` | number | Results per page (1-100) | 10 |
| `sort` | string | Sort order: `relevance`, `date`, `distance`, `random` | `relevance` |

### Response Types

#### Job Listing
```typescript
{
  id: string;
  title: string;
  description: string;
  company: {
    name: string;
    website: string;
    type: string;
  };
  salary: { min: number; max: number };
  contractType: string;
  careerLevel: string;
  categories: string[];
  industries: string[];
  workingHours: { min: number; max: number };
}
```

#### Geolocation Response
```typescript
{
  cityCenter: {
    latitude: number;
    longitude: number;
  };
  cityName: string;
}
```

## Development

### Running Tests

```bash
go test ./...
```

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter (requires golangci-lint)
golangci-lint run

# Check for vulnerabilities
go mod verify
```

### Building for Different Platforms

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o bin/server.exe ./cmd/server

# macOS
GOOS=darwin GOARCH=amd64 go build -o bin/server-darwin ./cmd/server

# Linux
GOOS=linux GOARCH=amd64 go build -o bin/server-linux ./cmd/server
```

## Error Handling

The server provides detailed error messages for common issues:

- **API Errors**: HTTP status codes and API response details
- **Validation Errors**: Parameter validation with specific field information
- **Network Errors**: Connection timeouts and network-related issues
- **Configuration Errors**: Invalid environment variables or missing settings

Example error response:
```json
{
  "error": "API error 400 at https://api.nationalevacaturebank.nl/api/jobs: Invalid city parameter"
}
```

## Security

- **Input Validation**: All user inputs are validated and sanitized
- **URL Encoding**: Proper encoding for API parameters
- **Error Context**: Error messages don't expose sensitive information
- **Timeout Protection**: Request timeouts prevent hanging operations

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go coding standards
- Add tests for new features
- Update documentation for API changes
- Ensure backwards compatibility for MCP tools

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Nationale Vacaturebank](https://www.nationalevacaturebank.nl) for providing the job search API
- [MCP (Model Context Protocol)](https://github.com/mark3labs/mcp-go) for the server framework
- [Anthropic](https://www.anthropic.com) for Claude integration support

## Support

- **Issues**: [GitHub Issues](https://github.com/voyagen/nationalevacaturebank-mcp-server/issues)
- **Documentation**: This README and inline code documentation
- **API Status**: Check [Nationale Vacaturebank API status](https://www.nationalevacaturebank.nl)

## Roadmap

- [ ] **Phase 2**: Add caching layer for improved performance
- [ ] **Phase 3**: Add retry logic with exponential backoff
- [ ] **Phase 4**: Add structured logging and metrics
- [ ] **Phase 5**: Add comprehensive test suite
- [ ] **Phase 6**: Add Docker containerization

---

**Note**: This server requires internet access to communicate with the Nationale Vacaturebank API. Ensure your network allows HTTPS requests to `api.nationalevacaturebank.nl`.