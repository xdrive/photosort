.PHONY: deps test test-cover

deps:
		@echo "Installing dependencies"
		go get -u github.com/golang/dep/cmd/dep
		dep ensure

test:
		@echo "Running tests"
		go test ./...

test-cover:
		@echo "Running tests with coverage"
		go test ./... -cover
	