package sfomuseum

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"github.com/whosonfirst/go-whosonfirst-uri"	
)

// CompileAircraftData will generate a list of `Aircraft` struct to be used as the source data for an `SFOMuseumLookup` instance.
// The list of aircraft are compiled by iterating over one or more source. `iterator_uri` is a valid `whosonfirst/go-whosonfirst-iterate` URI
// and `iterator_sources` are one more (iterator) URIs to process.
func CompileAircraftData(ctx context.Context, iterator_uri string, iterator_sources ...string) ([]*Aircraft, error) {

	lookup := make([]*Aircraft, 0)
	mu := new(sync.RWMutex)

	iter_cb := func(ctx context.Context, path string, fh io.ReadSeeker, args ...interface{}) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		_, uri_args, err := uri.ParseURI(path)

		if err != nil {
			return fmt.Errorf("Failed to parse %s, %w", path, err)
		}

		if uri_args.IsAlternate {
			return nil
		}

		body, err := io.ReadAll(fh)

		if err != nil {
			return fmt.Errorf("Failed to read '%s', %w", path, err)
		}

		wof_id, err := properties.Id(body)

		if err != nil {
			return fmt.Errorf("Failed to derive ID for %s, %w", path, err)
		}

		wof_name, err := properties.Name(body)

		if err != nil {
			return fmt.Errorf("Failed to derive name for %s, %w", path, err)
		}

		fl, err := properties.IsCurrent(body)

		if err != nil {
			return fmt.Errorf("Failed to determine is current for %s, %v", path, err)
		}

		sfom_id := int64(-1)

		sfom_rsp := gjson.GetBytes(body, "properties.sfomuseum:aircraft_id")

		if sfom_rsp.Exists() {
			sfom_id = sfom_rsp.Int()
		}

		a := &Aircraft{
			WhosOnFirstId: wof_id,
			SFOMuseumId:   sfom_id,
			Name:          wof_name,
			IsCurrent:     fl.Flag(),
		}

		concordances := properties.Concordances(body)

		if concordances != nil {

			code, ok := concordances["icao:designator"]

			if ok {
				a.ICAODesignator = fmt.Sprintf("%s", code)
			}

			id, ok := concordances["wd:id"]

			if ok {
				a.WikidataId = fmt.Sprintf("%s", id)
			}
		}

		mu.Lock()
		lookup = append(lookup, a)
		mu.Unlock()

		return nil
	}

	iter, err := iterator.NewIterator(ctx, iterator_uri, iter_cb)

	if err != nil {
		return nil, fmt.Errorf("Failed to create iterator, %w", err)
	}

	err = iter.IterateURIs(ctx, iterator_sources...)

	if err != nil {
		return nil, fmt.Errorf("Failed to iterate sources, %w", err)
	}

	return lookup, nil
}
