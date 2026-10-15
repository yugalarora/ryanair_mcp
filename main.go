package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"ryanair-mcp/internal/airports"
	"ryanair-mcp/internal/config"
	"ryanair-mcp/internal/planner"
	"ryanair-mcp/internal/ryanair"
)

const serverInstructions = `Ryanair Multi-Airport Flight Planner (Community / Unofficial Edition).
DISCLAIMER: This is an unofficial open-source tool and has no affiliation, endorsement, sponsorship, or connection with Ryanair DAC or any of its subsidiaries. All flight data is retrieved from publicly accessible Ryanair FareFinder APIs.

BEST PRACTICES FOR AGENTS (ELICITATION & JOURNEY EXPERIENCE):
1. Parameter Elicitation:
   - Before executing large flight scans, confirm or elicit the user's preferred currency (defaults to EUR), departure location, and travel radius.
   - Ask if dates are flexible (+/- 1-3 days): Ryanair fares frequently vary by 50-80% across adjacent days.
   - Ask about logistics: Does the user need outbound and inbound flights at the exact same airport (e.g. car parked at airport, public transit return), or are they open to open-jaw multi-city loops (e.g. fly to Malaga, rent a car / train through Andalusia, fly back from Seville or Barcelona)?
   - Preferred time windows: Ask if early morning or late night flights are acceptable, or if daytime departures are preferred.
2. Tool Usage:
   - Use 'plan_trip' for high-level requests (e.g. "fly from Brussels to Andalusia"): it automatically searches nearby airport clusters (CRL, BRU, LGG -> AGP, SVQ, XRY).
   - Use 'find_nearby_airports' to discover viable alternative hubs within driving distance (e.g. Hahn vs Cologne vs Frankfurt, or Charleroi vs Brussels).
   - Use 'get_daily_flights' when a user needs to see all departure times on a specific day, since the standard FareFinder API only returns the single cheapest flight.
`

func main() {
	cfg := config.LoadConfig()
	client := ryanair.NewClient()
	plnr := planner.NewPlanner(client)

	s := server.NewMCPServer(
		"ryanair-flight-planner",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithPromptCapabilities(true),
		server.WithInstructions(serverInstructions),
	)

	// 1. Tool: plan_trip
	registerPlanTrip(s, plnr, cfg)

	// 2. Tool: search_flights
	registerSearchFlights(s, client, cfg)

	// 3. Tool: search_return_flights
	registerSearchReturnFlights(s, client, cfg)

	// 4. Tool: get_daily_flights
	registerGetDailyFlights(s, client, cfg)

	// 5. Tool: find_nearby_airports
	registerFindNearbyAirports(s, cfg)

	// 6. Prompts for interactive workflows & elicitation
	registerPrompts(s, cfg)

	startServer(s, cfg)
}

