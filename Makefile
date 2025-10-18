build:
	go build ./lga/...

install:
	go install ./lga

new: build install

test:
	go test ./...
