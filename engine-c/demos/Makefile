include ../arch

VPATH= .. ../gadgets/$(ARCH) ../gadgets/any ../demos/$(ARCH) ../demos/any
CFLAGS= -std=c99 -I.. -I../gadgets/$(ARCH) -I../gadgets/any

D_SRCS:= $(notdir $(wildcard ../demos/$(ARCH)/*.c)) \
         $(notdir $(wildcard ../demos/any/*.c))
D_EXES:= $(patsubst %.c,%,$(D_SRCS))

G_SRCS:= $(notdir $(wildcard ../gadgets/$(ARCH)/*.c)) \
         $(notdir $(wildcard ../gadgets/any/*.c))
G_OBJS:= $(patsubst %.c,%.o,$(G_SRCS))

all: demos

demos: $(D_EXES)

$(D_EXES): $(G_OBJS) jet.o

deps: $(D_SRCS) $(G_SRCS) jet.c
	@$(CC) $(CFLAGS) -MM $^ >$@

info:
	@echo "ARCH   = $(ARCH)"
	@echo "VPATH  = $(VPATH)"
	@echo "CFLAGS = $(CFLAGS)"
	@echo "LDLIBS = $(LDLIBS)"
	@echo "D_SRCS = $(D_SRCS)"
	@echo "D_EXES = $(D_EXES)"
	@echo "G_SRCS = $(G_SRCS)"
	@echo "G_OBJS = $(G_OBJS)"

.c.o:
	@echo [CC] $<
	@$(CC) -c $(CFLAGS) $< -o $@
%: %.o
	@echo [LD] $@
	@$(CC) $^ $(LDLIBS) -o $@
	
.PHONY: all demos info

# dependencies, auto-generated
-include deps
