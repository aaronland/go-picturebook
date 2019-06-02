fmt:
	go fmt cmd/picturebook/*.go

tools:
	go build -o bin/picturebook cmd/picturebook/main.go
