package planner

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"ryanair-mcp/internal/airports"
	"ryanair-mcp/internal/ryanair"
)

// TripRequest defines the criteria for the multi-airport trip planner.
type TripRequest struct {
	Origin               string   `json:"origin"`                 // City name, region, or IATA code
	Destination          string   `json:"destination"`            // City name, region, or IATA code
	OriginRadiusKm       float64  `json:"origin_radius_km"`       // Search radius around origin (default: 150 km)
	DestinationRadiusKm  float64  `json:"destination_radius_km"`  // Search radius around destination (default: 150 km)
	OutboundDateFrom     string   `json:"outbound_date_from"`     // YYYY-MM-DD
	OutboundDateTo       string   `json:"outbound_date_to"`       // YYYY-MM-DD
	ReturnDateFrom       string   `json:"return_date_from"`       // Optional YYYY-MM-DD for return flights
	ReturnDateTo         string   `json:"return_date_to"`         // Optional YYYY-MM-DD for return flights
	MaxPrice             *float64 `json:"max_price"`              // Optional ceiling on price
	Currency             string   `json:"currency"`               // Currency code (default: EUR)
	MaxOptions           int      `json:"max_options"`            // Maximum itineraries to return (default: 15)
}

// AirportDistance notes an airport and how far it is from the requested reference point.
type AirportDistance struct {
	IATA       string  `json:"iata"`
	Name       string  `json:"name"`
	City       string  `json:"city"`
	Country    string  `json:"country"`
	DistanceKm float64 `json:"distance_km"`
}

// JourneyOption models a complete planned journey recommendation.
type JourneyOption struct {
	Type               string          `json:"type"` // "one_way" or "round_trip"
	TotalPrice         float64         `json:"total_price"`
	Currency           string          `json:"currency"`
	OriginAirport      AirportDistance `json:"origin_airport"`
	DestinationAirport AirportDistance `json:"destination_airport"`
	OutboundFlight     ryanair.Flight  `json:"outbound_flight"`
	ReturnFlight       *ryanair.Flight `json:"return_flight,omitempty"`
	Summary            string          `json:"summary"`
}

// PlanResult is the complete payload returned to the agent.
type PlanResult struct {
	OriginAnchor         string            `json:"origin_anchor"`
	DestinationAnchor    string            `json:"destination_anchor"`
	OriginAirportsFound  []AirportDistance `json:"origin_airports_checked"`
	DestAirportsFound    []AirportDistance `json:"dest_airports_checked"`
	JourneysFound        int               `json:"journeys_found"`
	RecommendedJourneys  []JourneyOption   `json:"recommended_journeys"`
}

// Planner coordinates multi-airport route discovery.
type Planner struct {
	client *ryanair.Client
}

// NewPlanner constructs a new Planner instance.
func NewPlanner(client *ryanair.Client) *Planner {
	if client == nil {
		client = ryanair.NewClient()
	}
	return &Planner{client: client}
}

