include ../arch

VPATH= .. ../gadgets/$(ARCH) ../gadgets/any ../demos/$(ARCH) ../demos/any
CFLAGS= -I.. -I../gadgets/$(ARCH) -I../gadgets/any

D_SRCS:= $(notdir $(wildcard ../demos/$(ARCH)/*.c)) \
         $(notdir $(wildcard ../demos/any/*.c))
D_EXES:= $(patsubst %.c,%,$(D_SRCS))

G_SRCS:= $(notdir $(wildcard ../gadgets/$(ARCH)/*.c)) \
         $(notdir $(wildcard ../gadgets/any/*.c))
G_OBJS:= $(patsubst %.c,%.o,$(G_SRCS))

all: demos
	# $(ARCH):
	@ls -l $(D_EXES)

demos: $(D_EXES)

$(D_EXES): $(G_OBJS) jet.o

deps: $(D_SRCS) $(G_SRCS) jet.c
	@$(CC) $(CFLAGS) -MM $^ >$@

info:
	@echo "ARCH   = $(ARCH)"
	@echo "VPATH  = $(VPATH)"
	@echo "CFLAGS = $(CFLAGS)"
	@echo "D_SRCS = $(D_SRCS)"
	@echo "D_EXES = $(D_EXES)"
	@echo "G_SRCS = $(G_SRCS)"
	@echo "G_OBJS = $(G_OBJS)"

.PHONY: all demos info

# dependencies, auto-generated
-include deps
