#include "engine.h"

#include <stdlib.h>
#include <string.h>

Message g_PrintBuffer[NMSGS];

Gadget* LookupGadget (const char* name, Message arg) {
    for (struct Lookup_t* p = g_Gadgets; p->s != 0; ++p)
        if (strcmp(name, p->s) == 0)
            return p->c(arg);
    return 0;
}

Gadget* NewGadget (uint8_t i, uint8_t o, size_t x,
                   void (*h)(Gadget*,int,Message)) {
    Gadget* gp = calloc(1, sizeof(Gadget) + x);
    gp->inlets = i;
    gp->outlets = o;
    gp->extra = (uint16_t) x;
    gp->handler = h;
    return gp;
}

void* ExtraData(Gadget *gp) {
    return gp + 1;
}

void FreeGadget (Gadget* gp) {
    if (gp != 0) {
        if (gp->onFree != 0)
            gp->onFree(gp);
        free(gp);
    }
}

static void CircuitHandler (Gadget* gp, int inlet, Message msg) {
    Gadget **gpp = ExtraData(gp);
    switch (inlet) {
        case 0:
            Emit(gpp[0], 0, msg); // FIXME only correct for first inlet
    }
}

static void FreeCircuit (Gadget* cp) {
    if (cp != 0)
        for (Gadget** gpp = (Gadget**) ExtraData(cp); *gpp != 0; ++gpp)
            FreeGadget(*gpp);
}

Gadget* NewCircuit (uint8_t i, uint8_t o, uint8_t g) {
    uint16_t extra = (g+1) * sizeof(Gadget*);
    Gadget *cp = NewGadget(i, o, extra, CircuitHandler);
    cp->onFree = FreeCircuit;
    return cp;
}

void Add (Gadget* cp, Gadget* gp, const Wire* w) {
    Gadget **gpp = ExtraData(cp);
    while (*gpp != 0)
        ++gpp;
    *gpp = gp;
    gp->parent = cp;
    gp->wires = w;
    if (gp->onAdded != 0)
        gp->onAdded(gp);
}

void Feed (Gadget* gp, int inlet, Message msg) {
    gp->handler(gp, inlet, msg);
}

void Emit (Gadget* gp, int outlet, Message msg) {
    Gadget **gpp = ExtraData(gp->parent);
    for (const Wire* wp = gp->wires; wp != 0 && wp->gid != 255; ++wp)
        if (outlet == wp->out)
            Feed(gpp[wp->gid], wp->in, msg);
}
