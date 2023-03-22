# go-sfomuseum-airfield-wasm

Go package for compiling methods from the [go-sfomuseum-airfield](https://github.com/sfomuseum/go-sfomuseum-airfield) package to a JavaScript-compatible WebAssembly (wasm) binary. It also provides a net/http middleware packages for appending the necessary static assets and HTML resources to use the wasm binary in web applications.

## Documentation

Documentation is incomplete at this time.

## Example

There is a full working example of this application in the `cmd/example` folder. To run this application type the following command:

```
$> make example
go run -mod vendor cmd/example/main.go -port 8080
2023/03/21 18:29:08 Listening for requests on localhost:8080
```

Then open `http://localhost:8080` in a  web browser. You should see something like this:

![](docs/images/sfomuseum-airfield-wasm.png)

## See also

* https://github.com/sfomuseum/go-sfomuseum-airfield