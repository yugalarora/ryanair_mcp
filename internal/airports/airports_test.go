package airports

import (
	"testing"
)

func TestFindNearbyAirports(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		radiusKm      float64
		limit         int
		expectAnchor  string
		expectCity    string
		minCount      int
	}{
		{
			name:         "Brussels hub lookup",
			query:        "Brussels",
			radiusKm:     150,
			limit:        5,
			expectCity:   "Brussels",
			minCount:     2,
		},
		{
			name:         "Direct IATA code DUB",
			query:        "DUB",
			radiusKm:     100,
			limit:        5,
			expectAnchor: "DUB",
			minCount:     1,
		},
		{
			name:       "London area airports",
			query:      "London",
			radiusKm:   120,
			limit:      10,
			expectCity: "London",
			minCount:   3,
		},
		{
			name:     "Regional alias Andalusia",
			query:    "Andalusia",
			radiusKm: 250,
			limit:    5,
			minCount: 2, // Finds SVQ, AGP, XRY
		},
		{
			name:     "Destination alias Tarifa",
			query:    "Tarifa",
			radiusKm: 150,
			limit:    5,
			minCount: 1, // Finds XRY / AGP
		},
		{
			name:     "Regional alias Puglia",
			query:    "Puglia",
			radiusKm: 150,
			limit:    5,
			minCount: 2, // Finds BRI and BDS
		},
		{
			name:         "New airport NRN Weeze",
			query:        "NRN",
			radiusKm:     150,
			limit:        5,
			expectAnchor: "NRN",
			minCount:     1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nearby, anchor, found := FindNearbyAirports(tc.query, tc.radiusKm, tc.limit)
			if !found {
				t.Fatalf("expected to find airports for %q, but got found=false", tc.query)
			}
			if tc.expectAnchor != "" && anchor.IATA != tc.expectAnchor {
				t.Errorf("expected anchor %s, got %s", tc.expectAnchor, anchor.IATA)
			}
			if tc.expectCity != "" && anchor.City != tc.expectCity {
				t.Errorf("expected anchor city %s, got %s", tc.expectCity, anchor.City)
			}
			if len(nearby) < tc.minCount {
				t.Errorf("expected at least %d nearby airports, got %d", tc.minCount, len(nearby))
			}
		})
	}
}

func TestDistanceKm(t *testing.T) {
	// Distance between Brussels CRL (50.4592, 4.4538) and BRU (50.9014, 4.4844) is ~49 km
	crl := Airport{Latitude: 50.4592, Longitude: 4.4538}
	bru := Airport{Latitude: 50.9014, Longitude: 4.4844}

	dist := DistanceKm(crl.Latitude, crl.Longitude, bru.Latitude, bru.Longitude)
	if dist < 45 || dist > 55 {
		t.Errorf("unexpected distance between CRL and BRU: got %.2f km, expected ~49 km", dist)
	}
}
