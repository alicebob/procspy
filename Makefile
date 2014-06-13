.PHONY: all
all: setup procspy

.PHONY: setup
setup: clean
	tar xf lsof_4.87_src.tar && mv lsof_4.87_src src
	cd src && ./Configure -n `uname -s|tr A-Z a-z`
	cd src && make # generates version.h and builds liblsof.a
	cp src/*.[ch] .
	cp src/lib/liblsof.a .
	rm -f main.c
	rm -f arg.c util.c # some unneeded stuff.
	# Fix conflict with global regex.h, gobuild adds '-I .'
	mv regex.h localregex.h
	#sed -e 's/regex.h/localregex.h/' -i lsof.h
	perl -i -pe's/regex.h/localregex.h/' lsof.h

.PHONY: procspy
procspy:
	go build

.PHONY: clean
clean:
	rm -rf *.c *.h src/ procspy