func registerPlanTrip(s *server.MCPServer, plnr *planner.Planner, cfg config.Config) {
	tool := mcp.NewTool("plan_trip",
		mcp.WithDescription("Intelligent multi-airport trip planner for Ryanair flights. Given origin and destination cities/regions (e.g. 'Brussels' to 'Faro/Algarve'), it automatically identifies nearby departure and arrival airports within a radius, queries flight schedules across all viable airport pairs concurrently, and ranks complete journey options by price and convenience. Supports both one-way and return trips."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithString("origin",
			mcp.Required(),
			mcp.Description("Origin city name, region, or 3-letter IATA code (e.g. 'Brussels', 'London', 'CRL', 'STN', 'DUB')."),
		),
		mcp.WithString("destination",
			mcp.Required(),
			mcp.Description("Destination city, region, or 3-letter IATA code (e.g. 'Faro', 'Algarve', 'Milan', 'FAO', 'BCN')."),
		),
		mcp.WithString("outbound_date_from",
			mcp.Required(),
			mcp.Description("Outbound travel date start in YYYY-MM-DD format (e.g. '2026-11-01')."),
		),
		mcp.WithString("outbound_date_to",
			mcp.Description("Outbound travel date end in YYYY-MM-DD format (optional, defaults to outbound_date_from)."),
		),
		mcp.WithString("return_date_from",
			mcp.Description("Optional: Return travel date start for round trips in YYYY-MM-DD format (e.g. '2026-11-08')."),
		),
		mcp.WithString("return_date_to",
			mcp.Description("Optional: Return travel date end for round trips in YYYY-MM-DD format."),
		),
		mcp.WithNumber("origin_radius_km",
			mcp.Description("Search radius in km for nearby departure airports (default: 150 km)."),
		),
		mcp.WithNumber("destination_radius_km",
			mcp.Description("Search radius in km for nearby arrival airports (default: 150 km)."),
		),
		mcp.WithNumber("max_price",
			mcp.Description("Optional maximum price filter."),
		),
		mcp.WithString("currency",
			mcp.Description("Currency code (default: 'EUR')."),
		),
		mcp.WithNumber("max_options",
			mcp.Description("Maximum number of journey recommendations to return (default: 15)."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("invalid arguments"), nil
		}

		origin, _ := args["origin"].(string)
		if origin == "" {
			origin = cfg.DefaultOrigin
		}
		if origin == "" {
			return mcp.NewToolResultError("origin parameter is required"), nil
		}

		dest, _ := args["destination"].(string)
		if dest == "" {
			return mcp.NewToolResultError("destination parameter is required"), nil
		}

		outboundDateFrom, _ := args["outbound_date_from"].(string)
		if outboundDateFrom == "" {
			return mcp.NewToolResultError("outbound_date_from is required"), nil
		}

		outboundDateTo, _ := args["outbound_date_to"].(string)
		if outboundDateTo == "" {
			outboundDateTo = outboundDateFrom
		}

		returnDateFrom, _ := args["return_date_from"].(string)
		returnDateTo, _ := args["return_date_to"].(string)
		if returnDateFrom != "" && returnDateTo == "" {
			returnDateTo = returnDateFrom
		}

		originRadius := cfg.DefaultRadiusKm
		if r, ok := args["origin_radius_km"].(float64); ok && r > 0 {
			originRadius = r
		}

		destRadius := cfg.DefaultRadiusKm
		if r, ok := args["destination_radius_km"].(float64); ok && r > 0 {
			destRadius = r
		}

		var maxPrice *float64
		if mp, ok := args["max_price"].(float64); ok && mp > 0 {
			maxPrice = &mp
		}

		currency, _ := args["currency"].(string)
		if currency == "" {
			currency = cfg.DefaultCurrency
		}

		maxOptions := cfg.DefaultMaxOptions
		if mo, ok := args["max_options"].(float64); ok && mo > 0 {
			maxOptions = int(mo)
		}

		req := planner.TripRequest{
			Origin:              origin,
			Destination:         dest,
			OutboundDateFrom:    outboundDateFrom,
			OutboundDateTo:      outboundDateTo,
			ReturnDateFrom:      returnDateFrom,
			ReturnDateTo:        returnDateTo,
			OriginRadiusKm:      originRadius,
			DestinationRadiusKm: destRadius,
			MaxPrice:            maxPrice,
			Currency:            currency,
			MaxOptions:          maxOptions,
		}

		result, err := plnr.PlanTrip(ctx, req)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("trip planning failed: %v", err)), nil
		}

		resJSON, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(resJSON)), nil
	})
}

