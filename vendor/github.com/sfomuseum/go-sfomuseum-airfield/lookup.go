package airfield

import (
	"context"
	"github.com/aaronland/go-roster"
	"net/url"
)

type Lookup interface {
	Find(context.Context, string) ([]interface{}, error)
	Append(context.Context, interface{}) error
}

var lookup_roster roster.Roster

type LookupInitializationFunc func(ctx context.Context, uri string) (Lookup, error)

func RegisterLookup(ctx context.Context, scheme string, init_func LookupInitializationFunc) error {

	err := ensureLookupRoster()

	if err != nil {
		return err
	}

	return lookup_roster.Register(ctx, scheme, init_func)
}

func ensureLookupRoster() error {

	if lookup_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		lookup_roster = r
	}

	return nil
}

func NewLookup(ctx context.Context, uri string) (Lookup, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := lookup_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(LookupInitializationFunc)
	return init_func(ctx, uri)
}
