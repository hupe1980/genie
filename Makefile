PROJECTNAME=$(shell basename "$(PWD)")

# Go related variables.
# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: setup
## setup: Setup installes dependencies
setup:
	@go mod tidy

.PHONY: lint
## test: Runs the linter
lint:
	@golangci-lint run --color=always --sort-results ./...

.PHONY: run
## run: Runs awsrecon
run:
	@go run -race main.go -h

.PHONY: test
## test: Runs go test with default values
test: 
	@go test -v -race -count=1 -coverprofile=coverage.out ./...

.PHONY: build
## build: Build from source
build:
	@go build -o genie .

.PHOMY: docker-build
## docker-build: Build a docker image
docker-build:
	docker build -t genie .

.PHONY: help
## help: Prints this help message
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo