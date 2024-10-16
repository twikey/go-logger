.PHONY: bench test cov

bench:
	go test -bench=. -benchmem

test:
	go test -v -cover ./...

cov:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@rm coverage.out