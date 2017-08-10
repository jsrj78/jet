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

Gadget* NewGadget (uint8_t i, uint8_t o, uint16_t x,
                   void (*h)(Gadget*,int,Message)) {
    Gadget* gp = calloc(1, sizeof(Gadget) + x);
    gp->inlets = i;
    gp->outlets = o;
    gp->extra = x;
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
    uint16_t extra = (g+1) * sizeof(Gadget*);
    Circuit *cp = (Circuit*) NewGadget(i, o, extra, CircuitHandler);
    return cp;
}

void Add (Circuit* cp, int pos, Gadget* gp) {
    cp->child[pos] = gp;
    gp->parent = cp;
}

void AddWires (Gadget* gp, const Wire* w) {
    gp->wires = w;
}

void Feed (Gadget* gp, int inlet, Message msg) {
    gp->handler(gp, inlet, msg);
}

void Emit (Gadget* gp, int outlet, Message msg) {
    (void) gp;
    (void) outlet;
    (void) msg;
    for (const Wire* wp = gp->wires; wp != 0; ++wp)
        if (wp->gid == 255)
            break;
        else if (outlet == wp->out)
            Feed(gp->parent->child[wp->gid], wp->in, msg);
}
