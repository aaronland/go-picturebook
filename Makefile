GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -tags libheif -ldflags="$(LDFLAGS)" -o bin/picturebook cmd/picturebook/main.go
