.PHONY: test
test:
	@go test -v -cover .
	@go vet .
	@golint .

.DEFAULT_GOAL := test
