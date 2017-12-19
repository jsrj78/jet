# Top level makefile for the JET engine in C.

ARCH ?= none

ARCHS = linux macos

all: $(ARCH)
	
$(ARCHS): build
	cd build && make -f ../demos/Makefile $@

build:
	mkdir $@

none:
	@echo "Please do 'make ARCH' where ARCH is one of these:"
	@echo "  $(ARCHS)"
	@echo "(or use 'export ARCH=...' to select a default)"

clean:
	rm -rf build

.PHONY: all $(ARCHS) none clean
