## test: runs all tests
test:
	@go test -v ./...

## cover: opens coverage in browser
cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

## coverage: displays test coverage
coverage:
	@go test -cover ./...

## build_cli: builds the command line tool celeritas and copies it to myapp
# compile everything in ./cmd/cli (here)
# and then install the binaries in ../myapp/(goframework) will be the name of the binaries
build_cli:
	@go build -o ../myapp/goframework ./cmd/cli

## build: builds the command line tool dist directory
build:
	@go build -o ./dist/goframework ./cmd/cli
