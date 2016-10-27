build:
	go build -v

fmt:
	gofmt -w -s *.go

test:
	go test -v -cover -race

