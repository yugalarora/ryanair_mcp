'use strict';

// Tool definitions and handlers. Every tool is read-only: it performs HTTP GET
// requests against public FareFinder endpoints and never books or changes
// anything. Ported from main.go.

const { findNearbyAirports } = require('./airports');
const { getOneWayFares, getRoundTripFares, getDailyFlights } = require('./ryanair');
const { planTrip } = require('./planner');

/** Runtime defaults, overridable by environment variables. */
function loadConfig(env) {
  const number = (raw, fallback) => {
    const parsed = Number(raw);
    return Number.isFinite(parsed) && parsed > 0 ? parsed : fallback;
  };

  return {
    defaultCurrency: env.RYANAIR_CURRENCY || env.DEFAULT_CURRENCY || 'EUR',
    defaultOrigin: env.RYANAIR_DEFAULT_ORIGIN || env.DEFAULT_ORIGIN || '',
    defaultRadiusKm: number(env.RYANAIR_RADIUS_KM || env.DEFAULT_RADIUS_KM, 150),
    defaultMaxOptions: Math.floor(number(env.RYANAIR_MAX_OPTIONS || env.DEFAULT_MAX_OPTIONS, 15)),
  };
}

/** An error whose message is meant to be shown to the agent verbatim. */
class ToolInputError extends Error {}

/** First non-empty string value among `keys`, so agents can use loose aliases. */
function pickString(args, keys) {
  for (const key of keys) {
    const value = args[key];
    if (typeof value === 'string' && value.trim() !== '') return value.trim();
  }
  return '';
}

/** Positive number among `keys`, else undefined. */
function pickNumber(args, keys) {
  for (const key of keys) {
    const value = Number(args[key]);
    if (Number.isFinite(value) && value > 0) return value;
  }
  return undefined;
}

function requireDate(value, field) {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(value)) {
    throw new ToolInputError(
      `"${field}" must be a calendar date in YYYY-MM-DD format, got "${value}".`
    );
  }
  const [year, month, day] = value.split('-').map(Number);
  const parsed = new Date(Date.UTC(year, month - 1, day));
  if (
    parsed.getUTCFullYear() !== year ||
    parsed.getUTCMonth() !== month - 1 ||
    parsed.getUTCDate() !== day
  ) {
    throw new ToolInputError(`"${field}" is not a real date: "${value}".`);
  }
  return value;
}

function requireIata(value, field) {
  const code = value.toUpperCase();
  if (!/^[A-Z]{3}$/.test(code)) {
    throw new ToolInputError(
      `"${field}" must be a 3-letter IATA airport code such as CRL, DUB or FAO, got "${value}". ` +
        'Use find_nearby_airports to look up a code from a city or region name.'
    );
  }
  return code;
}

function requirePresent(value, field, hint) {
  if (!value) {
    throw new ToolInputError(`"${field}" is required.${hint ? ' ' + hint : ''}`);
  }
  return value;
}

function normaliseCurrency(value, fallback) {
  const code = (value || fallback || 'EUR').toUpperCase();
  if (!/^[A-Z]{3}$/.test(code)) {
    throw new ToolInputError(
      `"currency" must be a 3-letter ISO currency code such as EUR, GBP or PLN, got "${value}".`
    );
  }
  return code;
}

function requireTime(value, field) {
  if (!value) return '';
  if (!/^([01]\d|2[0-3]):[0-5]\d$/.test(value)) {
    throw new ToolInputError(`"${field}" must be a 24-hour time in HH:MM format, got "${value}".`);
  }
  return value;
}

const ALIASES = {
  origin: ['origin', 'from', 'departure', 'orig'],
  destination: ['destination', 'to', 'arrival', 'dest'],
};

