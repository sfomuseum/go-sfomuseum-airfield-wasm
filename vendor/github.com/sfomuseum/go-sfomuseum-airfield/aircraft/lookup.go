package aircraft

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aaronland/go-roster"
	"github.com/sfomuseum/go-sfomuseum-airfield"	
)

type AircraftLookup interface {
	airfield.Lookup
}

func init() {
	ctx := context.Background()
	airfield.RegisterLookup(ctx, "aircraft", newAircraftLookup)
}

// Private method for airfield.RegisterLookup
func newAircraftLookup(ctx context.Context, uri string) (airfield.Lookup, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	// Rewrite aircraft://sfomuseum/github as sfomuseum://github

	u.Scheme = u.Host
	u.Host = ""

	path := strings.TrimLeft(u.Path, "/")
	p := strings.Split(path, "/")

	if len(p) > 0 {
		u.Host = p[0]
		u.Path = strings.Join(p[1:], "/")
	}

	return NewAircraftLookup(ctx, u.String())
}

var aircraft_lookup_roster roster.Roster

type AircraftLookupInitializationFunc func(ctx context.Context, uri string) (AircraftLookup, error)

func RegisterAircraftLookup(ctx context.Context, scheme string, init_func AircraftLookupInitializationFunc) error {

	err := ensureAircraftLookupRoster()

	if err != nil {
		return err
	}

	return aircraft_lookup_roster.Register(ctx, scheme, init_func)
}

func ensureAircraftLookupRoster() error {

	if aircraft_lookup_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		aircraft_lookup_roster = r
	}

	return nil
}

func NewAircraftLookup(ctx context.Context, uri string) (AircraftLookup, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := aircraft_lookup_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(AircraftLookupInitializationFunc)
	return init_func(ctx, uri)
}
