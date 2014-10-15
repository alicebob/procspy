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
	GOOS=darwin go build
	GOOS=linux go build

.PHONY: install
install:
	go install
