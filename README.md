# Ryanair Flight Planner MCP Server

[![CI](https://github.com/yugalarora/ryanair_mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/yugalarora/ryanair_mcp/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/yugalarora/ryanair_mcp?display_name=tag)](https://github.com/yugalarora/ryanair_mcp/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/yugalarora/ryanair_mcp)](https://goreportcard.com/report/github.com/yugalarora/ryanair_mcp)

An ultra-fast, lightweight **Model Context Protocol (MCP)** server written in **Go** that equips autonomous AI assistants (Antigravity, Claude Desktop, Cursor, etc.) with real-time flight search, multi-airport regional discovery, and intelligent travel planning capabilities.

> [!WARNING]
> **Unofficial Community Project & Disclaimer**  
> This project is an independent community development and is **not affiliated with, authorized by, endorsed by, sponsored by, or in any way connected to Ryanair DAC or any of its subsidiaries or affiliates**. The official Ryanair website can be found at [ryanair.com](https://www.ryanair.com). All flight information is retrieved from publicly accessible FareFinder API endpoints for informational purposes.

---

## 🌟 Why Ryanair MCP?

Ryanair operates Europe's largest point-to-point flight network across hundreds of regional secondary airports (e.g. Frankfurt-Hahn, Brussels-Charleroi, Milan-Bergamo, Beauvais-Paris). Traditional flight search engines frequently fail to show budget options because travelers do not know which secondary airports serve their intended destination region.

This server solves that by providing:
1. **Multi-Airport Clustered Discovery**: Enter a region (e.g. *"Andalusia"*, *"Algarve"*, *"Puglia"*) or city (e.g. *"Brussels"*), and the engine concurrently checks all departure and arrival airports within an adjustable kilometer radius.
2. **Full Daily Flight Timetables**: Standard FareFinder APIs only return the single cheapest flight per day. Our `get_daily_flights` engine queries multiple departure windows to surface **all available flights** (morning, afternoon, evening) on that date.
3. **Agent Elicitation & Interactive Prompts**: Built-in MCP Prompts (`plan_vacation`, `weekend_getaway`) and system instructions that guide LLMs to ask users about date flexibility, currency, and airport constraints before executing expensive queries.
4. **Pure Go / Zero Runtime Dependencies**: Compiles to a single static binary. No Python, Node.js, or browser headless drivers required.

---

## 🛠️ Tools Exposed

| Tool Name | Purpose | Description & Key Parameters |
| :--- | :--- | :--- |
| **`plan_trip`** | **Multi-Airport Trip Planner** | Given an origin and destination (city, region, or IATA code), discovers nearby airports within a radius, concurrently scans all route pairings, and ranks results by price and convenience.<br>• `origin` (required): e.g. `'Brussels'`, `'London'`, `'DUB'`<br>• `destination` (required): e.g. `'Andalusia'`, `'Algarve'`, `'BCN'`<br>• `outbound_date_from`: `YYYY-MM-DD`<br>• `return_date_from`: Optional for round trips<br>• `origin_radius_km` / `destination_radius_km`: (default: 150 km)<br>• `currency`: (default: `EUR`) |
| **`get_daily_flights`** | **Full Schedule Discovery** | Overcomes standard 1-flight-per-day limits by scanning morning, afternoon, and evening windows to return all scheduled flights and prices for a specific route and day.<br>• `origin`: 3-letter IATA code<br>• `destination`: 3-letter IATA code<br>• `date`: `YYYY-MM-DD`<br>• `currency`: `EUR` |
| **`find_nearby_airports`** | **Geographic Airport Search** | Computes Haversine great-circle distances against a built-in database of 80+ Ryanair hubs to return all airports within a given radius of a city or IATA code.<br>• `query`: City name or IATA code<br>• `radius_km`: Search radius (default: 150 km)<br>• `limit`: Max airports to return |
| **`search_flights`** | **One-Way Direct Search** | Searches cheapest direct flights between two specific airports with departure time and price ceiling filters.<br>• `origin`, `destination`, `date_from`, `date_to`, `time_from`, `time_to`, `max_price`, `currency` |
| **`search_return_flights`** | **Round-Trip Search** | Searches round-trip combinations between two specific airport codes across outbound and inbound date windows.<br>• `origin`, `destination`, `outbound_date_from`, `outbound_date_to`, `return_date_from`, `return_date_to`, `currency` |

---

## 💬 Prompts & Elicitation

The server registers standard MCP Prompts that autonomous agents can invoke or inspect:

* **`plan_vacation`**: Interactive workflow that asks the agent to elicit user preferences (date flexibility, same-airport return vs open-jaw loop, preferred flight times) before executing airport radius scans.
* **`weekend_getaway`**: Automatically checks all non-stop Ryanair routes from a user's local airport to find the lowest-cost weekend getaways.

---

## ⚙️ Configuration & Environment Variables

All parameters can be customized without touching code via environment variables:

| Environment Variable | Default | Description |
| :--- | :--- | :--- |
| `RYANAIR_CURRENCY` | `EUR` | Default currency code (`EUR`, `GBP`, `USD`, `PLN`, etc.). |
| `RYANAIR_RADIUS_KM` | `150` | Default search radius (in kilometers) for nearby airports. |
| `RYANAIR_DEFAULT_ORIGIN` | `""` | Optional fallback origin city or airport code (e.g. `CRL`, `STN`). |
| `RYANAIR_MAX_OPTIONS` | `15` | Default maximum recommendations returned by `plan_trip`. |
| `MCP_TRANSPORT` | `stdio` | Transport mechanism: `stdio` (default) or `http`. |
| `MCP_PORT` | `8080` | Port used when `MCP_TRANSPORT=http`. |

---

## 📦 Installation & Setup

### Option 1: Download Pre-built Binary
Download the pre-compiled binary for your operating system and CPU architecture from the [GitHub Releases](https://github.com/yugalarora/ryanair_mcp/releases) page:
* **macOS**: `ryanair-mcp-v1.0.0-darwin-arm64.tar.gz` (Apple Silicon M1/M2/M3/M4) or `darwin-amd64` (Intel)
* **Linux**: `ryanair-mcp-v1.0.0-linux-amd64.tar.gz` or `linux-arm64.tar.gz`
* **Windows**: `ryanair-mcp-v1.0.0-windows-amd64.zip`

Extract the binary and place it in your preferred path (e.g. `/usr/local/bin` or `C:\tools\`).

### Option 2: Build from Source
Requires **Go 1.22+**:

```bash
git clone https://github.com/yugalarora/ryanair_mcp.git
cd ryanair_mcp

# Build executable
GOPROXY=direct go build -o ryanair-mcp-server .
```

### Option 3: Run with Docker
```bash
# Build Docker image
docker build -t ryanair-mcp-server .

# Run in stdio mode
docker run -i --rm ryanair-mcp-server
```

---

## 🔌 Connecting to AI Clients

### 1. Antigravity CLI / IDE
Add the server to `~/.gemini/config/mcp_config.json`:

```json
{
  "mcpServers": {
    "ryanair-planner": {
      "command": "/path/to/ryanair-mcp-server",
      "args": [],
      "env": {
        "RYANAIR_CURRENCY": "EUR",
        "RYANAIR_RADIUS_KM": "150"
      }
    }
  }
}
```

### 2. Claude Desktop
Add to `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or `%APPDATA%\Claude\claude_desktop_config.json` (Windows):

```json
{
  "mcpServers": {
    "ryanair-planner": {
      "command": "/path/to/ryanair-mcp-server"
    }
  }
}
```

### 3. Cursor
Add to `.cursor/mcp.json` in your project or global Cursor settings:

```json
{
  "mcpServers": {
    "ryanair-planner": {
      "command": "/path/to/ryanair-mcp-server"
    }
  }
}
```

---

## 🌐 Publishing to MCP Registries

If you are publishing or distributing this MCP server, the ecosystem supports several registries:

1. **[Smithery.ai](https://smithery.ai/)**:
   - The primary registry and package manager for MCP servers.
   - This repository already includes `smithery.yaml`.
   - Once published to GitHub, users can install it via:
     ```bash
     npx -y @smithery/cli install @yugalarora/ryanair_mcp --client claude
     ```
2. **[Glama.ai](https://glama.ai/mcp/servers)**:
   - Community registry for discovering MCP servers. Submit your GitHub repository URL on the Glama directory to index tools and documentation automatically.
3. **[PulseMCP](https://www.pulsemcp.com/)**:
   - Curated directory of Model Context Protocol implementations.
4. **Official Model Context Protocol Servers Repository**:
   - You can submit a pull request to [modelcontextprotocol/servers](https://github.com/modelcontextprotocol/servers) under the community implementations list.

---

## 🧪 Testing

Run automated tests:

```bash
GOPROXY=direct go test -v ./...
```

Verify stdio JSON-RPC initialization:

```bash
printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}\n{"jsonrpc":"2.0","id":2,"method":"prompts/list","params":{}}\n' | ./ryanair-mcp-server
```

---

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for local development workflows and guidelines.

---

## 🔒 Security

For security vulnerability disclosures, please refer to [SECURITY.md](SECURITY.md).

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).
