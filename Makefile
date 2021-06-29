install:
	go install
test:
	go test -cover $(shell go list ./...)