---
name: finding-cheap-flights
description: Use when someone wants to find or compare flights in Europe, plan a trip around price rather than a fixed airport, check what a route costs on given dates, or find where they can fly cheaply from a home airport. Covers choosing between the flight search tools, which questions to settle first, and how to present fares honestly.
---

# Finding cheap flights in Europe

Budget European flying rewards flexibility about *which airport*, not just which day. The same
week in the Algarve can differ by a factor of three depending on whether you leave from
Charleroi or Zaventem and land at Faro or Seville. The job is to search that space, then present
the trade-offs plainly, because the cheapest fare is frequently not the best trip.

## Settle these before a broad search

Each of these changes the result materially, so ask about any the person has not already told you.
Ask them together in one go, not as an interrogation, and skip any that are already obvious.

- **Date flexibility.** Are the dates fixed, or movable by a day or two? Adjacent days routinely
  differ by 50% or more. If they are flexible, widen the date range rather than running separate
  searches.
- **Same airport on the return?** A car parked at the airport forces a round trip through one
  airport. Otherwise an open-jaw route is open: fly into Malaga, travel overland, fly home from
  Seville. Open-jaw often beats a round trip outright, and it is invisible unless you ask.
- **Currency.** Fares are quoted in whatever currency is requested; get it right so the numbers
  mean something.
- **Acceptable departure times.** A 06:10 departure needs a 04:00 start, and a 23:40 arrival may
  mean no onward transport. Worth knowing before recommending a fare that depends on one.

## Choosing a tool

| Situation | Tool |
| --- | --- |
| Places named, not airports: "Brussels to Andalusia" | `plan_trip` |
| Need the airport codes for a city or region, or want to see what else is in range | `find_nearby_airports` |
| Both airport codes known, one direction | `search_flights` |
| Both airport codes known, there and back | `search_return_flights` |
| Route decided, need actual departure times | `get_daily_flights` |
| "Where can I go cheaply from here?" | `search_flights` with no destination |

`plan_trip` already expands both ends into nearby airports, so there is no need to call
`find_nearby_airports` first. Call it when the person wants to *see* the airport options, or when
`plan_trip` found nothing and you want to check whether the place name resolved sensibly.

## Reading the results

- **Only the cheapest flight per day.** `plan_trip`, `search_flights` and `search_return_flights`
  return one cheapest fare per route per day. A single result does not mean a single flight. Once a
  route is chosen, call `get_daily_flights` to get the real timetable before discussing times.
- **Distance is part of the price.** Every airport comes with `distance_km` from the place that was
  actually asked for. Report it next to the fare. Saving 30 EUR and gaining a two-hour each-way
  transfer plus a car hire is usually a worse trip, and the person can only judge that if you show
  both numbers.
- **Partial failures are reported, not hidden.** `plan_trip` returns `partial_failures` when some
  airport pairings could not be searched. If it is present, say so rather than presenting the
  ranked list as exhaustive.
- **Times are local to each airport** and carry no timezone. A departure of `2026-12-10T20:55:00`
  is 20:55 local at the departure airport. Do not convert or claim UTC.

## Presenting fares

Give a short ranked table — route, airports with their distances, dates and times, total price —
then a sentence on which option you would take and why. Three to five options is enough; a dump of
fifteen is not useful.

Two things to always state:

1. **Fares move.** These are live prices read at the moment of the search, and they change
   constantly. A price quoted now may not be there in an hour.
2. **Booking happens at ryanair.com.** This tooling is read-only and cannot book, hold, change or
   cancel anything. Do not imply a seat is reserved or a price is locked.

If a fare looks implausible — a few euros for a long route, or a price far below everything around
it — flag it as possibly a data artefact rather than presenting it as a sure thing.

## Unofficial tool

This searches publicly accessible Ryanair FareFinder data and is not affiliated with, authorised by
or endorsed by Ryanair. Say so if someone asks where the data comes from, or seems to think they are
talking to the airline. Do not speak as Ryanair, and do not make commitments about baggage rules,
check-in, changes or refunds; those are the airline's to state, and this tool has no access to them.
