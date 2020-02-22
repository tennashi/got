lint:
	go vet ./...
	golint -set_exit_status ./...

fmt: lint
	goimports -w ./

test:
	ENV=test go test -cover -race ./...

build:
	go build -ldflags="-s -w" -o dist/got cmd/got/main.go
