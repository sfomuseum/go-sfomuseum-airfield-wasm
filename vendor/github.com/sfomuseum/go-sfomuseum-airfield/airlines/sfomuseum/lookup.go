package sfomuseum

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/sfomuseum/go-sfomuseum-airfield/airlines"
	"github.com/sfomuseum/go-sfomuseum-airfield/data"	
)

var lookup_table *sync.Map
var lookup_idx int64

var lookup_init sync.Once
var lookup_init_err error

type SFOMuseumLookupFunc func(context.Context)

type SFOMuseumLookup struct {
	airlines.AirlinesLookup
}

func init() {
	ctx := context.Background()
	airlines.RegisterAirlinesLookup(ctx, "sfomuseum", NewSFOMuseumLookup)

	lookup_idx = int64(0)
}

// NewLookup will return an `airlines.AirlinesLookup` instance. By default the lookup table is derived from precompiled (embedded) data in `data/sfomuseum.json`
// by passing in `sfomuseum://` as the URI. It is also possible to create a new lookup table with the following URI options:
// 	`sfomuseum://github`
// This will cause the lookup table to be derived from the data stored at https://raw.githubusercontent.com/sfomuseum/go-sfomuseum-airlines/main/data/sfomuseum.json. This might be desirable if there have been updates to the underlying data that are not reflected in the locally installed package's pre-compiled data.
//	`sfomuseum://iterator?uri={URI}&source={SOURCE}`
// This will cause the lookup table to be derived, at runtime, from data emitted by a `whosonfirst/go-whosonfirst-iterate` instance. `{URI}` should be a valid `whosonfirst/go-whosonfirst-iterate/iterator` URI and `{SOURCE}` is one or more URIs for the iterator to process.
func NewSFOMuseumLookup(ctx context.Context, uri string) (airlines.AirlinesLookup, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	// Account for both:
	// airfield.NewLookup(ctx, "airlines://sfomuseum/github")
	// airlines.NewAirLinesLookup(ctx, "sfomuseum://github")

	var source string

	switch u.Host {
	case "sfomuseum":
		source = u.Path
	default:
		source = u.Host
	}

	switch source {
	case "iterator":

		q := u.Query()

		iterator_uri := q.Get("uri")
		iterator_sources := q["source"]

		return NewSFOMuseumLookupFromIterator(ctx, iterator_uri, iterator_sources...)

	case "github":

		rsp, err := http.Get(DATA_GITHUB)

		if err != nil {
			return nil, fmt.Errorf("Failed to load remote data from Github, %w", err)
		}

		lookup_func := NewSFOMuseumLookupFuncWithReader(ctx, rsp.Body)
		return NewSFOMuseumLookupWithLookupFunc(ctx, lookup_func)

	default:

		fs := data.FS
		fh, err := fs.Open(DATA_JSON)

		if err != nil {
			return nil, fmt.Errorf("Failed to load local precompiled data, %w", err)
		}

		lookup_func := NewSFOMuseumLookupFuncWithReader(ctx, fh)
		return NewSFOMuseumLookupWithLookupFunc(ctx, lookup_func)
	}
}

// NewLookup will return an `SFOMuseumLookupFunc` function instance that, when invoked, will populate an `airlines.AirlinesLookup` instance with data stored in `r`.
// `r` will be closed when the `SFOMuseumLookupFunc` function instance is invoked.
// It is assumed that the data in `r` will be formatted in the same way as the procompiled (embedded) data stored in `data/sfomuseum.json`.
func NewSFOMuseumLookupFuncWithReader(ctx context.Context, r io.ReadCloser) SFOMuseumLookupFunc {

	defer r.Close()

	var airlines_list []*Airline

	dec := json.NewDecoder(r)
	err := dec.Decode(&airlines_list)

	if err != nil {

		lookup_func := func(ctx context.Context) {
			lookup_init_err = err
		}

		return lookup_func
	}

	return NewSFOMuseumLookupFuncWithAirlines(ctx, airlines_list)
}

