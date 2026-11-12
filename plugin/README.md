# EuroHop Flight Planner

Find cheap flights across Europe by searching every airport in range, not just the obvious one.

Ryanair runs Europe's largest point-to-point network, and much of it lands at secondary airports:
Frankfurt-Hahn instead of Frankfurt, Brussels-Charleroi instead of Zaventem, Milan-Bergamo instead
of Malpensa, Beauvais instead of Paris. Ordinary flight search hides those fares, because you have
to already know the airport code to type it in. This plugin starts from the place you actually want
to go and works out which airports serve it.

> **Unofficial community project.** This plugin is an independent open-source tool. It is **not
> affiliated with, authorised by, endorsed by, or connected to Ryanair DAC** or any of its
> subsidiaries. It reads publicly accessible FareFinder endpoints for informational purposes only.
> Book at [ryanair.com](https://www.ryanair.com). Fares change constantly, so a price shown here may
> differ by the time you book.

## Install

Find **EuroHop Flight Planner** in the Claude directory, or install it directly in Claude Code:

```
/plugin install eurohop-flight-planner@yugalarora/ryanair_mcp
```

## What it does

Ask in plain language — "cheapest way to get from Brussels to the Algarve in early December, back
the following weekend" — and Claude will expand both ends into nearby airports, search every viable
pairing at once, and rank the results by what the whole trip costs.

Two things it does that a normal fare search will not:

- **Multi-airport discovery.** Name a city or a region (*Andalusia*, *Algarve*, *Puglia*, *Sicily*)
  and it checks every airport within your chosen radius on both ends, then reports how far each
  airport is from where you actually want to be. A fare that is cheaper from an airport 200 km away
  is often not the better trip, so the distance is always shown alongside the price.
- **Full daily timetables.** The underlying fare API returns only the single cheapest flight per
  day. `get_daily_flights` queries four separate time windows and merges them, so you see the
  morning, afternoon and evening departures rather than one price with no time attached.

## Tools

| Tool | What it does |
| --- | --- |
| `plan_trip` | Plans a one-way or round trip between two places given as cities, regions or airport codes. Expands both ends into nearby airports, searches every pairing, ranks journeys by total price. |
| `find_nearby_airports` | Lists Ryanair-served airports within a radius of a city, region or airport code, sorted by distance. Turns a place name into airport codes. |
| `search_flights` | Cheapest direct one-way flights from one airport within a date and time window. Omit the destination to list everywhere it flies. |
| `search_return_flights` | Cheapest direct round-trip combinations between two airports across outbound and inbound date windows. |
| `get_daily_flights` | Every scheduled flight on one route on one date, with times and prices. |

All five are read-only. They issue HTTP GET requests and nothing else: no booking, no changes, no
cancellations, no account access. The plugin has no credentials and asks for none.

## Prompts

- **`plan_vacation`** — a guided workflow that confirms your constraints (date flexibility, whether
  the return must use the same airport, acceptable departure times) before running a broad search.
- **`weekend_getaway`** — scans every destination served from your home airport on the chosen dates
  and prices the cheapest complete round trips.

## Settings

Configured when you enable the plugin; every one is optional.

| Setting | Default | What it does |
| --- | --- | --- |
| Default currency | `EUR` | ISO currency code used when a search does not name one. |
| Home airport or city | *(none)* | Departure point when a search does not name one, e.g. `CRL`, `STN`, `Dublin`. |
| Default search radius | `150` km | How far from a named city or region to look for alternative airports. |

## What it sends, and where

The MCP server runs locally on your own machine as a Node process. When you call a tool, it makes
HTTPS GET requests to `services-api.ryanair.com`, carrying only the search parameters you supplied:
airport codes, dates, times, currency and an optional price ceiling. It sends no personal data, no
credentials and no conversation content, stores nothing on disk, and contacts no other host. The
airport directory is a static table compiled into the plugin, so airport and distance lookups are
answered entirely offline.

## Requirements

Node.js 18 or later, which supplies the built-in `fetch` the server uses. The server has **no npm
dependencies** — it implements the MCP stdio protocol directly, so nothing is downloaded or
installed when you add the plugin.

## Coverage

About 120 airports across Ireland, the UK, the Benelux countries, France, Spain, Portugal, Italy,
Germany, Austria, Switzerland, Poland, Czechia, Slovakia, Hungary, Scandinavia, Greece, Croatia,
Cyprus, Malta and Morocco, including regional aliases so "Algarve", "Costa del Sol" and "Salento"
resolve to the right airports. Airports outside that directory still work in the search tools if you
pass the three-letter IATA code directly.

## Licence

MIT. See [LICENSE](LICENSE).
