.PHONY: all
all: test build install

.PHONY: build
build: 
	go build
	go vet
	golint .

.PHONY: test
test:
	go test
	GOOS=darwin go build
	GOOS=linux go build

.PHONY: install
install:
	go install

.PHONY: bench
bench:
	go test -bench .
