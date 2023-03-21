package main

import (
	_ "github.com/sfomuseum/go-sfomuseum-airfield/aircraft/sfomuseum"
	_ "github.com/sfomuseum/go-sfomuseum-airfield/airlines/sfomuseum"
	_ "github.com/sfomuseum/go-sfomuseum-airfield/airports/sfomuseum"
)
	
import (
	"log"
	"syscall/js"

	"github.com/sfomuseum/go-sfomuseum-airfield-wasm"
	"github.com/sfomuseum/go-sfomuseum-airfield"		
)

func main() {

	airport_lookup, err := airfield.NewLookup(ctx, "airports://sfomuseum")

	if err != nil {
		log.Fatalf("Failed to create airports lookup, %v", err)
	}
	
	lookup_airport_func := wasm.LookupAirportFunc(airport_lookup)
	defer lookup_airport_func.Release()

	js.Global().Set("sfomuseum_lookup_airport", lookup_airport_func)

	c := make(chan struct{}, 0)

	log.Println("SFO Museum airfield lookup WASM binary initialized")
	<-c
}
