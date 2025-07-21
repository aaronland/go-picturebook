GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/contour cmd/contour/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/contour-svg cmd/contour-svg/main.go

wasmjs:
	GOOS=js GOARCH=wasm \
		go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags wasmjs \
		-o www/wasm/contour.wasm \
		cmd/contour-wasm-js/main.go
