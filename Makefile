.PHONY: test
test:
	@mkdir -p `pwd`/docs/coverage
	@go test -coverprofile `pwd`/docs/coverage/coverprofile.out ./... | { grep -v 'no test files'; true; }
	@echo "\n== RESUME ==="
	@go tool cover -func `pwd`/docs/coverage/coverprofile.out
	@go tool cover -html=`pwd`/docs/coverage/coverprofile.out -o docs/coverage/index.html