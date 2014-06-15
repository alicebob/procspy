.PHONY: all
all: setup lib install

.PHONY: setup
setup: clean
	tar xf lsof_4.87_src.tar && mv lsof_4.87_src src
	chmod -R u+w src/*
	cd src && ./Configure -n `uname -s|tr A-Z a-z`
	#cd src && > main.c
	#cd src && > print.c
	# cd src && perl -i -pe's/main.[co]//g' Makefile
	cp print.c src/print.c
	cd src && make
	cd src && rm -f main.o
	cd src && ar cr libmylsof.a *.o
	cd src && ranlib libmylsof.a
	#cp src/libmylsof.a .
	#cp src/lib/liblsof.a .
	# Fix conflict with global regex.h, gobuild adds '-I .'
	cd src && mv regex.h localregex.h
	#sed -e 's/regex.h/localregex.h/' -i lsof.h
	cd src && perl -i -pe's/regex.h/localregex.h/' lsof.h

.PHONY: lib
lib:
	sudo install -D src/libmylsof.a /usr/local/lib/procspy/libmylsof.a
	sudo install -D src/lib/liblsof.a /usr/local/lib/procspy/liblsof.a

.PHONY: install
install:
	go install

.PHONY: clean
clean:
	rm -rf src/ procspy
