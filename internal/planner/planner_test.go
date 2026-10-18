package planner

import (
	"testing"
)

func TestTripRequestDefaults(t *testing.T) {
	req := TripRequest{
		Origin:           "Brussels",
		Destination:      "Malaga",
		OutboundDateFrom: "2026-11-01",
	}

	p := NewPlanner(nil)
	if p == nil {
		t.Fatal("expected NewPlanner to return non-nil instance")
	}

	if req.Origin != "Brussels" || req.Destination != "Malaga" {
		t.Errorf("unexpected trip request fields: %+v", req)
	}
}
