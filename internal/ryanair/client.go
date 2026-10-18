package ryanair

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"
)

const (
	DefaultBaseURL = "https://services-api.ryanair.com/farfnd/v4/"
	DefaultTimeout = 15 * time.Second
)

// Client executes HTTP requests against Ryanair's Farfnd V4 FareFinder REST APIs.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Ryanair API client.
func NewClient() *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// Flight represents a single flight segment with schedule and pricing details.
type Flight struct {
	Origin             string    `json:"origin"`
	OriginName         string    `json:"origin_name"`
	OriginCity         string    `json:"origin_city"`
	OriginCountry      string    `json:"origin_country"`
	Destination        string    `json:"destination"`
	DestinationName    string    `json:"destination_name"`
	DestinationCity    string    `json:"destination_city"`
	DestinationCountry string    `json:"destination_country"`
	DepartureTime      time.Time `json:"departure_time"`
	ArrivalTime        time.Time `json:"arrival_time"`
	FlightNumber       string    `json:"flight_number"`
	Price              float64   `json:"price"`
	Currency           string    `json:"currency"`
}

// Trip represents a round-trip journey with outbound and return flights.
type Trip struct {
	Outbound   Flight  `json:"outbound"`
	Inbound    Flight  `json:"inbound"`
	TotalPrice float64 `json:"total_price"`
	Currency   string  `json:"currency"`
}

// rawFlightSegment models the JSON response from farfnd/v4.
type rawFlightSegment struct {
	DepartureAirport struct {
		CountryName string `json:"countryName"`
		IataCode    string `json:"iataCode"`
		Name        string `json:"name"`
		City        struct {
			Name string `json:"name"`
		} `json:"city"`
	} `json:"departureAirport"`
	ArrivalAirport struct {
		CountryName string `json:"countryName"`
		IataCode    string `json:"iataCode"`
		Name        string `json:"name"`
		City        struct {
			Name string `json:"name"`
		} `json:"city"`
	} `json:"arrivalAirport"`
	DepartureDate string `json:"departureDate"`
	ArrivalDate   string `json:"arrivalDate"`
	FlightNumber  string `json:"flightNumber"`
	Price         struct {
		Value        float64 `json:"value"`
		CurrencyCode string  `json:"currencyCode"`
	} `json:"price"`
}

type rawOneWayResponse struct {
	Fares []struct {
		Outbound rawFlightSegment `json:"outbound"`
	} `json:"fares"`
}

type rawRoundTripResponse struct {
	Fares []struct {
		Outbound rawFlightSegment `json:"outbound"`
		Inbound  rawFlightSegment `json:"inbound"`
		Summary  struct {
			Price struct {
				Value        float64 `json:"value"`
				CurrencyCode string  `json:"currencyCode"`
			} `json:"price"`
		} `json:"summary"`
	} `json:"fares"`
}

func parseSegment(raw rawFlightSegment) Flight {
	depTime, _ := time.Parse("2006-01-02T15:04:05", raw.DepartureDate)
	arrTime, _ := time.Parse("2006-01-02T15:04:05", raw.ArrivalDate)

	fn := raw.FlightNumber
	if len(fn) > 2 && fn[:2] == "FR" && fn[2] != ' ' {
		fn = "FR " + fn[2:]
	}

	return Flight{
		Origin:             raw.DepartureAirport.IataCode,
		OriginName:         raw.DepartureAirport.Name,
		OriginCity:         raw.DepartureAirport.City.Name,
		OriginCountry:      raw.DepartureAirport.CountryName,
		Destination:        raw.ArrivalAirport.IataCode,
		DestinationName:    raw.ArrivalAirport.Name,
		DestinationCity:    raw.ArrivalAirport.City.Name,
		DestinationCountry: raw.ArrivalAirport.CountryName,
		DepartureTime:      depTime,
		ArrivalTime:        arrTime,
		FlightNumber:       fn,
		Price:              raw.Price.Value,
		Currency:           raw.Price.CurrencyCode,
	}
}

