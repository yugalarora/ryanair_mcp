'use strict';

// Multi-airport trip planner: expands an origin and destination into clusters
// of nearby airports, searches every viable pairing, and ranks the resulting
// journeys by price. Ported from internal/planner/planner.go.

const { findNearbyAirports } = require('./airports');
const { getOneWayFares, getRoundTripFares } = require('./ryanair');

const MAX_AIRPORTS_PER_END = 5;
const CONCURRENCY = 4; // Keep the request rate gentle on the public API.

const MONTHS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

/** Render "2026-12-10T20:55:00" as "10 Dec 20:55" for human-readable summaries. */
function formatDateTime(value) {
  const match = /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2})/.exec(String(value || ''));
  if (!match) return String(value || 'unknown time');
  const [, , month, day, hour, minute] = match;
  const monthName = MONTHS[Number(month) - 1] || month;
  return `${day} ${monthName} ${hour}:${minute}`;
}

/** Run `tasks` with at most `limit` in flight at once, preserving order. */
async function runPooled(tasks, limit) {
  const results = new Array(tasks.length);
  let next = 0;

  async function worker() {
    while (next < tasks.length) {
      const index = next++;
      results[index] = await tasks[index]();
    }
  }

  const workers = [];
  for (let i = 0; i < Math.min(limit, tasks.length); i += 1) {
    workers.push(worker());
  }
  await Promise.all(workers);
  return results;
}

/** Shrink a nearby-airport record to the fields we report back. */
function toAirportRef(airport) {
  return {
    iata: airport.iata,
    name: airport.name,
    city: airport.city,
    country: airport.country,
    distance_km: airport.distance_km,
  };
}

/**
 * Resolve one end of the trip into a list of candidate airports.
 * Falls back to treating the query as a literal IATA code when it is not in
 * the directory, so routes we have no coordinates for still work.
 */
function resolveEnd(query, radiusKm) {
  const resolved = findNearbyAirports(query, radiusKm, MAX_AIRPORTS_PER_END);

  if (!resolved) {
    const code = String(query).trim().toUpperCase();
    return {
      anchor: null,
      anchorLabel: String(query),
      airports: [{ iata: code, name: code, city: code, country: '', distance_km: 0 }],
    };
  }

  return {
    anchor: resolved.anchor,
    anchorLabel: `${resolved.anchor.city} (${resolved.anchor.iata}, ${resolved.anchor.country})`,
    airports: resolved.nearby.map(toAirportRef),
  };
}

/** Append "[85km from Brussels]" notes when an airport is not the anchor itself. */
function distanceNotes(originRef, destRef, originAnchor, destAnchor) {
  let notes = '';
  if (originRef.distance_km > 0 && originAnchor) {
    notes += ` [${Math.round(originRef.distance_km)}km from ${originAnchor.city}]`;
  }
  if (destRef.distance_km > 0 && destAnchor) {
    notes += ` [${Math.round(destRef.distance_km)}km from ${destAnchor.city}]`;
  }
  return notes;
}

/**
 * Plan a one-way or round trip across every viable airport pairing.
 * Returns the ranked journeys plus the airport clusters that were searched.
 */
