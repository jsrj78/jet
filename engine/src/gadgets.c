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
    Gadget* gp = NewGadget(sizeof(Message), PrintHandler);
    *(Message*) ExtraData(gp) = msg;
    return gp;
}

static Gadget* MakePassGadget (Message msg) {
    (void) msg;
    return NewGadget(0, Emit);
}

static Gadget* MakeInletGadget (Message msg) {
    (void) msg;
    return NewGadget(sizeof(Wire), 0);
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
    return NewGadget(sizeof(Wire), OutletHandler);
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
    Gadget* gp = NewGadget(sizeof(Message), SwapHandler);
    *(Message*) ExtraData(gp) = msg;
    return gp;
}

static void ChangeHandler (Gadget* gp, int inlet, Message msg) {
    assert(inlet == 0);

    Message arg = *(Message*) ExtraData(gp);
    if (msg != arg) {
        *(Message*) ExtraData(gp) = msg;
        Emit(gp, 0, msg);
    }
}

static Gadget* MakeChangeGadget (Message msg) {
    (void) msg;
    Gadget* gp = NewGadget(sizeof(Message), ChangeHandler);
    *(Message*) ExtraData(gp) = -1;
    return gp;
}

static void MosesHandler (Gadget* gp, int inlet, Message msg) {
    assert(inlet == 0);

    Message arg = *(Message*) ExtraData(gp);
    Emit(gp, msg >= arg, msg);
}

static Gadget* MakeMosesGadget (Message msg) {
    Gadget* gp = NewGadget(sizeof(Message), MosesHandler);
    *(Message*) ExtraData(gp) = msg;
    return gp;
}

struct Lookup_t g_Gadgets[] = {
    { "print",  MakePrintGadget  },
    { "pass",   MakePassGadget   },
    { "inlet",  MakeInletGadget  },
    { "outlet", MakeOutletGadget },
    { "swap",   MakeSwapGadget   },
    { "change", MakeChangeGadget },
    { "moses",  MakeMosesGadget  },
    { 0, 0 }
};
