DOCKER_IMAGE=jenkins-tests
LINT_VERSION=v1.43.0

.PHONY:install
install:
	go install

.PHONY:lint
lint:
	golangci-lint run --config golangci.yml --timeout 5m

.PHONY:lint-install
lint-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(LINT_VERSION)

.PHONY: test
test:
	go test -cover $(shell go list ./... | grep -v /integration)