func registerSearchFlights(s *server.MCPServer, client *ryanair.Client, cfg config.Config) {
	tool := mcp.NewTool("search_flights",
		mcp.WithDescription("Search direct one-way Ryanair flights from an origin airport code within a date window. Returns the cheapest flights per route."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithString("origin",
			mcp.Required(),
			mcp.Description("3-letter IATA departure airport code (e.g. 'CRL', 'DUB', 'STN')."),
		),
		mcp.WithString("destination",
			mcp.Description("Optional: 3-letter IATA arrival airport code (e.g. 'FAO', 'BCN'). If omitted, searches all destinations."),
		),
		mcp.WithString("date_from",
			mcp.Required(),
			mcp.Description("Departure date start in YYYY-MM-DD format."),
		),
		mcp.WithString("date_to",
			mcp.Description("Optional departure date end in YYYY-MM-DD format (defaults to date_from)."),
		),
		mcp.WithString("time_from",
			mcp.Description("Optional earliest departure time HH:MM (default: '00:00')."),
		),
		mcp.WithString("time_to",
			mcp.Description("Optional latest departure time HH:MM (default: '23:59')."),
		),
		mcp.WithString("currency",
			mcp.Description("Currency code (default: 'EUR')."),
		),
		mcp.WithNumber("max_price",
			mcp.Description("Optional maximum price ceiling."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("invalid arguments"), nil
		}

		origin, _ := args["origin"].(string)
		if origin == "" {
			return mcp.NewToolResultError("origin parameter is required"), nil
		}
		dest, _ := args["destination"].(string)
		dateFrom, _ := args["date_from"].(string)
		if dateFrom == "" {
			return mcp.NewToolResultError("date_from parameter is required"), nil
		}
		dateTo, _ := args["date_to"].(string)
		if dateTo == "" {
			dateTo = dateFrom
		}
		timeFrom, _ := args["time_from"].(string)
		timeTo, _ := args["time_to"].(string)
		currency, _ := args["currency"].(string)
		if currency == "" {
			currency = cfg.DefaultCurrency
		}

		var maxPrice *float64
		if mp, ok := args["max_price"].(float64); ok && mp > 0 {
			maxPrice = &mp
		}

		flights, err := client.GetOneWayFares(ctx, origin, dest, dateFrom, dateTo, timeFrom, timeTo, currency, maxPrice)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("flight search failed: %v", err)), nil
		}

		resJSON, _ := json.MarshalIndent(map[string]interface{}{
			"count":   len(flights),
			"flights": flights,
		}, "", "  ")
		return mcp.NewToolResultText(string(resJSON)), nil
	})
}

func registerSearchReturnFlights(s *server.MCPServer, client *ryanair.Client, cfg config.Config) {
	tool := mcp.NewTool("search_return_flights",
		mcp.WithDescription("Search cheapest round-trip Ryanair flight combinations between specific airport codes across outbound and inbound date windows."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithString("origin",
			mcp.Required(),
			mcp.Description("3-letter IATA departure airport code (e.g. 'DUB', 'CRL', 'STN')."),
		),
		mcp.WithString("destination",
			mcp.Description("Optional: 3-letter IATA arrival airport code (e.g. 'FAO', 'BCN')."),
		),
		mcp.WithString("outbound_date_from",
			mcp.Required(),
			mcp.Description("Outbound date start in YYYY-MM-DD format."),
		),
		mcp.WithString("outbound_date_to",
			mcp.Required(),
			mcp.Description("Outbound date end in YYYY-MM-DD format."),
		),
		mcp.WithString("return_date_from",
			mcp.Required(),
			mcp.Description("Return date start in YYYY-MM-DD format."),
		),
		mcp.WithString("return_date_to",
			mcp.Required(),
			mcp.Description("Return date end in YYYY-MM-DD format."),
		),
		mcp.WithString("currency",
			mcp.Description("Currency code (default: 'EUR')."),
		),
		mcp.WithNumber("max_price",
			mcp.Description("Optional maximum total price filter for the round-trip."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("invalid arguments"), nil
		}

		origin, _ := args["origin"].(string)
		if origin == "" {
			return mcp.NewToolResultError("origin parameter is required"), nil
		}
		dest, _ := args["destination"].(string)
		outboundFrom, _ := args["outbound_date_from"].(string)
		outboundTo, _ := args["outbound_date_to"].(string)
		returnFrom, _ := args["return_date_from"].(string)
		returnTo, _ := args["return_date_to"].(string)
		currency, _ := args["currency"].(string)
		if currency == "" {
			currency = cfg.DefaultCurrency
		}

		var maxPrice *float64
		if mp, ok := args["max_price"].(float64); ok && mp > 0 {
			maxPrice = &mp
		}

		trips, err := client.GetRoundTripFares(ctx, origin, dest, outboundFrom, outboundTo, returnFrom, returnTo, currency, maxPrice)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("return flight search failed: %v", err)), nil
		}

		resJSON, _ := json.MarshalIndent(map[string]interface{}{
			"count": len(trips),
			"trips": trips,
		}, "", "  ")
		return mcp.NewToolResultText(string(resJSON)), nil
	})
}

