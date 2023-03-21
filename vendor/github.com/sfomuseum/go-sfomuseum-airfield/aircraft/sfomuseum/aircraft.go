package sfomuseum

import (
	"context"
	"fmt"
	
	"github.com/sfomuseum/go-sfomuseum-airfield"
	"github.com/sfomuseum/go-sfomuseum-airfield/aircraft"
)

type Aircraft struct {
	WhosOnFirstId  int64  `json:"wof:id"`
	Name           string `json:"wof:name"`
	SFOMuseumId    int64  `json:"sfomuseum:aircraft_id"`
	ICAODesignator string `json:"icao:designator,omitempty"`
	WikidataId     string `json:"wd:id,omitempty"`
	IsCurrent      int64  `json:"mz:is_current"`
}

func (a *Aircraft) String() string {
	return fmt.Sprintf("%s \"%s\" %d (%d) (%s) Is current: %d", a.ICAODesignator, a.Name, a.WhosOnFirstId, a.SFOMuseumId, a.WikidataId, a.IsCurrent)
}

// Return the current Aircraft matching 'code'. Multiple matches throw an error.
func FindCurrentAircraft(ctx context.Context, code string) (*Aircraft, error) {

	lookup, err := NewSFOMuseumLookup(ctx, "")

	if err != nil {
		return nil, fmt.Errorf("Failed to create new lookup, %w", err)
	}

	return FindCurrentAircraftWithLookup(ctx, lookup, code)
}

// Return the current Aircraft matching 'code' with a custom airfield.Lookup instance. Multiple matches throw an error.
func FindCurrentAircraftWithLookup(ctx context.Context, lookup airfield.Lookup, code string) (*Aircraft, error) {

	current, err := FindAircraftCurrentWithLookup(ctx, lookup, code)

	if err != nil {
		return nil, err
	}

	switch len(current) {
	case 0:
		return nil, aircraft.NotFound{code}
	case 1:
		return current[0], nil
	default:
		return nil, aircraft.MultipleCandidates{code}
	}

}

// Returns all Aircraft instances matching 'code' that are marked as current.
func FindAircraftCurrent(ctx context.Context, code string) ([]*Aircraft, error) {

	lookup, err := NewSFOMuseumLookup(ctx, "")

	if err != nil {
		return nil, fmt.Errorf("Failed to create new lookup, %w", err)
	}

	return FindAircraftCurrentWithLookup(ctx, lookup, code)
}

// Returns all Aircraft instances matching 'code' that are marked as current with a custom airfield.Lookup instance.
func FindAircraftCurrentWithLookup(ctx context.Context, lookup airfield.Lookup, code string) ([]*Aircraft, error) {

	rsp, err := lookup.Find(ctx, code)

	if err != nil {
		return nil, aircraft.NotFound{code}
	}

	current := make([]*Aircraft, 0)

	for _, r := range rsp {

		g := r.(*Aircraft)

		// if g.IsCurrent == 0 {
		if g.IsCurrent != 1 {
			continue
		}

		current = append(current, g)
	}

	return current, nil
}
