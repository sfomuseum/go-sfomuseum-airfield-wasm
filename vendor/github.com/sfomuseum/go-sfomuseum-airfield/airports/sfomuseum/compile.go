package sfomuseum

import (
	"context"
	"fmt"
	"io"
	_ "log"
	"sync"

	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"github.com/whosonfirst/go-whosonfirst-uri"	
)

func CompileAirportsData(ctx context.Context, iterator_uri string, iterator_sources ...string) ([]*Airport, error) {

	lookup := make([]*Airport, 0)
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
			return fmt.Errorf("Failed to read %s, %w", path, err)
		}

		pt_rsp := gjson.GetBytes(body, "properties.sfomuseum:placetype")

		if pt_rsp.String() != "airport" {
			return nil
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

		sfom_rsp := gjson.GetBytes(body, "properties.sfomuseum:airport_id")

		if sfom_rsp.Exists() {
			sfom_id = sfom_rsp.Int()
		}

		a := &Airport{
			WhosOnFirstId: wof_id,
			SFOMuseumId:   sfom_id,
			Name:          wof_name,
			IsCurrent:     fl.Flag(),
		}

		concordances := properties.Concordances(body)

		if concordances != nil {

			iata_code, ok := concordances["iata:code"]

			if ok && iata_code != "" {
				a.IATACode = fmt.Sprintf("%s", iata_code)
			}

			icao_code, ok := concordances["icao:code"]

			if ok && icao_code != "" {
				a.ICAOCode = fmt.Sprintf("%s", icao_code)
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
