fmt:
	go fmt cmd/picturebook/*.go

tools:
	go build -mod vendor -o bin/picturebook cmd/picturebook/main.go
