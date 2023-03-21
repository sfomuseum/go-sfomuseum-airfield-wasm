package sfomuseum

import (
	"context"
	"fmt"
	
	"github.com/sfomuseum/go-sfomuseum-airfield"
	"github.com/sfomuseum/go-sfomuseum-airfield/airlines"
)

type Airline struct {
	WhosOnFirstId int64  `json:"wof:id"`
	Name          string `json:"wof:name"`
	Role          string `json:"sfomuseum:airline_role"`
	SFOMuseumId   int64  `json:"sfomuseum:airline_id"`
	IATACode      string `json:"iata:code,omitempty"`
	ICAOCode      string `json:"icao:code,omitempty"`
	ICAOCallsign  string `json:"icao:callsign,omitempty"`
	WikidataId    string `json:"wd:id,omitempty"`
	IsCurrent     int64  `json:"mz:is_current"`
}

func (a *Airline) String() string {
	return fmt.Sprintf("%s %s %s \"%s\" %d (%d) (%s) Is current: %d", a.IATACode, a.ICAOCode, a.ICAOCallsign, a.Name, a.WhosOnFirstId, a.SFOMuseumId, a.WikidataId, a.IsCurrent)
}

// Return the current Airline matching 'code'. Multiple matches throw an error.
func FindCurrentAirline(ctx context.Context, code string, roles ...string) (*Airline, error) {

	lookup, err := NewSFOMuseumLookup(ctx, "")

	if err != nil {
		return nil, fmt.Errorf("Failed to create new lookup, %w", err)
	}

	return FindCurrentAirlineWithLookup(ctx, lookup, code, roles...)
}

// Return the current Airline matching 'code' with a custom airfield.Lookup instance. Multiple matches throw an error.
func FindCurrentAirlineWithLookup(ctx context.Context, lookup airfield.Lookup, code string, roles ...string) (*Airline, error) {

	current, err := FindAirlinesCurrentWithLookup(ctx, lookup, code)

	if err != nil {
		return nil, err
	}

	if len(roles) > 0 {

		candidates := make([]*Airline, 0)

		for _, a := range current {

			ok := false

			for _, r := range roles {

				if a.Role == r {
					ok = true
					break
				}
			}

			if ok {
				candidates = append(candidates, a)
			}
		}

		current = candidates
	}

	switch len(current) {
	case 0:
		return nil, airlines.NotFound{code}
	case 1:
		return current[0], nil
	default:
		return nil, airlines.MultipleCandidates{code}
	}

}

// Returns all Airline instances matching 'code' that are marked as current.
func FindAirlinesCurrent(ctx context.Context, code string) ([]*Airline, error) {

	lookup, err := NewSFOMuseumLookup(ctx, "")

	if err != nil {
		return nil, fmt.Errorf("Failed to create new lookup, %w", err)
	}

	return FindAirlinesCurrentWithLookup(ctx, lookup, code)
}

// Returns all Airline instances matching 'code' that are marked as current with a custom airfield.Lookup instance.
func FindAirlinesCurrentWithLookup(ctx context.Context, lookup airfield.Lookup, code string) ([]*Airline, error) {

	rsp, err := lookup.Find(ctx, code)

	if err != nil {
		return nil, airlines.NotFound{code}
	}

	current := make([]*Airline, 0)

	for _, r := range rsp {

		g := r.(*Airline)

		// if g.IsCurrent == 0 {
		if g.IsCurrent != 1 {
			continue
		}

		current = append(current, g)
	}

	return current, nil
}
