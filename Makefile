setup:
	go mod download
	go get -u golang.org/x/lint/golint
	go get -u golang.org/x/tools/cmd/goimports

lint:
	go tool vet ./
	golint -set_exit_status ./...

fmt: lint
	goimports -w ./

test: fmt
	ENV=test go test -cover -race ./...

build:
	go build -ldflags="-s -w" -o dist/got ./main.go
