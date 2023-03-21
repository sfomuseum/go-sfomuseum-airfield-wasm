package airports

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aaronland/go-roster"
	"github.com/sfomuseum/go-sfomuseum-airfield"	
)

type AirportsLookup interface {
	airfield.Lookup
	// FindAirport(context.Context, string) ([]*Airport, error)
}

func init() {
	ctx := context.Background()
	airfield.RegisterLookup(ctx, "airports", newLookup)
}

func newLookup(ctx context.Context, uri string) (airfield.Lookup, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	// Rewrite airports://sfomuseum/github as sfomuseum://github

	u.Scheme = u.Host
	u.Host = ""

	path := strings.TrimLeft(u.Path, "/")
	p := strings.Split(path, "/")

	if len(p) > 0 {
		u.Host = p[0]
		u.Path = strings.Join(p[1:], "/")
	}

	return NewAirportsLookup(ctx, u.String())
}

var airports_lookup_roster roster.Roster

type AirportsLookupInitializationFunc func(ctx context.Context, uri string) (AirportsLookup, error)

func RegisterAirportsLookup(ctx context.Context, scheme string, init_func AirportsLookupInitializationFunc) error {

	err := ensureAirportsLookupRoster()

	if err != nil {
		return err
	}

	return airports_lookup_roster.Register(ctx, scheme, init_func)
}

func ensureAirportsLookupRoster() error {

	if airports_lookup_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		airports_lookup_roster = r
	}

	return nil
}

func NewAirportsLookup(ctx context.Context, uri string) (AirportsLookup, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := airports_lookup_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(AirportsLookupInitializationFunc)
	return init_func(ctx, uri)
}
