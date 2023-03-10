#=====================================#
# HELPERS #
#=====================================#

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


#=====================================#
# DEVELOPMENT #
#=====================================#

## run: run the ./main.go application
.PHONY: run
run:
	go run ./main.go

## build: build the ./main.go application
.PHONY: build
build:
	go build -o WeightChangerFAGRIP/WeightChangerFAGRIP.exe

#=====================================#
# QUALITY CONTROL #
#=====================================#

## audit: tidy dependencies and format, vet and test all code

## go fmt ./... : command to format all .go files in the project directory, according to the Go standard.
## go vet ./... : runs a variety of analyzers which carry out static analysis of your code and warn you
## go test -race -vet=off ./... : command to run all tests in the project directory
## staticcheck tool : to carry out some additional static analysis checks.
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## go mod tidy : prune any unused dependencies from the go.mod and go.sum files, and add any missing dependencies
## go mod verify : check that the dependencies on your computer (located in your module cache located at $GOPATH/pkg/mod)
## haven’t been changed since they were downloaded and that they match the cryptographic hashes in your go.sum file
## go mod vendor: copy the necessary source code from your module cache into a new vendor directory in your project root
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor