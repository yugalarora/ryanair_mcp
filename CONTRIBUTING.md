# Contributing to Ryanair MCP Server

Thank you for your interest in contributing to the **Ryanair MCP Server**! This is an open-source, community-driven project that enables autonomous agents and LLMs to plan regional European flights and discover nearby airport alternatives using the Model Context Protocol (MCP).

---

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md). Please be respectful and constructive in all discussions.

---

## Development Setup

### Prerequisites
- **Go 1.22+** (tested through Go 1.25)
- Git
- An MCP-compatible client for end-to-end testing (e.g. Antigravity, Claude Desktop, Cursor, or `mcp-cli`)

### Clone and Build
```bash
git clone https://github.com/yugalarora/ryanair_mcp.git
cd ryanair_mcp

# Build locally
GOPROXY=direct go build -o ryanair-mcp-server .
```

### Running Tests
```bash
GOPROXY=direct go test -v ./...
```

### Testing MCP Protocol Interactively
You can test the server directly over standard I/O by piping JSON-RPC messages:
```bash
# Test MCP initialize
printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n' | ./ryanair-mcp-server

# Test tool list
printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}\n' | ./ryanair-mcp-server
```

---

## Guidelines for Contributions

1. **Maintain Zero Dependencies Where Feasible**: The project is written in pure Go without external scripting runtimes (no Python/Node requirement). Keep new dependencies minimal and standard.
2. **Respect Ryanair Rate Limits**: All queries hit publicly available FareFinder REST endpoints. Do not remove rate limiting or add aggressive parallel scraping that could cause service denial.
3. **No Personal or Sensitive Hardcoding**: Never hardcode personal paths, developer usernames, or home airport assumptions. All defaults must be configurable via environment variables or tool arguments.
4. **Airport Directory Updates**: When adding new airport nodes to `internal/airports/airports.go`, provide exact IATA codes, correct GPS coordinates, and common search aliases.
5. **Elicitation Best Practices**: Tools and prompts should encourage agent elicitation rather than guessing user preferences (currency, airport radius, date flexibility).

---

## Submitting a Pull Request

1. Fork the repository and create a feature branch (`git checkout -b feature/amazing-feature`).
2. Commit your changes with clear, descriptive messages.
3. Run tests and verify the code compiles without warnings:
   ```bash
   go vet ./...
   go test ./...
   ```
4. Push to your branch and open a Pull Request against `main`.
5. Provide a summary of what changed and example queries tested.
