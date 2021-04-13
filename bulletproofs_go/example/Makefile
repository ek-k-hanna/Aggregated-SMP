.PHONY: clean

UNAME := $(shell uname)

ifeq ($(UNAME), Darwin)
    DYLIB_EXT := .dylib
else
    DYLIB_EXT := .so
endif

all: ../c_binding/libbulletproofs$(DYLIB_EXT) main

../c_binding/libbulletproofs$(DYLIB_EXT):
	$(MAKE) -C c_binding

main:
	g++ -o $@ main.cxx -I../c_binding -L../c_binding -lbulletproofs

clean:
	rm -f main ../c_binding/libbulletproofs.*
