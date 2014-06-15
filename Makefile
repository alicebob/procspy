.PHONY: all
all: setup procspy

.PHONY: setup
setup: clean
	tar xf lsof_4.87_src.tar && mv lsof_4.87_src src
	cd src && ./Configure -n `uname -s|tr A-Z a-z`
	cd src && make
	cd src && rm -f main.o
	cd src && ar cr libmylsof.a *.o
	cd src && ranlib libmylsof.a
	cp src/*.[h] .
	cp src/libmylsof.a .
	cp src/lib/liblsof.a .
	# Fix conflict with global regex.h, gobuild adds '-I .'
	mv regex.h localregex.h
	#sed -e 's/regex.h/localregex.h/' -i lsof.h
	perl -i -pe's/regex.h/localregex.h/' lsof.h

.PHONY: procspy
procspy:
	go build

.PHONY: clean
clean:
	rm -rf *.[ah] src/ procspy
