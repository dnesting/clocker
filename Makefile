all: test clocker

clocker: cmd/clocker/*.go *.go go.mod
	go build -o $@ cmd/clocker/*.go

test:
	go test ./...

clean:
	rm clocker

.PHONY: all test clean
