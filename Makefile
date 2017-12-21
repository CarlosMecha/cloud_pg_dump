
build:
	@govendor build -o bin/cloud_pg_dump

test:
	@govendor test -race -cover +local

.PHONY: test build