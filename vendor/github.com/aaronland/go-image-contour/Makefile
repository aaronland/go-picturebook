GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/contour cmd/contour/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/contour-svg cmd/contour-svg/main.go
