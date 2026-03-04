GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

vuln:
	govulncheck -show verbose ./...

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/picturebook cmd/picturebook/main.go
