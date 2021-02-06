install_deps:
	go mod download

build:
	go build
# Standard go test
test:
	go test ./...

# Make sure no unnecessary dependecies are present
go-mod-tidy:
	go mod tidy -v
	git diff-index --quiet HEAD

# Run all tests & linters in CI
ci: test
