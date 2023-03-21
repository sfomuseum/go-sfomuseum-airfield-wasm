GOMOD=vendor

cli:
	@make cli-lookup
	# @make cli-create
	@make cli-stats

cli-create:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/create-airline cmd/create-airline/main.go

cli-lookup:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/lookup cmd/lookup/main.go

cli-stats:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/tailnumbers cmd/tailnumbers/main.go

compile:
	@make compile-airlines
	@make compile-airports
	@make compile-aircraft
	@make cli-lookup

compile-airlines:
	go run -mod $(GOMOD) cmd/compile-flysfo-airlines-data/main.go -iterator-uri 'git:///tmp?exclude=properties.edtf:deprecated=.*' https://github.com/sfomuseum-data/sfomuseum-data-enterprise.git
	go run -mod $(GOMOD) cmd/compile-sfomuseum-airlines-data/main.go  -iterator-uri 'git:///tmp?exclude=properties.edtf:deprecated=.*' https://github.com/sfomuseum-data/sfomuseum-data-enterprise.git

compile-airports:
	go run -mod $(GOMOD) cmd/compile-sfomuseum-airports-data/main.go  -iterator-uri 'git:///tmp?include=properties.sfomuseum:placetype=airport&exclude=properties.edtf:deprecated=.*' https://github.com/sfomuseum-data/sfomuseum-data-whosonfirst.git

compile-aircraft:
	go run -mod $(GOMOD) cmd/compile-sfomuseum-aircraft-data/main.go -iterator-uri 'git:///tmp?exclude=properties.edtf:deprecated=.*' https://github.com/sfomuseum-data/sfomuseum-data-aircraft.git

local-scan:
	/usr/local/bin/sonar-scanner/bin/sonar-scanner -Dsonar.projectKey=go-sfomuseum-airfield -Dsonar.sources=. -Dsonar.host.url=http://localhost:9000 -Dsonar.login=$(TOKEN)
