.PHONY: all
all: test build

.PHONY: build
build: 
	go build
	go vet
	golint .

.PHONY: test
test:
	go test

.PHONY: install
install:
	go install
