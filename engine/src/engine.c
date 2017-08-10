#include "engine.h"

#include <stdlib.h>
#include <string.h>

Message g_PrintBuffer[NMSGS];

Gadget* LookupGadget (const char* name, Message arg) {
    for (struct Lookup_t* p = g_Gadgets; p->s != 0; ++p)
        if (strcmp(name, p->s) == 0)
            return (p->c)(arg);
    return 0;
}

Gadget* NewGadget (uint8_t i, uint8_t o, size_t x,
                   void (*h)(Gadget*,int,Message)) {
    Gadget* gp = calloc(1, sizeof(Gadget) + (o+1) * sizeof(Wire*) + x);
    gp->inlets = i;
    gp->outlets = o;
    gp->handler = h;
    return gp;
}

static void CircuitHandler (Gadget* gp, int inlet, Message msg) {
    switch (inlet) {
        case 0:
            (void) gp;
            (void) msg;
    }
}

Circuit* NewCircuit (uint8_t i, uint8_t o, uint8_t g) {
    size_t extra = (g+1) * sizeof(Gadget*);
    Circuit *cp = (Circuit*) NewGadget(i, o, extra, CircuitHandler);
    return cp;
}

void Add (Circuit* cp, int pos, Gadget* gp) {
    SetNthGadget(cp, pos, gp);
    gp->parent = cp;
}

void Feed (Gadget* gp, int inlet, Message msg) {
    gp->handler(gp, inlet, msg);
}

void Emit (Gadget* gp, int outlet, Message msg) {
    (void) gp;
    (void) outlet;
    (void) msg;
}
