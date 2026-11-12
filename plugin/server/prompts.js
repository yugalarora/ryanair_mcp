'use strict';

// User-invoked MCP prompts. These are workflows the user explicitly chooses;
// they are not injected into unrelated conversations.

function buildPrompts(config) {
  return [
    {
      name: 'plan_vacation',
      title: 'Plan a vacation',
      description:
        'Walk through planning a trip: confirm the traveller\'s constraints, then search ' +
        'alternative airports and rank the resulting itineraries.',
      arguments: [
        {
          name: 'origin',
          description: "Departure city or airport code, e.g. 'Brussels', 'CRL'.",
          required: true,
        },
        {
          name: 'destination',
          description: "Destination city or region, e.g. 'Andalusia', 'Algarve', 'BCN'.",
          required: true,
        },
        { name: 'departure_date', description: 'Target departure date, YYYY-MM-DD.', required: true },
        { name: 'return_date', description: 'Optional return date, YYYY-MM-DD.', required: false },
        { name: 'currency', description: 'Preferred currency code.', required: false },
      ],
      build(args) {
        const origin = args.origin || '';
        const destination = args.destination || '';
        const departureDate = args.departure_date || '';
        const returnDate = args.return_date || 'not specified (treat as one-way unless told otherwise)';
        const currency = args.currency || config.defaultCurrency;

        return {
          description: 'Vacation planning workflow',
          text: `Plan a flight itinerary from ${origin} to ${destination}, departing ${departureDate}.
Return date: ${returnDate}
Currency: ${currency}

Steps:
1. Ask me about anything that would change the search, unless I have already said:
   - Are my dates flexible by one to three days? Adjacent days often differ a lot in price.
   - Must the return flight land at the same airport I left from, or is an open-jaw
     route acceptable (for example out to Malaga, back from Seville)?
   - Are early-morning or late-night departures acceptable?
2. Use find_nearby_airports on both ${origin} and ${destination} to see which airports are in range.
3. Use plan_trip to search the viable airport pairings.
4. For the best one or two routes, use get_daily_flights to show the full timetable for that day.
5. Summarise the options as a table ranked by total price, noting for each how far the
   airport is from where I actually want to be, and the departure and arrival times.`,
        };
      },
    },

    {
      name: 'weekend_getaway',
      title: 'Find a weekend getaway',
      description:
        'Find the cheapest weekend trips from a home airport by scanning every destination ' +
        'served on the chosen dates.',
      arguments: [
        {
          name: 'origin',
          description: "Home airport IATA code or city, e.g. 'DUB', 'CRL', 'STN'.",
          required: true,
        },
        { name: 'date_from', description: 'Earliest outbound date, YYYY-MM-DD.', required: true },
        { name: 'date_to', description: 'Optional latest outbound date, YYYY-MM-DD.', required: false },
      ],
      build(args) {
        const origin = args.origin || '';
        const dateFrom = args.date_from || '';
        const dateTo = args.date_to || dateFrom;

        return {
          description: 'Weekend getaway finder',
          text: `Find the cheapest weekend getaways from ${origin} departing between ${dateFrom} and ${dateTo}, priced in ${config.defaultCurrency}.

Steps:
1. Use search_flights from ${origin} with no destination, to list every destination served
   on those dates with its cheapest fare.
2. Take the five cheapest destinations.
3. For each, use search_return_flights to price a return on the Sunday or Monday after.
4. Present the best three to five complete round trips: destination, both flight times,
   and total price. Mention if any fare looks like an error or a one-off.`,
        };
      },
    },
  ];
}

module.exports = { buildPrompts };
