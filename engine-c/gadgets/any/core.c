#include "core.h"

#include <stdio.h>

void core_jetConfig (JConfig* cptr, JValue info) {
    cptr->inlets = cptr->outlets = 1;
}

void core_jetDispatch (JGadget* gptr, int inlet, JValue msg) {
}

void core_initConfig (JConfig* cptr, JValue info) {
    cptr->outlets = 1;
}

void core_initDispatch (JGadget* gptr, int inlet, JValue msg) {
}

void core_passConfig (JConfig* cptr, JValue info) {
    cptr->inlets = cptr->outlets = 1;
}

void core_passDispatch (JGadget* gptr, int inlet, JValue msg) {
    jEmit(gptr, 0, msg);
}

void core_printConfig (JConfig* cptr, JValue info) {
    cptr->inlets = 1;
}

void core_printDispatch (JGadget* gptr, int inlet, JValue msg) {
    printf("print: ref %u, val %d\n", msg.ref, msg.val);
}
