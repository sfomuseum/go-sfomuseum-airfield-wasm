package wasm

import (
	"context"
	"encoding/json"
	"syscall/js"

	"github.com/sfomuseum/go-sfomuseum-airfield/aircraft"	
)

func LookupAircraftFunc(lookup aircraft.AircraftLookup) js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		code := args[0].String()
		
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			go func() {

				ctx := context.Background()
				rsp, err := lookup.Find(ctx, code)

				if err != nil {
					reject.Invoke(err.Error())
					return
				}

				enc_rsp, err := json.Marshal(rsp)

				if err != nil {
					reject.Invoke(err.Error())
					return
				}

				resolve.Invoke(string(enc_rsp))
			}()

			return nil
		})
		
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
