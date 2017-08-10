#include "engine.h"

static int printIndex;

void ResetPrint (void) {
    printIndex = 0;
}

static void PrintHandler (Gadget* gp, int inlet, Message msg) {
    Message arg = *((Message*) ExtraData(gp));
    switch (inlet) {
        case 0:
            if (arg != 0)
                g_PrintBuffer[printIndex++] = arg;
            if (printIndex < NMSGS)
                g_PrintBuffer[printIndex++] = msg;
    }
}

static Gadget* MakePrintGadget (Message msg) {
    Gadget* gp = NewGadget(1, 0, sizeof(Message), PrintHandler);
    *((Message*) ExtraData(gp)) = msg;
    return gp;
}

static void PassHandler (Gadget* gp, int inlet, Message msg) {
    switch (inlet) {
        case 0:
            Emit(gp, 0, msg);
    }
}

static Gadget* MakePassGadget (Message msg) {
    (void) msg;
    return NewGadget(1, 1, 0, PassHandler);
}

static void InletHandler (Gadget* gp, int inlet, Message msg) {
    (void) gp;
    (void) inlet;
    (void) msg;
}

static Gadget* MakeInletGadget (Message msg) {
    (void) msg;
    return NewGadget(0, 1, sizeof(Wire), InletHandler);
}

static void OutletHandler (Gadget* gp, int inlet, Message msg) {
    switch (inlet) {
        case 0:
            Emit(gp->parent, 0, msg); // FIXME only correct for first outlet
    }
}

static Gadget* MakeOutletGadget (Message msg) {
    (void) msg;
    return NewGadget(1, 1, sizeof(Wire), OutletHandler);
}

struct Lookup_t g_Gadgets[] = {
    { "print",  MakePrintGadget  },
    { "pass",   MakePassGadget   },
    { "inlet",  MakeInletGadget  },
    { "outlet", MakeOutletGadget },
    { 0, 0 }
};
