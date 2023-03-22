package main

import (
	_ "github.com/sfomuseum/go-sfomuseum-airfield/aircraft/sfomuseum"
	_ "github.com/sfomuseum/go-sfomuseum-airfield/airlines/sfomuseum"
	_ "github.com/sfomuseum/go-sfomuseum-airfield/airports/sfomuseum"
)
	
import (
	"context"
	"log"
	"syscall/js"

	"github.com/sfomuseum/go-sfomuseum-airfield-wasm"
	"github.com/sfomuseum/go-sfomuseum-airfield"		
)

func main() {

	ctx := context.Background()
	
	airport_lookup, err := airfield.NewLookup(ctx, "airports://sfomuseum")

	if err != nil {
		log.Fatalf("Failed to create airports lookup, %v", err)
	}

	airline_lookup, err := airfield.NewLookup(ctx, "airlines://sfomuseum")

	if err != nil {
		log.Fatalf("Failed to create airlines lookup, %v", err)
	}

	aircraft_lookup, err := airfield.NewLookup(ctx, "aircraft://sfomuseum")

	if err != nil {
		log.Fatalf("Failed to create aircrafts lookup, %v", err)
	}
	
	lookup_airport_func := wasm.LookupAirportFunc(airport_lookup)
	defer lookup_airport_func.Release()

	lookup_airline_func := wasm.LookupAirlineFunc(airline_lookup)
	defer lookup_airline_func.Release()

	lookup_aircraft_func := wasm.LookupAircraftFunc(aircraft_lookup)
	defer lookup_aircraft_func.Release()
	
	js.Global().Set("sfomuseum_lookup_airport", lookup_airport_func)
	js.Global().Set("sfomuseum_lookup_airline", lookup_airline_func)
	js.Global().Set("sfomuseum_lookup_aircraft", lookup_aircraft_func)	

	c := make(chan struct{}, 0)

	log.Println("SFO Museum airfield lookup WASM binary initialized")
	<-c
}
