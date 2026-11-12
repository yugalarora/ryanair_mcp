'use strict';

// Read-only client for Ryanair's publicly accessible FareFinder (farfnd v4)
// endpoints. Every request is an HTTP GET; nothing is ever booked, modified
// or written. Ported from internal/ryanair/client.go.

const BASE_URL = 'https://services-api.ryanair.com/farfnd/v4/';
const REQUEST_TIMEOUT_MS = 15000;
const USER_AGENT = 'eurohop-flight-planner/1.0 (+https://github.com/yugalarora/ryanair_mcp)';

/** Normalise "FR2941" to "FR 2941" for readability. */
function formatFlightNumber(raw) {
  if (typeof raw !== 'string') return '';
  if (raw.length > 2 && raw.slice(0, 2) === 'FR' && raw[2] !== ' ') {
    return 'FR ' + raw.slice(2);
  }
  return raw;
}

/** Flatten one FareFinder segment into our flat flight shape. */
function parseSegment(raw) {
  const dep = (raw && raw.departureAirport) || {};
  const arr = (raw && raw.arrivalAirport) || {};
  const price = (raw && raw.price) || {};

  return {
    origin: dep.iataCode || '',
    origin_name: dep.name || '',
    origin_city: (dep.city && dep.city.name) || '',
    origin_country: dep.countryName || '',
    destination: arr.iataCode || '',
    destination_name: arr.name || '',
    destination_city: (arr.city && arr.city.name) || '',
    destination_country: arr.countryName || '',
    // Local airport time as reported by the API, no timezone conversion.
    departure_time: raw.departureDate || '',
    arrival_time: raw.arrivalDate || '',
    flight_number: formatFlightNumber(raw.flightNumber),
    price: typeof price.value === 'number' ? price.value : 0,
    currency: price.currencyCode || '',
  };
}

/** GET a FareFinder endpoint and return the parsed JSON body. */
async function getJson(endpoint, params) {
  const url = new URL(BASE_URL + endpoint);
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== null && value !== '') {
      url.searchParams.set(key, String(value));
    }
  }

  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);

  try {
    const response = await fetch(url, {
      method: 'GET',
      headers: { Accept: 'application/json', 'User-Agent': USER_AGENT },
      signal: controller.signal,
    });

    if (!response.ok) {
      throw new Error(
        `Ryanair FareFinder returned HTTP ${response.status} for ${endpoint}. ` +
          'The route may not exist, the date may be outside the published schedule, ' +
          'or the service may be rate limiting. Try a different date or airport pair.'
      );
    }

    return await response.json();
  } catch (err) {
    if (err && err.name === 'AbortError') {
      throw new Error(
        `Ryanair FareFinder did not respond within ${REQUEST_TIMEOUT_MS / 1000}s for ${endpoint}.`
      );
    }
    throw err;
  } finally {
    clearTimeout(timer);
  }
}

/**
 * Cheapest one-way fares departing `origin` in the given date and time window.
 * `destination` is optional; omitting it returns every destination served.
 */
async function getOneWayFares({
  origin,
  destination,
  dateFrom,
  dateTo,
  timeFrom,
  timeTo,
  currency,
  maxPrice,
}) {
  const data = await getJson('oneWayFares', {
    departureAirportIataCode: origin,
    arrivalAirportIataCode: destination,
    outboundDepartureDateFrom: dateFrom,
    outboundDepartureDateTo: dateTo,
    outboundDepartureTimeFrom: timeFrom || '00:00',
    outboundDepartureTimeTo: timeTo || '23:59',
    currency: currency || 'EUR',
    priceValueTo: maxPrice > 0 ? Math.round(maxPrice) : undefined,
  });

  const fares = (data && data.fares) || [];
  return fares.filter((f) => f && f.outbound).map((f) => parseSegment(f.outbound));
}

/** Cheapest round-trip combinations across the outbound and inbound windows. */
async function getRoundTripFares({
  origin,
  destination,
  outDateFrom,
  outDateTo,
  retDateFrom,
  retDateTo,
  currency,
  maxPrice,
}) {
  const data = await getJson('roundTripFares', {
    departureAirportIataCode: origin,
    arrivalAirportIataCode: destination,
    outboundDepartureDateFrom: outDateFrom,
    outboundDepartureDateTo: outDateTo,
    inboundDepartureDateFrom: retDateFrom,
    inboundDepartureDateTo: retDateTo,
    outboundDepartureTimeFrom: '00:00',
    outboundDepartureTimeTo: '23:59',
    inboundDepartureTimeFrom: '00:00',
    inboundDepartureTimeTo: '23:59',
    currency: currency || 'EUR',
    priceValueTo: maxPrice > 0 ? Math.round(maxPrice) : undefined,
  });

  const fares = (data && data.fares) || [];
  const trips = [];

  for (const fare of fares) {
    if (!fare || !fare.outbound || !fare.inbound) continue;

    const outbound = parseSegment(fare.outbound);
    const inbound = parseSegment(fare.inbound);

    const summaryPrice =
      fare.summary && fare.summary.price && typeof fare.summary.price.value === 'number'
        ? fare.summary.price.value
        : 0;

    trips.push({
      outbound,
      inbound,
      total_price: summaryPrice > 0 ? summaryPrice : outbound.price + inbound.price,
      currency: outbound.currency,
    });
  }

  return trips;
}

// FareFinder returns only the single cheapest flight per time window, so to
// see a full day's schedule we query several narrower windows and merge.
const DAY_WINDOWS = [
  ['00:00', '11:00'],
  ['11:01', '15:30'],
  ['15:31', '19:30'],
  ['19:31', '23:59'],
];

/**
 * Every distinct scheduled flight on one route for one date, ordered by
 * departure time. Windows that error are skipped so a partial outage still
 * returns the flights we did find.
 */
async function getDailyFlights({ origin, destination, date, currency }) {
  const results = await Promise.allSettled(
    DAY_WINDOWS.map(([timeFrom, timeTo]) =>
      getOneWayFares({
        origin,
        destination,
        dateFrom: date,
        dateTo: date,
        timeFrom,
        timeTo,
        currency,
      })
    )
  );

  const seen = new Map();
  let failures = 0;

  for (const result of results) {
    if (result.status !== 'fulfilled') {
      failures += 1;
      continue;
    }
    for (const flight of result.value) {
      seen.set(`${flight.flight_number}_${flight.departure_time}`, flight);
    }
  }

  if (failures === DAY_WINDOWS.length) {
    const firstError = results.find((r) => r.status === 'rejected');
    throw firstError && firstError.reason instanceof Error
      ? firstError.reason
      : new Error('Every FareFinder time window failed for this route and date.');
  }

  const flights = Array.from(seen.values());
  flights.sort((a, b) => String(a.departure_time).localeCompare(String(b.departure_time)));
  return { flights, partial: failures > 0 };
}

module.exports = { getOneWayFares, getRoundTripFares, getDailyFlights, parseSegment };
