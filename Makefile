GOROOT=$(shell go env GOROOT)
GOMOD=vendor

rebuild:
	@make wasm

wasm:
	GOOS=js GOARCH=wasm go build -mod $(GOMOD) -ldflags="-s -w" -o static/wasm/sfomuseum_airfield.wasm cmd/lookup/main.go

example:
	go run -mod $(GOMOD) cmd/example/main.go -port 8080
