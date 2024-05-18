TEST_PRIVATE_KEY = 1baa694c49154f63b1503c7138f184c80f221670f035403ff428a65183bab247

test:
	@go clean -testcache && go test ./pkg/...

integration-test:
	@go clean -testcache && go test ./tests/...

inbox:
	@go run cmd/main.go -inbox ./tests/util/query.test.json -user $(TEST_PRIVATE_KEY)

send-mail:
	@go run cmd/main.go -send ./tests/util/good.kmail.json -user $(TEST_PRIVATE_KEY)