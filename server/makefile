integration-test: # integration tests
	@echo "Running integration tests"
	@go clean -testcache && go test ./tests/...

test: # unit tests
	@echo "Running unit tests"
	@go clean -testcache && go test ./pkg/...

stress-test: # run make test 10 times
	@echo "Running stress tests"
	@for i in {1..20}; do \
		make test || exit 1; \
	done

stress-integration-test: # run make integration-test 10 times
	@echo "Running stress integration tests"
	@for i in {1..20}; do \
		make integration-test || exit 1; \
	done

list: # list all commands
	@cat Makefile | grep -E '^[a-zA-Z0-9_-]+:.*'