// NewLookup will return an `SFOMuseumLookupFunc` function instance that, when invoked, will populate an `airlines.AirlinesLookup` instance with data stored in `airlines_list`.
func NewSFOMuseumLookupFuncWithAirlines(ctx context.Context, airlines_list []*Airline) SFOMuseumLookupFunc {

	lookup_func := func(ctx context.Context) {

		table := new(sync.Map)

		for _, data := range airlines_list {

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			appendData(ctx, table, data)
		}

		lookup_table = table
	}

	return lookup_func
}

// NewLookupWithLookupFunc will return an `airlines.AirlinesLookup` instance derived by data compiled using `lookup_func`.
func NewSFOMuseumLookupWithLookupFunc(ctx context.Context, lookup_func SFOMuseumLookupFunc) (airlines.AirlinesLookup, error) {

	fn := func() {
		lookup_func(ctx)
	}

	lookup_init.Do(fn)

	if lookup_init_err != nil {
		return nil, lookup_init_err
	}

	l := SFOMuseumLookup{}
	return &l, nil
}

func NewSFOMuseumLookupFromIterator(ctx context.Context, iterator_uri string, iterator_sources ...string) (airlines.AirlinesLookup, error) {

	airline_data, err := CompileAirlinesData(ctx, iterator_uri, iterator_sources...)

	if err != nil {
		return nil, fmt.Errorf("Failed to compile airline data, %w", err)
	}

	lookup_func := NewSFOMuseumLookupFuncWithAirlines(ctx, airline_data)
	return NewSFOMuseumLookupWithLookupFunc(ctx, lookup_func)
}

func (l *SFOMuseumLookup) Find(ctx context.Context, code string) ([]interface{}, error) {

	pointers, ok := lookup_table.Load(code)

	if !ok {
		return nil, airlines.NotFound{code}
	}

	airline := make([]interface{}, 0)

	for _, p := range pointers.([]string) {

		if !strings.HasPrefix(p, "pointer:") {
			return nil, fmt.Errorf("Invalid pointer, %s", p)
		}

		row, ok := lookup_table.Load(p)

		if !ok {
			return nil, fmt.Errorf("Invalid pointer, %s", p)
		}

		airline = append(airline, row.(*Airline))
	}

	return airline, nil
}

func (l *SFOMuseumLookup) Append(ctx context.Context, data interface{}) error {
	return appendData(ctx, lookup_table, data.(*Airline))
}

func appendData(ctx context.Context, table *sync.Map, data *Airline) error {

	idx := atomic.AddInt64(&lookup_idx, 1)

	pointer := fmt.Sprintf("pointer:%d", idx)
	table.Store(pointer, data)

	str_wofid := strconv.FormatInt(data.WhosOnFirstId, 10)
	str_sfomid := strconv.FormatInt(data.SFOMuseumId, 10)

	possible_codes := []string{
		str_wofid,
		str_sfomid,
		fmt.Sprintf("wof:id=%s", str_wofid),
		fmt.Sprintf("sfomuseum:airline_id=%s", str_sfomid),
	}

	if data.IATACode != "" {
		possible_codes = append(possible_codes, data.IATACode)
		possible_codes = append(possible_codes, fmt.Sprintf("iata:code=%s", data.IATACode))
	}

	if data.ICAOCode != "" {
		possible_codes = append(possible_codes, data.ICAOCode)
		possible_codes = append(possible_codes, fmt.Sprintf("icao:code=%s", data.ICAOCode))
	}

	if data.ICAOCallsign != "" {
		possible_codes = append(possible_codes, data.ICAOCallsign)
		possible_codes = append(possible_codes, fmt.Sprintf("icao:callsign=%s", data.ICAOCallsign))
	}

	if data.WikidataId != "" {
		possible_codes = append(possible_codes, data.WikidataId)
		possible_codes = append(possible_codes, fmt.Sprintf("wikidata:id=%s", data.WikidataId))
	}

	for _, code := range possible_codes {

		if code == "" {
			continue
		}

		pointers := make([]string, 0)
		has_pointer := false

		others, ok := table.Load(code)

		if ok {

			pointers = others.([]string)
		}

		for _, dupe := range pointers {

			if dupe == pointer {
				has_pointer = true
				break
			}
		}

		if has_pointer {
			continue
		}

		pointers = append(pointers, pointer)
		table.Store(code, pointers)
	}

	return nil
}
