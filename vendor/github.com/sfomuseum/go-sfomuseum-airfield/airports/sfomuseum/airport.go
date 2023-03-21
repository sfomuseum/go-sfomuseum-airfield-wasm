package sfomuseum

import (
	"context"
	"fmt"
	
	"github.com/sfomuseum/go-sfomuseum-airfield"
	"github.com/sfomuseum/go-sfomuseum-airfield/airports"
)

type Airport struct {
	WhosOnFirstId int64  `json:"wof:id"`
	Name          string `json:"wof:name"`
	SFOMuseumId   int64  `json:"sfomuseum:airport_id"`
	IATACode      string `json:"iata:code"`
	ICAOCode      string `json:"icao:code"`
	WikidataId    string `json:"wd:id,omitempty"`
	IsCurrent     int64  `json:"mz:is_current"`
}

func (a *Airport) String() string {
	return fmt.Sprintf("%s %s \"%s\" %d (%d) (%s) Is Current: %d", a.IATACode, a.ICAOCode, a.Name, a.WhosOnFirstId, a.SFOMuseumId, a.WikidataId, a.IsCurrent)
}

// Return the current Airport matching 'code'. Multiple matches throw an error.
func FindCurrentAirport(ctx context.Context, code string) (*Airport, error) {

	lookup, err := NewSFOMuseumLookup(ctx, "")

	if err != nil {
		return nil, fmt.Errorf("Failed to create new lookup, %w", err)
	}

	return FindCurrentAirportWithLookup(ctx, lookup, code)
}

// Return the current Airport matching 'code' with a custom airfield.Lookup instance. Multiple matches throw an error.
func FindCurrentAirportWithLookup(ctx context.Context, lookup airfield.Lookup, code string) (*Airport, error) {

	current, err := FindAirportsCurrentWithLookup(ctx, lookup, code)

	if err != nil {
		return nil, err
	}

	switch len(current) {
	case 0:
		return nil, airports.NotFound{code}
	case 1:
		return current[0], nil
	default:
		return nil, airports.MultipleCandidates{code}
	}

}

// Returns all Airport instances matching 'code' that are marked as current.
func FindAirportsCurrent(ctx context.Context, code string) ([]*Airport, error) {

	lookup, err := NewSFOMuseumLookup(ctx, "")

	if err != nil {
		return nil, fmt.Errorf("Failed to create new lookup, %w", err)
	}

	return FindAirportsCurrentWithLookup(ctx, lookup, code)
}

// Returns all Airport instances matching 'code' that are marked as current with a custom airfield.Lookup instance.
func FindAirportsCurrentWithLookup(ctx context.Context, lookup airfield.Lookup, code string) ([]*Airport, error) {

	rsp, err := lookup.Find(ctx, code)

	if err != nil {
		return nil, fmt.Errorf("Failed to find (sfomuseum) airport '%s', %w", code, err)
	}

	current := make([]*Airport, 0)

	for _, r := range rsp {

		g := r.(*Airport)

		// if g.IsCurrent == 0 {
		if g.IsCurrent != 1 {
			continue
		}

		current = append(current, g)
	}

	return current, nil
}
