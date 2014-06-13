.PHONY: all
all: main

.PHONY: c
c:
	cd src && yes n | ./Configure
	cp src/*.[ch] .
	rm main.c

.PHONY: main
main:
	go build