func registerGetDailyFlights(s *server.MCPServer, client *ryanair.Client, cfg config.Config) {
	tool := mcp.NewTool("get_daily_flights",
		mcp.WithDescription("Get ALL scheduled Ryanair flights and their prices between two specific airport codes on a single day. Overcomes standard FareFinder limitation that only shows the single cheapest flight, letting you view full morning, afternoon, and evening flight schedules."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithString("origin",
			mcp.Required(),
			mcp.Description("3-letter IATA departure airport code (e.g. 'CRL')."),
		),
		mcp.WithString("destination",
			mcp.Required(),
			mcp.Description("3-letter IATA arrival airport code (e.g. 'FAO')."),
		),
		mcp.WithString("date",
			mcp.Required(),
			mcp.Description("Flight date in YYYY-MM-DD format (e.g. '2026-11-01')."),
		),
		mcp.WithString("currency",
			mcp.Description("Currency code (default: 'EUR')."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("invalid arguments"), nil
		}

		origin, _ := args["origin"].(string)
		dest, _ := args["destination"].(string)
		date, _ := args["date"].(string)
		currency, _ := args["currency"].(string)
		if currency == "" {
			currency = cfg.DefaultCurrency
		}

		if origin == "" || dest == "" || date == "" {
			return mcp.NewToolResultError("origin, destination, and date parameters are required"), nil
		}

		flights, err := client.GetDailyFlights(ctx, origin, dest, date, currency)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("daily flight lookup failed: %v", err)), nil
		}

		resJSON, _ := json.MarshalIndent(map[string]interface{}{
			"origin":      origin,
			"destination": dest,
			"date":        date,
			"count":       len(flights),
			"flights":     flights,
		}, "", "  ")
		return mcp.NewToolResultText(string(resJSON)), nil
	})
}

func registerFindNearbyAirports(s *server.MCPServer, cfg config.Config) {
	tool := mcp.NewTool("find_nearby_airports",
		mcp.WithDescription("Find all Ryanair-serviced European airports within a given radius (km) of any city name, region, or 3-letter IATA code, sorted by distance."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("City name, region, or 3-letter IATA airport code (e.g. 'Brussels', 'London', 'Milan', 'Frankfurt', 'CRL')."),
		),
		mcp.WithNumber("radius_km",
			mcp.Description("Maximum search radius in kilometers (default: 150 km)."),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of airports to return (default: 10)."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("invalid arguments"), nil
		}

		query, _ := args["query"].(string)
		if query == "" {
			return mcp.NewToolResultError("query parameter is required"), nil
		}

		radius := cfg.DefaultRadiusKm
		if r, ok := args["radius_km"].(float64); ok && r > 0 {
			radius = r
		}

		limit := 10
		if l, ok := args["limit"].(float64); ok && l > 0 {
			limit = int(l)
		}

		nearby, anchor, found := airports.FindNearbyAirports(query, radius, limit)
		if !found {
			return mcp.NewToolResultText(fmt.Sprintf("No airport or city matching %q found in Ryanair European directory.", query)), nil
		}

		resJSON, _ := json.MarshalIndent(map[string]interface{}{
			"query":           query,
			"anchor_airport":  anchor,
			"radius_km":       radius,
			"airports_found":  len(nearby),
			"nearby_airports": nearby,
		}, "", "  ")
		return mcp.NewToolResultText(string(resJSON)), nil
	})
}