// GetOneWayFares calls the oneWayFares endpoint.
func (c *Client) GetOneWayFares(
	ctx context.Context,
	origin string,
	destination string,
	dateFrom string,
	dateTo string,
	timeFrom string,
	timeTo string,
	currency string,
	maxPrice *float64,
) ([]Flight, error) {
	endpoint := c.BaseURL + "oneWayFares"
	params := url.Values{}
	params.Set("departureAirportIataCode", origin)
	params.Set("outboundDepartureDateFrom", dateFrom)
	params.Set("outboundDepartureDateTo", dateTo)

	if timeFrom == "" {
		timeFrom = "00:00"
	}
	if timeTo == "" {
		timeTo = "23:59"
	}
	params.Set("outboundDepartureTimeFrom", timeFrom)
	params.Set("outboundDepartureTimeTo", timeTo)

	if currency != "" {
		params.Set("currency", currency)
	} else {
		params.Set("currency", "EUR")
	}

	if destination != "" {
		params.Set("arrivalAirportIataCode", destination)
	}
	if maxPrice != nil && *maxPrice > 0 {
		params.Set("priceValueTo", fmt.Sprintf("%.0f", *maxPrice))
	}

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// Browser user agent to avoid CDN blocks
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var parsed rawOneWayResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	var results []Flight
	for _, item := range parsed.Fares {
		results = append(results, parseSegment(item.Outbound))
	}

	return results, nil
}

// GetRoundTripFares calls the roundTripFares endpoint.
func (c *Client) GetRoundTripFares(
	ctx context.Context,
	origin string,
	destination string,
	outDateFrom string,
	outDateTo string,
	retDateFrom string,
	retDateTo string,
	currency string,
	maxPrice *float64,
) ([]Trip, error) {
	endpoint := c.BaseURL + "roundTripFares"
	params := url.Values{}
	params.Set("departureAirportIataCode", origin)
	params.Set("outboundDepartureDateFrom", outDateFrom)
	params.Set("outboundDepartureDateTo", outDateTo)
	params.Set("inboundDepartureDateFrom", retDateFrom)
	params.Set("inboundDepartureDateTo", retDateTo)
	params.Set("outboundDepartureTimeFrom", "00:00")
	params.Set("outboundDepartureTimeTo", "23:59")
	params.Set("inboundDepartureTimeFrom", "00:00")
	params.Set("inboundDepartureTimeTo", "23:59")

	if currency != "" {
		params.Set("currency", currency)
	} else {
		params.Set("currency", "EUR")
	}

	if destination != "" {
		params.Set("arrivalAirportIataCode", destination)
	}
	if maxPrice != nil && *maxPrice > 0 {
		params.Set("priceValueTo", fmt.Sprintf("%.0f", *maxPrice))
	}

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var parsed rawRoundTripResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	var results []Trip
	for _, item := range parsed.Fares {
		out := parseSegment(item.Outbound)
		in := parseSegment(item.Inbound)
		total := out.Price + in.Price
		if item.Summary.Price.Value > 0 {
			total = item.Summary.Price.Value
		}

		results = append(results, Trip{
			Outbound:   out,
			Inbound:    in,
			TotalPrice: total,
			Currency:   out.Currency,
		})
	}

	return results, nil
}

// GetDailyFlights scans multiple departure time slices across a single day
// to retrieve all distinct scheduled flights on that route (since FareFinder only returns 1 per time window).
func (c *Client) GetDailyFlights(ctx context.Context, origin, dest, date, currency string) ([]Flight, error) {
	timeWindows := [][2]string{
		{"00:00", "11:00"},
		{"11:01", "15:30"},
		{"15:31", "19:30"},
		{"19:31", "23:59"},
	}

	type result struct {
		flights []Flight
		err     error
	}

	ch := make(chan result, len(timeWindows))
	var wg sync.WaitGroup

	for _, w := range timeWindows {
		wg.Add(1)
		go func(tFrom, tTo string) {
			defer wg.Done()
			fls, err := c.GetOneWayFares(ctx, origin, dest, date, date, tFrom, tTo, currency, nil)
			ch <- result{flights: fls, err: err}
		}(w[0], w[1])
	}

	wg.Wait()
	close(ch)

	seen := make(map[string]Flight)
	for res := range ch {
		if res.err == nil {
			for _, f := range res.flights {
				key := fmt.Sprintf("%s_%s", f.FlightNumber, f.DepartureTime.Format(time.RFC3339))
				seen[key] = f
			}
		}
	}

	var distinct []Flight
	for _, f := range seen {
		distinct = append(distinct, f)
	}

	sort.Slice(distinct, func(i, j int) bool {
		return distinct[i].DepartureTime.Before(distinct[j].DepartureTime)
	})

	return distinct, nil
}
