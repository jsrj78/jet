#include "core.h"

#include <stdio.h>

void core_jetConfig (Config_t* cptr, Value_t info) {
    cptr->inlets = cptr->outlets = 1;
}

void core_jetDispatch (Gadget_t* gptr, int inlet, Value_t msg) {
}

void core_initConfig (Config_t* cptr, Value_t info) {
    cptr->outlets = 1;
}

void core_initDispatch (Gadget_t* gptr, int inlet, Value_t msg) {
}

void core_passConfig (Config_t* cptr, Value_t info) {
    cptr->inlets = cptr->outlets = 1;
}

void core_passDispatch (Gadget_t* gptr, int inlet, Value_t msg) {
    jEmit(gptr, 0, msg);
}

void core_printConfig (Config_t* cptr, Value_t info) {
    cptr->inlets = 1;
}

void core_printDispatch (Gadget_t* gptr, int inlet, Value_t msg) {
    printf("print: ");
    jPrint(msg);
    printf("\n");
}
