.PHONY: all
all: setup main

.PHONY: setup
setup: clean
	tar xf lsof_4.87_src.tar && mv lsof_4.87_src src
	cd src && ./Configure -n `uname -s|tr A-Z a-z`
	cd src && make # builds version.h
	cp src/*.[ch] .
	rm -f main.c
	rm -f arg.c util.c # some unneeded stuff.
	# Fix conflict with global regex.h, gobuild adds '-I .'
	mv regex.h localregex.h
	sed -i 's/regex.h/localregex.h/g' lsof.h

.PHONY: main
main:
	go build

.PHONY: clean
clean:
	rm -rf *.c *.h src/ lsof
