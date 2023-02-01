PROJECT_NAME = "vcache"

.PHONY: all dep test lint 

all: init lint test race 

init:
  go mod tidy
  go mod vendor

pre-push: lint test race 

lint: 
  golangci-lint run -v ./...

test: 
  go test -v ./...

race: dep ## Run data race detector
  go test -race -v ./...

dep: ## Get the dependencies
  go get -v -d ./...
