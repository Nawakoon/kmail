integration-test:
	@echo "Running integration tests"
	@go clean -testcache && go test ./tests/...

test:
	@echo "Running unit tests"
	@go clean -testcache && go test ./pkg/...

stress-test: # run make test 10 times, if any test fails, it will stop
	@echo "Running stress tests"
	@for i in {1..20}; do \
		make test || exit 1; \
	done

stress-integration-test: # run make integration-test 10 times, if any test fails, it will stop
	@echo "Running stress integration tests"
	@for i in {1..20}; do \
		make integration-test || exit 1; \
	done