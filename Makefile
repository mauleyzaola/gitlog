LINT_VERSION=v1.37.1

install:
	go install

lint:
	golangci-lint run --config golangci.yml --timeout 5m

lint-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(LINT_VERSION)

test:
	go test -cover $(shell go list ./...)