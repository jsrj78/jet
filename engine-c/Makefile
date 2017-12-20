# Top level makefile for the JET engine in C.
# Configures the architecture and defers to demos/ for the actual builds.

ARCH= none
-include arch

ARCHS= linux macos

default: $(ARCH)
	
$(ARCHS)::
	@echo 'ARCH=$@' >arch

$(ARCHS):: | build
	cd build && make -f ../demos/Makefile all

dry:
	cd build && make -f ../demos/Makefile -n

%::
	cd build && make -f ../demos/Makefile $@

build:
	mkdir $@

none:
	@echo "Please run 'make ARCH', where ARCH is one of:"
	@echo "  $(ARCHS)"
	@echo "(after that you can just use 'make')"

clean:
	rm -rf arch build

.PHONY: default arch $(ARCHS) dry none clean