/** Build the tool list and handlers against a loaded config. */
function buildTools(config) {
  const currencyProperty = {
    type: 'string',
    description: `ISO 4217 currency code for returned fares (default: ${config.defaultCurrency}).`,
  };

  return [
    {
      name: 'plan_trip',
      title: 'Plan a multi-airport trip',
      description:
        'Plan a one-way or round-trip Ryanair journey between two places given as city names, ' +
        'region names, or IATA codes. Expands each end into the airports within a radius ' +
        '(so "Brussels" covers CRL, BRU and LGG, and "Andalusia" covers AGP, SVQ and XRY), ' +
        'searches every viable airport pairing, and returns journeys ranked by total price. ' +
        'Use this for open-ended requests where the traveller has not fixed an airport; use ' +
        'search_flights or search_return_flights when both airport codes are already known.',
      annotations: { title: 'Plan a multi-airport trip', readOnlyHint: true, idempotentHint: true, openWorldHint: true },
      inputSchema: {
        type: 'object',
        properties: {
          origin: {
            type: 'string',
            description: "Departure city, region, or IATA code, e.g. 'Brussels', 'London', 'DUB'.",
          },
          destination: {
            type: 'string',
            description: "Arrival city, region, or IATA code, e.g. 'Algarve', 'Andalusia', 'BCN'.",
          },
          outbound_date_from: {
            type: 'string',
            description: 'Earliest outbound departure date, YYYY-MM-DD.',
          },
          outbound_date_to: {
            type: 'string',
            description: 'Latest outbound departure date, YYYY-MM-DD. Defaults to outbound_date_from.',
          },
          return_date_from: {
            type: 'string',
            description: 'Earliest return date, YYYY-MM-DD. Omit for a one-way trip.',
          },
          return_date_to: {
            type: 'string',
            description: 'Latest return date, YYYY-MM-DD. Defaults to return_date_from.',
          },
          origin_radius_km: {
            type: 'number',
            description: `Radius in km around the origin to include alternative departure airports (default: ${config.defaultRadiusKm}).`,
          },
          destination_radius_km: {
            type: 'number',
            description: `Radius in km around the destination to include alternative arrival airports (default: ${config.defaultRadiusKm}).`,
          },
          max_price: {
            type: 'number',
            description: 'Optional upper price limit per journey, in the requested currency.',
          },
          currency: currencyProperty,
          max_options: {
            type: 'number',
            description: `Maximum ranked journeys to return (default: ${config.defaultMaxOptions}).`,
          },
        },
        required: ['origin', 'destination', 'outbound_date_from'],
        additionalProperties: true,
      },
      async handler(args) {
        const origin = pickString(args, ALIASES.origin) || config.defaultOrigin;
        requirePresent(
          origin,
          'origin',
          'Give a city, region, or IATA code, or set RYANAIR_DEFAULT_ORIGIN.'
        );

        const destination = requirePresent(pickString(args, ALIASES.destination), 'destination');

        const outboundDateFrom = requireDate(
          requirePresent(
            pickString(args, [
              'outbound_date_from',
              'outbound_date',
              'date_from',
              'date',
              'departure_date',
            ]),
            'outbound_date_from'
          ),
          'outbound_date_from'
        );

        const outboundDateToRaw = pickString(args, ['outbound_date_to', 'date_to']);
        const outboundDateTo = outboundDateToRaw
          ? requireDate(outboundDateToRaw, 'outbound_date_to')
          : outboundDateFrom;

        const returnDateFromRaw = pickString(args, [
          'return_date_from',
          'return_date',
          'inbound_date_from',
          'inbound_date',
        ]);
        const returnDateFrom = returnDateFromRaw
          ? requireDate(returnDateFromRaw, 'return_date_from')
          : '';

        const returnDateToRaw = pickString(args, ['return_date_to', 'inbound_date_to']);
        const returnDateTo = returnDateToRaw
          ? requireDate(returnDateToRaw, 'return_date_to')
          : returnDateFrom;

        if (outboundDateTo < outboundDateFrom) {
          throw new ToolInputError(
            `"outbound_date_to" (${outboundDateTo}) is before "outbound_date_from" (${outboundDateFrom}).`
          );
        }
        if (returnDateFrom && returnDateFrom < outboundDateFrom) {
          throw new ToolInputError(
            `"return_date_from" (${returnDateFrom}) is before "outbound_date_from" (${outboundDateFrom}).`
          );
        }

        return planTrip({
          origin,
          destination,
          outboundDateFrom,
          outboundDateTo,
          returnDateFrom,
          returnDateTo,
          originRadiusKm: pickNumber(args, ['origin_radius_km']) || config.defaultRadiusKm,
          destinationRadiusKm:
            pickNumber(args, ['destination_radius_km']) || config.defaultRadiusKm,
          maxPrice: pickNumber(args, ['max_price']),
          currency: normaliseCurrency(pickString(args, ['currency']), config.defaultCurrency),
          maxOptions: pickNumber(args, ['max_options']) || config.defaultMaxOptions,
        });
      },
    },

    {
      name: 'search_flights',
      title: 'Search one-way flights',
      description:
        'Search the cheapest direct one-way Ryanair flights leaving a specific airport within a ' +
        'date and time window. Returns one cheapest flight per destination and date. Omit ' +
        'destination to list every destination served from that airport, which is how to find ' +
        'the cheapest places to fly to on a given date.',
      annotations: { title: 'Search one-way flights', readOnlyHint: true, idempotentHint: true, openWorldHint: true },
      inputSchema: {
        type: 'object',
        properties: {
          origin: { type: 'string', description: "3-letter IATA departure code, e.g. 'CRL'." },
          destination: {
            type: 'string',
            description: "Optional 3-letter IATA arrival code, e.g. 'FAO'. Omit to search all destinations.",
          },
          date_from: { type: 'string', description: 'Earliest departure date, YYYY-MM-DD.' },
          date_to: {
            type: 'string',
            description: 'Latest departure date, YYYY-MM-DD. Defaults to date_from.',
          },
          time_from: { type: 'string', description: "Earliest departure time, HH:MM (default '00:00')." },
          time_to: { type: 'string', description: "Latest departure time, HH:MM (default '23:59')." },
          max_price: { type: 'number', description: 'Optional upper price limit per flight.' },
          currency: currencyProperty,
        },
        required: ['origin', 'date_from'],
        additionalProperties: true,
      },
      async handler(args) {
        const origin = requireIata(
          requirePresent(pickString(args, ALIASES.origin), 'origin'),
          'origin'
        );

        const destinationRaw = pickString(args, ALIASES.destination);
        const destination = destinationRaw ? requireIata(destinationRaw, 'destination') : '';

        const dateFrom = requireDate(
          requirePresent(
            pickString(args, [
              'date_from',
              'date',
              'departure_date',
              'outbound_date_from',
              'outbound_date',
            ]),
            'date_from'
          ),
          'date_from'
        );

        const dateToRaw = pickString(args, ['date_to', 'outbound_date_to']);
        const dateTo = dateToRaw ? requireDate(dateToRaw, 'date_to') : dateFrom;

        if (dateTo < dateFrom) {
          throw new ToolInputError(`"date_to" (${dateTo}) is before "date_from" (${dateFrom}).`);
        }

        const currency = normaliseCurrency(pickString(args, ['currency']), config.defaultCurrency);

        const flights = await getOneWayFares({
          origin,
          destination,
          dateFrom,
          dateTo,
          timeFrom: requireTime(pickString(args, ['time_from']), 'time_from'),
          timeTo: requireTime(pickString(args, ['time_to']), 'time_to'),
          currency,
          maxPrice: pickNumber(args, ['max_price']),
        });

        return {
          origin,
          destination: destination || 'all destinations',
          date_range: { from: dateFrom, to: dateTo },
          currency,
          count: flights.length,
          flights,
        };
      },
    },

    {
      name: 'search_return_flights',
      title: 'Search round-trip flights',
      description:
        'Search the cheapest direct round-trip Ryanair fare combinations between two airports ' +
        'across an outbound and an inbound date window. Returns outbound and inbound pairs with ' +
        'their combined price. Use plan_trip instead when the traveller is flexible about which ' +
        'airport to use.',
      annotations: { title: 'Search round-trip flights', readOnlyHint: true, idempotentHint: true, openWorldHint: true },
      inputSchema: {
        type: 'object',
        properties: {
          origin: { type: 'string', description: "3-letter IATA departure code, e.g. 'DUB'." },
          destination: {
            type: 'string',
            description: "Optional 3-letter IATA arrival code, e.g. 'FAO'. Omit to search all destinations.",
          },
          outbound_date_from: { type: 'string', description: 'Earliest outbound date, YYYY-MM-DD.' },
          outbound_date_to: {
            type: 'string',
            description: 'Latest outbound date, YYYY-MM-DD. Defaults to outbound_date_from.',
          },
          return_date_from: { type: 'string', description: 'Earliest return date, YYYY-MM-DD.' },
          return_date_to: {
            type: 'string',
            description: 'Latest return date, YYYY-MM-DD. Defaults to return_date_from.',
          },
          max_price: { type: 'number', description: 'Optional upper limit on the total round-trip price.' },
          currency: currencyProperty,
        },
        required: ['origin', 'outbound_date_from', 'return_date_from'],
        additionalProperties: true,
      },
      async handler(args) {
        const origin = requireIata(
          requirePresent(pickString(args, ALIASES.origin), 'origin'),
          'origin'
        );

        const destinationRaw = pickString(args, ALIASES.destination);
        const destination = destinationRaw ? requireIata(destinationRaw, 'destination') : '';

        const outDateFrom = requireDate(
          requirePresent(
            pickString(args, ['outbound_date_from', 'outbound_date', 'date_from', 'date']),
            'outbound_date_from'
          ),
          'outbound_date_from'
        );

        const outDateToRaw = pickString(args, ['outbound_date_to', 'date_to']);
        const outDateTo = outDateToRaw
          ? requireDate(outDateToRaw, 'outbound_date_to')
          : outDateFrom;

        const retDateFrom = requireDate(
          requirePresent(
            pickString(args, [
              'return_date_from',
              'return_date',
              'inbound_date_from',
              'inbound_date',
            ]),
            'return_date_from',
            'Use search_flights for a one-way search.'
          ),
          'return_date_from'
        );

        const retDateToRaw = pickString(args, ['return_date_to', 'inbound_date_to']);
        const retDateTo = retDateToRaw ? requireDate(retDateToRaw, 'return_date_to') : retDateFrom;

        if (retDateFrom < outDateFrom) {
          throw new ToolInputError(
            `"return_date_from" (${retDateFrom}) is before "outbound_date_from" (${outDateFrom}).`
          );
        }

        const currency = normaliseCurrency(pickString(args, ['currency']), config.defaultCurrency);

        const trips = await getRoundTripFares({
          origin,
          destination,
          outDateFrom,
          outDateTo,
          retDateFrom,
          retDateTo,
          currency,
          maxPrice: pickNumber(args, ['max_price']),
        });

        trips.sort((a, b) => a.total_price - b.total_price);

        return {
          origin,
          destination: destination || 'all destinations',
          outbound_date_range: { from: outDateFrom, to: outDateTo },
          return_date_range: { from: retDateFrom, to: retDateTo },
          currency,
          count: trips.length,
          trips,
        };
      },
    },

    {
      name: 'get_daily_flights',
      title: 'Get a full day of flights',
      description:
        'List every scheduled Ryanair flight on one route for one date, with prices and times. ' +
        'FareFinder returns only the single cheapest flight per time window, so this queries ' +
        'four windows across the day and merges them to recover the morning, afternoon and ' +
        'evening departures. Use it once a route is chosen and the traveller needs departure times.',
      annotations: { title: 'Get a full day of flights', readOnlyHint: true, idempotentHint: true, openWorldHint: true },
      inputSchema: {
        type: 'object',
        properties: {
          origin: { type: 'string', description: "3-letter IATA departure code, e.g. 'CRL'." },
          destination: { type: 'string', description: "3-letter IATA arrival code, e.g. 'FAO'." },
          date: { type: 'string', description: 'Flight date, YYYY-MM-DD.' },
          currency: currencyProperty,
        },
        required: ['origin', 'destination', 'date'],
        additionalProperties: true,
      },
      async handler(args) {
        const origin = requireIata(
          requirePresent(pickString(args, ALIASES.origin), 'origin'),
          'origin'
        );
        const destination = requireIata(
          requirePresent(pickString(args, ALIASES.destination), 'destination'),
          'destination'
        );
        const date = requireDate(
          requirePresent(
            pickString(args, [
              'date',
              'date_from',
              'departure_date',
              'outbound_date_from',
              'outbound_date',
            ]),
            'date'
          ),
          'date'
        );
        const currency = normaliseCurrency(pickString(args, ['currency']), config.defaultCurrency);

        const { flights, partial } = await getDailyFlights({
          origin,
          destination,
          date,
          currency,
        });

        return {
          origin,
          destination,
          date,
          currency,
          count: flights.length,
          flights,
          note: partial
            ? 'Some time windows could not be queried, so this schedule may be incomplete.'
            : undefined,
        };
      },
    },

    {
      name: 'find_nearby_airports',
      title: 'Find nearby airports',
      description:
        'List the Ryanair-served airports within a radius of a city, region, or IATA code, ' +
        'sorted by great-circle distance. Matches against a built-in directory of about 120 ' +
        'European and North African airports, including region aliases such as "Algarve", ' +
        '"Andalusia" and "Puglia". Use it to turn a place name into airport codes, or to find ' +
        'cheaper secondary airports near a traveller.',
      annotations: { title: 'Find nearby airports', readOnlyHint: true, idempotentHint: true, openWorldHint: false },
      inputSchema: {
        type: 'object',
        properties: {
          query: {
            type: 'string',
            description: "City, region, or 3-letter IATA code, e.g. 'Brussels', 'Algarve', 'CRL'.",
          },
          radius_km: {
            type: 'number',
            description: `Search radius in km (default: ${config.defaultRadiusKm}).`,
          },
          limit: { type: 'number', description: 'Maximum airports to return (default: 10).' },
        },
        required: ['query'],
        additionalProperties: true,
      },
      async handler(args) {
        const query = requirePresent(
          pickString(args, ['query', 'city', 'airport', 'origin', 'destination', 'location']),
          'query'
        );

        const radiusKm = pickNumber(args, ['radius_km']) || config.defaultRadiusKm;
        const limit = pickNumber(args, ['limit']) || 10;

        const resolved = findNearbyAirports(query, radiusKm, limit);
        if (!resolved) {
          return {
            query,
            radius_km: radiusKm,
            airports_found: 0,
            nearby_airports: [],
            note:
              `No airport, city, or region matching "${query}" is in the built-in directory. ` +
              'It covers Ryanair-served airports in Europe and North Africa. Try a larger ' +
              'nearby city, or pass a 3-letter IATA code directly to the search tools.',
          };
        }

        return {
          query,
          anchor_airport: {
            iata: resolved.anchor.iata,
            name: resolved.anchor.name,
            city: resolved.anchor.city,
            country: resolved.anchor.country,
          },
          radius_km: radiusKm,
          airports_found: resolved.nearby.length,
          nearby_airports: resolved.nearby,
        };
      },
    },
  ];
}

module.exports = { buildTools, loadConfig, ToolInputError };
