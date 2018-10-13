test:
	go test ./... -v -short

integration_test: vendor_deps
	./build/integration/int-test.sh

vendor_deps:
	go mod vendor

.PHONY: integration_test vendor_deps test
