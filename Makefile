setup-lint:
	go get -u golang.org/x/lint/golint
	go get -u golang.org/x/tools/cmd/goimports

lint:
	go vet
	golint -set_exit_status ./...

fmt: lint
	goimports -w ./

test: fmt
	ENV=test go test -cover -race ./...

build:
	go build -ldflags="-s -w" -o dist/got cmd/got/main.go
