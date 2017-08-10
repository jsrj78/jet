#include "engine.h"

static int printIndex;

void ResetPrint (void) {
    printIndex = 0;
}

static void PrintHandler (Gadget* gp, int inlet, Message msg) {
    switch (inlet) {
        case 0:
            if (gp->arg != 0)
                g_PrintBuffer[printIndex++] = gp->arg;
            if (printIndex < NMSGS)
                g_PrintBuffer[printIndex++] = msg;
    }
}

static Gadget* MakePrintGadget (Message msg) {
    Gadget* gp = NewGadget(1, 0, 0, PrintHandler);
    gp->arg = msg;
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

struct Lookup_t g_Gadgets[] = {
    { "print", MakePrintGadget },
    { "pass", MakePassGadget },
    { 0, 0 }
};