func registerPrompts(s *server.MCPServer, cfg config.Config) {
	// 1. Prompt: plan_vacation
	vacationPrompt := mcp.NewPrompt("plan_vacation",
		mcp.WithPromptDescription("Interactive multi-airport vacation planning workflow guiding the agent to elicit travel constraints and locate optimal flights"),
		mcp.WithArgument("origin",
			mcp.ArgumentDescription("User's starting city or departure airport code (e.g. 'Brussels', 'CRL', 'London')"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("destination",
			mcp.ArgumentDescription("Target vacation destination or region (e.g. 'Andalusia', 'Algarve', 'Puglia', 'BCN')"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("departure_date",
			mcp.ArgumentDescription("Target departure date in YYYY-MM-DD format"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("return_date",
			mcp.ArgumentDescription("Optional return date in YYYY-MM-DD format"),
		),
		mcp.WithArgument("currency",
			mcp.ArgumentDescription("Preferred currency code (e.g. EUR, GBP, USD)"),
		),
	)

	s.AddPrompt(vacationPrompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := request.Params.Arguments
		origin := args["origin"]
		dest := args["destination"]
		depDate := args["departure_date"]
		retDate := args["return_date"]
		curr := args["currency"]
		if curr == "" {
			curr = cfg.DefaultCurrency
		}

		instructions := fmt.Sprintf(`Goal: Plan an optimal flight itinerary from %s to %s starting on %s.
Return date: %s
Currency: %s

Agent Steps:
1. Elicit any missing preferences from the user if not specified:
   - Are dates flexible by +/- 1-3 days? (Checking adjacent days can save 50%%+ on fares).
   - Does the return flight need to land at the exact same airport as departure (e.g. car parked at airport)?
   - Any preferred flight times (morning vs evening)?
2. Discover nearby airports using 'find_nearby_airports' for origin (%s) and destination (%s) within 150 km.
3. Call 'plan_trip' to search all viable airport pairs concurrently.
4. If the cheapest route has multiple daily departures, call 'get_daily_flights' to present the full schedule.
5. Provide a clear comparison table of options ranked by price and convenience.`,
			origin, dest, depDate, retDate, curr, origin, dest)

		msg := mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(instructions))
		return mcp.NewGetPromptResult("Vacation Planning Workflow", []mcp.PromptMessage{msg}), nil
	})

	// 2. Prompt: weekend_getaway
	weekendPrompt := mcp.NewPrompt("weekend_getaway",
		mcp.WithPromptDescription("Find the cheapest weekend flight getaways from your home airport to any European destination"),
		mcp.WithArgument("origin",
			mcp.ArgumentDescription("Departure airport IATA code or city (e.g. 'DUB', 'CRL', 'STN')"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("date_from",
			mcp.ArgumentDescription("Travel date start in YYYY-MM-DD format"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("date_to",
			mcp.ArgumentDescription("Optional travel date end in YYYY-MM-DD format"),
		),
	)

	s.AddPrompt(weekendPrompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := request.Params.Arguments
		origin := args["origin"]
		dateFrom := args["date_from"]
		dateTo := args["date_to"]
		if dateTo == "" {
			dateTo = dateFrom
		}

		instructions := fmt.Sprintf(`Goal: Find the top low-cost weekend getaways departing from %s between %s and %s in %s.
Agent Steps:
1. Call 'search_flights' from %s without specifying destination to retrieve all available direct Ryanair routes and fares.
2. Filter for the top 5 lowest-fare destinations.
3. For each top destination, search return flights using 'search_flights' for the following Sunday/Monday.
4. Present the top 3-5 complete weekend getaway combinations with departure/return times and total cost.`,
			origin, dateFrom, dateTo, cfg.DefaultCurrency, origin)

		msg := mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(instructions))
		return mcp.NewGetPromptResult("Weekend Getaway Finder", []mcp.PromptMessage{msg}), nil
	})
}

func startServer(s *server.MCPServer, cfg config.Config) {
	if cfg.Transport == "http" {
		addr := ":" + cfg.Port
		httpServer := server.NewStreamableHTTPServer(s)
		log.Printf("Starting Ryanair MCP HTTP server on %s", addr)
		if err := http.ListenAndServe(addr, httpServer); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	} else {
		log.Println("Starting Ryanair MCP server on stdio...")
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Stdio server error: %v", err)
		}
	}
}
