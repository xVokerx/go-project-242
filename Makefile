build:
	go build -o bin/hexlet-path-size.exe ./cmd/hexlet-path-size 

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

tests:
	go test -v -timeout 30s
