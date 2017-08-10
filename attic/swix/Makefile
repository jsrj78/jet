# Makefile for running tests, using http://cpputest.github.io

COMPONENT_NAME = swix

SRC_DIRS = src
TEST_SRC_DIRS = tests
INCLUDE_DIRS = $(CPPUTEST_HOME)/include src
CPPUTEST_HOME ?= /usr/local
CPPUTEST_CFLAGS += -std=gnu99

# see https://forums.pragprog.com/forums/123/topics/15200
export CPPUTEST_WARNINGFLAGS += -Wno-reserved-id-macro \
				-Wno-keyword-macro \
				-Wno-c++11-long-long

SILENCE = @

# file was not installed by "brew install cpputest", using a copy instead
include tests/MakefileWorker.mk

# this generates a tags file for vim
tags: .
	ctags -w */*.{h,c,cpp}
