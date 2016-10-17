test:
		go test -cover -race ./...

bench:
		go test -bench=. ./...

race:
		go test -v -race ./...

cover:
		./script/coverage

coveralls:
		./script/coverage --coveralls
