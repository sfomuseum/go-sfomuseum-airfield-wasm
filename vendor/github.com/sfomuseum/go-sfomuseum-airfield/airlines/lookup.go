package airlines

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aaronland/go-roster"
	"github.com/sfomuseum/go-sfomuseum-airfield"	
)

type AirlinesLookup interface {
	airfield.Lookup
	// FindAirline(context.Context, string) ([]*Airline, error)
}

func init() {
	ctx := context.Background()
	airfield.RegisterLookup(ctx, "airlines", newLookup)
}

func newLookup(ctx context.Context, uri string) (airfield.Lookup, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	// Rewrite airlines://sfomuseum/github as sfomuseum://github

	u.Scheme = u.Host
	u.Host = ""

	path := strings.TrimLeft(u.Path, "/")
	p := strings.Split(path, "/")

	if len(p) > 0 {
		u.Host = p[0]
		u.Path = strings.Join(p[1:], "/")
	}

	return NewAirlinesLookup(ctx, u.String())
}

var airlines_lookup_roster roster.Roster

type AirlinesLookupInitializationFunc func(ctx context.Context, uri string) (AirlinesLookup, error)

func RegisterAirlinesLookup(ctx context.Context, scheme string, init_func AirlinesLookupInitializationFunc) error {

	err := ensureAirlinesLookupRoster()

	if err != nil {
		return err
	}

	return airlines_lookup_roster.Register(ctx, scheme, init_func)
}

func ensureAirlinesLookupRoster() error {

	if airlines_lookup_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		airlines_lookup_roster = r
	}

	return nil
}

func NewAirlinesLookup(ctx context.Context, uri string) (AirlinesLookup, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := airlines_lookup_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(AirlinesLookupInitializationFunc)
	return init_func(ctx, uri)
}