// PlanTrip resolves nearby airports for both endpoints and searches all viable route pairings.
func (p *Planner) PlanTrip(ctx context.Context, req TripRequest) (*PlanResult, error) {
	if req.Currency == "" {
		req.Currency = "EUR"
	}
	if req.OriginRadiusKm <= 0 {
		req.OriginRadiusKm = 150.0
	}
	if req.DestinationRadiusKm <= 0 {
		req.DestinationRadiusKm = 150.0
	}
	if req.MaxOptions <= 0 {
		req.MaxOptions = 15
	}
	if req.OutboundDateTo == "" {
		req.OutboundDateTo = req.OutboundDateFrom
	}

	// 1. Resolve candidate origin airports
	nearbyOrigins, origAnchor, origFound := airports.FindNearbyAirports(req.Origin, req.OriginRadiusKm, 5)
	var originDistances []AirportDistance
	if origFound {
		for _, o := range nearbyOrigins {
			originDistances = append(originDistances, AirportDistance{
				IATA:       o.IATA,
				Name:       o.Name,
				City:       o.City,
				Country:    o.Country,
				DistanceKm: o.DistanceKm,
			})
		}
	} else {
		// Treat as direct IATA
		originDistances = append(originDistances, AirportDistance{
			IATA:       req.Origin,
			Name:       req.Origin,
			DistanceKm: 0,
		})
	}

	// 2. Resolve candidate destination airports
	nearbyDests, destAnchor, destFound := airports.FindNearbyAirports(req.Destination, req.DestinationRadiusKm, 5)
	var destDistances []AirportDistance
	if destFound {
		for _, d := range nearbyDests {
			destDistances = append(destDistances, AirportDistance{
				IATA:       d.IATA,
				Name:       d.Name,
				City:       d.City,
				Country:    d.Country,
				DistanceKm: d.DistanceKm,
			})
		}
	} else {
		// Treat as direct IATA
		destDistances = append(destDistances, AirportDistance{
			IATA:       req.Destination,
			Name:       req.Destination,
			DistanceKm: 0,
		})
	}

	isRoundTrip := req.ReturnDateFrom != ""
	if isRoundTrip && req.ReturnDateTo == "" {
		req.ReturnDateTo = req.ReturnDateFrom
	}

	// 3. Concurrently search route combinations across Origins x Destinations
	type pair struct {
		orig AirportDistance
		dest AirportDistance
	}

	var pairs []pair
	for _, o := range originDistances {
		for _, d := range destDistances {
			if o.IATA != d.IATA {
				pairs = append(pairs, pair{orig: o, dest: d})
			}
		}
	}

	var mu sync.Mutex
	var journeys []JourneyOption

	// Use worker pool of 4 concurrent workers to be gentle on Ryanair API
	sem := make(chan struct{}, 4)
	var wg sync.WaitGroup

	for _, pr := range pairs {
		wg.Add(1)
		go func(item pair) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if isRoundTrip {
				trips, err := p.p_searchRoundTrip(ctx, item.orig.IATA, item.dest.IATA, req)
				if err == nil && len(trips) > 0 {
					mu.Lock()
					for _, t := range trips {
						summary := fmt.Sprintf("Round-trip: %s (%s) -> %s (%s) on %s, returning on %s",
							item.orig.City, item.orig.IATA, item.dest.City, item.dest.IATA,
							t.Outbound.DepartureTime.Format("02 Jan 15:04"),
							t.Inbound.DepartureTime.Format("02 Jan 15:04"))

						if item.orig.DistanceKm > 0 {
							summary += fmt.Sprintf(" [%.0fkm from %s]", item.orig.DistanceKm, origAnchor.City)
						}
						if item.dest.DistanceKm > 0 {
							summary += fmt.Sprintf(" [%.0fkm from %s]", item.dest.DistanceKm, destAnchor.City)
						}

						inboundCopy := t.Inbound
						journeys = append(journeys, JourneyOption{
							Type:               "round_trip",
							TotalPrice:         t.TotalPrice,
							Currency:           t.Currency,
							OriginAirport:      item.orig,
							DestinationAirport: item.dest,
							OutboundFlight:     t.Outbound,
							ReturnFlight:       &inboundCopy,
							Summary:            summary,
						})
					}
					mu.Unlock()
				}
			} else {
				flights, err := p.p_searchOneWay(ctx, item.orig.IATA, item.dest.IATA, req)
				if err == nil && len(flights) > 0 {
					mu.Lock()
					for _, f := range flights {
						summary := fmt.Sprintf("One-way: %s (%s) -> %s (%s) on %s for %.2f %s",
							item.orig.City, item.orig.IATA, item.dest.City, item.dest.IATA,
							f.DepartureTime.Format("02 Jan 15:04"),
							f.Price, f.Currency)

						if item.orig.DistanceKm > 0 {
							summary += fmt.Sprintf(" [%.0fkm from %s]", item.orig.DistanceKm, origAnchor.City)
						}
						if item.dest.DistanceKm > 0 {
							summary += fmt.Sprintf(" [%.0fkm from %s]", item.dest.DistanceKm, destAnchor.City)
						}

						journeys = append(journeys, JourneyOption{
							Type:               "one_way",
							TotalPrice:         f.Price,
							Currency:           f.Currency,
							OriginAirport:      item.orig,
							DestinationAirport: item.dest,
							OutboundFlight:     f,
							Summary:            summary,
						})
					}
					mu.Unlock()
				}
			}
		}(pr)
	}

	wg.Wait()

	// Sort journeys by total price ascending
	sort.Slice(journeys, func(i, j int) bool {
		return journeys[i].TotalPrice < journeys[j].TotalPrice
	})

	totalFound := len(journeys)
	if len(journeys) > req.MaxOptions {
		journeys = journeys[:req.MaxOptions]
	}

	origAnchorName := req.Origin
	if origFound {
		origAnchorName = fmt.Sprintf("%s (%s, %s)", origAnchor.City, origAnchor.IATA, origAnchor.Country)
	}

	destAnchorName := req.Destination
	if destFound {
		destAnchorName = fmt.Sprintf("%s (%s, %s)", destAnchor.City, destAnchor.IATA, destAnchor.Country)
	}

	return &PlanResult{
		OriginAnchor:        origAnchorName,
		DestinationAnchor:   destAnchorName,
		OriginAirportsFound: originDistances,
		DestAirportsFound:   destDistances,
		JourneysFound:       totalFound,
		RecommendedJourneys: journeys,
	}, nil
}

func (p *Planner) p_searchOneWay(ctx context.Context, orig, dest string, req TripRequest) ([]ryanair.Flight, error) {
	return p.client.GetOneWayFares(
		ctx,
		orig,
		dest,
		req.OutboundDateFrom,
		req.OutboundDateTo,
		"00:00",
		"23:59",
		req.Currency,
		req.MaxPrice,
	)
}

func (p *Planner) p_searchRoundTrip(ctx context.Context, orig, dest string, req TripRequest) ([]ryanair.Trip, error) {
	return p.client.GetRoundTripFares(
		ctx,
		orig,
		dest,
		req.OutboundDateFrom,
		req.OutboundDateTo,
		req.ReturnDateFrom,
		req.ReturnDateTo,
		req.Currency,
		req.MaxPrice,
	)
}
