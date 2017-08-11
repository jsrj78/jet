#include "engine.h"

#include <assert.h>
#include <string.h>

static int printIndex;

void ResetPrint (void) {
    printIndex = 0;
    memset(g_PrintBuffer, 0, NMSGS * sizeof(Message));
}

static void PrintHandler (Gadget* gp, int inlet, Message msg) {
    assert(inlet == 0);

    Message arg = *(Message*) ExtraData(gp);
    if (arg != 0 && printIndex < NMSGS)
        g_PrintBuffer[printIndex++] = arg;
    if (printIndex < NMSGS)
        g_PrintBuffer[printIndex++] = msg;
}

static Gadget* MakePrintGadget (Message msg) {
    Gadget* gp = NewGadget(1, 0, sizeof(Message), PrintHandler);
    *(Message*) ExtraData(gp) = msg;
    return gp;
}

static Gadget* MakePassGadget (Message msg) {
    (void) msg;
    return NewGadget(1, 1, 0, Emit);
}

static Gadget* MakeInletGadget (Message msg) {
    (void) msg;
    return NewGadget(0, 1, sizeof(Wire), 0);
}

static void OutletHandler (Gadget* gp, int inlet, Message msg) {
    assert(inlet == 0);

    // scan the sibling gadgets to find the matching outlet
    Gadget **gpp = ExtraData(gp->parent);
    int outlet = 0;
    while (*gpp != gp)
        if ((*gpp++)->handler == OutletHandler)
            ++outlet;

    Emit(gp->parent, outlet, msg);
}

static Gadget* MakeOutletGadget (Message msg) {
    (void) msg;
    return NewGadget(1, 1, sizeof(Wire), OutletHandler);
}

static void SwapHandler (Gadget* gp, int inlet, Message msg) {
    switch (inlet) {
        case 0:
            Emit(gp, 1, msg);
            Emit(gp, 0, *(Message*) ExtraData(gp));
            break;
        case 1:
            *(Message*) ExtraData(gp) = msg;
            break;
        default:
            assert(0);
    }
}

static Gadget* MakeSwapGadget (Message msg) {
    Gadget* gp = NewGadget(1, 1, sizeof(Message), SwapHandler);
    *(Message*) ExtraData(gp) = msg;
    return gp;
}

struct Lookup_t g_Gadgets[] = {
    { "print",  MakePrintGadget  },
    { "pass",   MakePassGadget   },
    { "inlet",  MakeInletGadget  },
    { "outlet", MakeOutletGadget },
    { "swap",   MakeSwapGadget   },
    { 0, 0 }
};
