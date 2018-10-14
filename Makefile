build:
	go build ./...
	mv ./server ./server.out

full: | gen_migrations test integration_test

test:
	go test ./... -v -short

integration_test: | vendor_deps ci

ci:
	./build/integration/int-test.sh


gen_migrations:
	go get -u github.com/jteeuwen/go-bindata/...
	go-bindata -prefix "migrations/" -nometadata -pkg main ./migrations
	mv bindata.go cmd/server/db_migrations_bindata.go
	go fmt cmd/server/db_migrations_bindata.go

vendor_deps:
	go mod vendor

.PHONY: integration_test vendor_deps test build ci