async function planTrip(request) {
  const currency = request.currency || 'EUR';
  const originRadiusKm = request.originRadiusKm > 0 ? request.originRadiusKm : 150;
  const destRadiusKm = request.destinationRadiusKm > 0 ? request.destinationRadiusKm : 150;
  const maxOptions = request.maxOptions > 0 ? request.maxOptions : 15;

  const outboundDateFrom = request.outboundDateFrom;
  const outboundDateTo = request.outboundDateTo || outboundDateFrom;

  const isRoundTrip = Boolean(request.returnDateFrom);
  const returnDateFrom = request.returnDateFrom || '';
  const returnDateTo = request.returnDateTo || returnDateFrom;

  const originEnd = resolveEnd(request.origin, originRadiusKm);
  const destEnd = resolveEnd(request.destination, destRadiusKm);

  const pairs = [];
  for (const originRef of originEnd.airports) {
    for (const destRef of destEnd.airports) {
      if (originRef.iata !== destRef.iata) pairs.push({ originRef, destRef });
    }
  }

  if (pairs.length === 0) {
    throw new Error(
      `Origin "${request.origin}" and destination "${request.destination}" resolved to the same ` +
        'airport, so there is no route to search. Widen a radius or name a different endpoint.'
    );
  }

  const tasks = pairs.map(({ originRef, destRef }) => async () => {
    try {
      if (isRoundTrip) {
        const trips = await getRoundTripFares({
          origin: originRef.iata,
          destination: destRef.iata,
          outDateFrom: outboundDateFrom,
          outDateTo: outboundDateTo,
          retDateFrom: returnDateFrom,
          retDateTo: returnDateTo,
          currency,
          maxPrice: request.maxPrice,
        });

        return {
          journeys: trips.map((trip) => ({
            type: 'round_trip',
            total_price: trip.total_price,
            currency: trip.currency || currency,
            origin_airport: originRef,
            destination_airport: destRef,
            outbound_flight: trip.outbound,
            return_flight: trip.inbound,
            summary:
              `Round-trip: ${originRef.city} (${originRef.iata}) to ${destRef.city} ` +
              `(${destRef.iata}) on ${formatDateTime(trip.outbound.departure_time)}, ` +
              `returning ${formatDateTime(trip.inbound.departure_time)} for ` +
              `${trip.total_price.toFixed(2)} ${trip.currency || currency}` +
              distanceNotes(originRef, destRef, originEnd.anchor, destEnd.anchor),
          })),
          error: null,
        };
      }

      const flights = await getOneWayFares({
        origin: originRef.iata,
        destination: destRef.iata,
        dateFrom: outboundDateFrom,
        dateTo: outboundDateTo,
        timeFrom: '00:00',
        timeTo: '23:59',
        currency,
        maxPrice: request.maxPrice,
      });

      return {
        journeys: flights.map((flight) => ({
          type: 'one_way',
          total_price: flight.price,
          currency: flight.currency || currency,
          origin_airport: originRef,
          destination_airport: destRef,
          outbound_flight: flight,
          summary:
            `One-way: ${originRef.city} (${originRef.iata}) to ${destRef.city} ` +
            `(${destRef.iata}) on ${formatDateTime(flight.departure_time)} for ` +
            `${flight.price.toFixed(2)} ${flight.currency || currency}` +
            distanceNotes(originRef, destRef, originEnd.anchor, destEnd.anchor),
        })),
        error: null,
      };
    } catch (err) {
      // One dead airport pair must not fail the whole plan; record and continue.
      return {
        journeys: [],
        error: `${originRef.iata}-${destRef.iata}: ${err.message}`,
      };
    }
  });

  const outcomes = await runPooled(tasks, CONCURRENCY);

  const journeys = [];
  const errors = [];
  for (const outcome of outcomes) {
    journeys.push(...outcome.journeys);
    if (outcome.error) errors.push(outcome.error);
  }

  if (journeys.length === 0 && errors.length === pairs.length) {
    throw new Error(
      `All ${pairs.length} airport pairings failed. First error: ${errors[0]}`
    );
  }

  journeys.sort((a, b) => a.total_price - b.total_price);

  return {
    trip_type: isRoundTrip ? 'round_trip' : 'one_way',
    origin_anchor: originEnd.anchorLabel,
    destination_anchor: destEnd.anchorLabel,
    currency,
    outbound_date_range: { from: outboundDateFrom, to: outboundDateTo },
    return_date_range: isRoundTrip ? { from: returnDateFrom, to: returnDateTo } : null,
    origin_airports_checked: originEnd.airports,
    dest_airports_checked: destEnd.airports,
    airport_pairs_searched: pairs.length,
    journeys_found: journeys.length,
    recommended_journeys: journeys.slice(0, maxOptions),
    // Present only when some pairings failed, so the agent can say so.
    partial_failures: errors.length > 0 ? errors : undefined,
  };
}

module.exports = { planTrip, formatDateTime };
