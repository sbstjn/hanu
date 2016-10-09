test:
	go test -cover ./...

race:
	go test -v -race ./...

cover:
	@./script/coverage
